package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/Pixie2468/git-branch-explorer/internal/commands"
)

type focusState int

const (
	focusBranches focusState = iota
	focusCommits
	focusDetails
)

type model struct {
	focus focusState

	branches     []string
	branchCursor int
	branchOffset int

	commits      []string
	commitCursor int
	commitOffset int

	commitDetails     string
	detailLines       []string // word-wrapped lines ready to render
	detailScrollOffset int     // how many lines scrolled down in details

	loading bool
	errText string

	width  int
	height int

	svc *commands.Service
}

func InitialModel(gitService *commands.Service) model {
	return model{
		svc:     gitService,
		focus:   focusBranches,
		loading: true,
		width:   80,
		height:  24,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		loadBranches(m.svc),
		loadCommits(m.svc),
	)
}

// listViewport is the number of list rows that actually fit in the pane.
//
// Pane layout on screen (outside → inside):
//
//	border-top   (1)
//	pad-top      (1)
//	title        (1)
//	title-margin (1)   ← MarginBottom(1) on titleStyle
//	[ list rows ]      ← this is what we calculate
//	pad-bottom   (1)
//	border-bot   (1)
//	             ─────
//	             6 rows of chrome
//
// The view wraps everything in "\n" + ui + "\n" + helpbar + "\n" = 3 extra rows.
// Total pane screen height = m.height - 3
// List rows = (m.height - 3) - 6 = m.height - 9
func (m model) listViewport() int {
	n := m.height - 9
	if n < 3 {
		return 3
	}
	return n
}

// detailViewport is the number of detail text lines that fit in the pane
// (same chrome math as listViewport, title takes 2 rows).
func (m model) detailViewport() int {
	return m.listViewport()
}

// clampOffset keeps the cursor inside the visible window.
func clampOffset(cursor, offset, viewport, total int) int {
	if total <= viewport {
		return 0
	}
	if cursor < offset {
		return cursor
	}
	if cursor >= offset+viewport {
		return cursor - viewport + 1
	}
	return offset
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Re-wrap detail lines for the new width
		if m.commitDetails != "" {
			m.detailLines = wrapDetailLines(m.commitDetails, m.detailWidth())
		}
		return m, nil

	case BranchesDataMsg:
		m.branches = msg.Data
		if m.commits != nil {
			m.loading = false
		}

	case CommitsDataMsg:
		m.commits = msg.Data
		m.loading = false
		if len(m.commits) > 0 {
			return m, loadCommitDetails(m.svc, m.commits[0])
		}

	case DetailsDataMsg:
		m.commitDetails = msg.Data
		m.detailScrollOffset = 0
		m.detailLines = wrapDetailLines(m.commitDetails, m.detailWidth())

	case errMsg:
		m.errText = msg.err.Error()
		m.loading = false

	case tea.KeyMsg:
		m.errText = ""

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		// ── Focus switching ──────────────────────────────────────
		case "tab":
			switch m.focus {
			case focusBranches:
				m.focus = focusCommits
			case focusCommits:
				m.focus = focusDetails
			case focusDetails:
				m.focus = focusBranches
			}

		case "shift+tab":
			switch m.focus {
			case focusBranches:
				m.focus = focusDetails
			case focusCommits:
				m.focus = focusBranches
			case focusDetails:
				m.focus = focusCommits
			}

		// ── Navigation ──────────────────────────────────────────
		case "up", "k":
			vp := m.listViewport()
			switch m.focus {
			case focusBranches:
				if m.branchCursor > 0 {
					m.branchCursor--
					m.branchOffset = clampOffset(m.branchCursor, m.branchOffset, vp, len(m.branches))
				}
			case focusCommits:
				if m.commitCursor > 0 {
					m.commitCursor--
					m.commitOffset = clampOffset(m.commitCursor, m.commitOffset, vp, len(m.commits))
				}
			case focusDetails:
				if m.detailScrollOffset > 0 {
					m.detailScrollOffset--
				}
			}

		case "down", "j":
			vp := m.listViewport()
			switch m.focus {
			case focusBranches:
				if m.branchCursor < len(m.branches)-1 {
					m.branchCursor++
					m.branchOffset = clampOffset(m.branchCursor, m.branchOffset, vp, len(m.branches))
				}
			case focusCommits:
				if m.commitCursor < len(m.commits)-1 {
					m.commitCursor++
					m.commitOffset = clampOffset(m.commitCursor, m.commitOffset, vp, len(m.commits))
				}
			case focusDetails:
				maxScroll := len(m.detailLines) - m.detailViewport()
				if maxScroll > 0 && m.detailScrollOffset < maxScroll {
					m.detailScrollOffset++
				}
			}

		// ── Actions ─────────────────────────────────────────────
		case "enter":
			switch m.focus {
			case focusBranches:
				if len(m.branches) > 0 {
					if err := m.svc.Switch(m.branches[m.branchCursor]); err != nil {
						m.errText = err.Error()
						return m, nil
					}
					m.commitCursor = 0
					m.commitOffset = 0
					m.commitDetails = ""
					m.detailLines = nil
					m.detailScrollOffset = 0
					return m, loadCommits(m.svc)
				}
			case focusCommits:
				if len(m.commits) > 0 {
					m.detailScrollOffset = 0
					return m, loadCommitDetails(m.svc, m.commits[m.commitCursor])
				}
			}
		}
	}

	return m, nil
}

// detailWidth returns the available inner width for the details pane text.
// Must stay in sync with view.go's detailW calculation.
func (m model) detailWidth() int {
	usable := m.width - 12
	if usable < 30 {
		usable = 30
	}
	bW := usable * 22 / 100
	cW := usable * 22 / 100
	dW := usable - bW - cW
	if dW < 16 {
		dW = 16
	}
	// subtract horizontal padding (2 left + 2 right = 4) from the content width
	w := dW - 4
	if w < 8 {
		w = 8
	}
	return w
}
