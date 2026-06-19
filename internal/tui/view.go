package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	if m.loading {
		v := tea.NewView("\n  Loading...\n")
		v.AltScreen = true
		return v
	}

	vp := m.listViewport() // rows available inside each list pane

	// ── Fixed pane widths (percentage of terminal) ──────────────
	// 3 panes × (2 border + 2 padding) = 12 chars of horizontal chrome
	usable := max(m.width-12, 30)
	branchW := usable * 22 / 100
	commitW := usable * 22 / 100
	detailW := usable - branchW - commitW
	if branchW < 12 {
		branchW = 12
	}
	if commitW < 10 {
		commitW = 10
	}
	if detailW < 16 {
		detailW = 16
	}

	paneH := max(m.height-5, 7)

	// ── Branches pane ───────────────────────────────────────────
	branchPane := paneStyle(m.focus == focusBranches).
		Width(branchW).
		Height(paneH).
		Render(
			titleStyle.Render("Branches") + "\n" +
				renderList(m.branches, m.branchCursor, m.branchOffset, vp, m.focus == focusBranches),
		)

	// ── Commits pane ────────────────────────────────────────────
	commitPane := paneStyle(m.focus == focusCommits).
		Width(commitW).
		Height(paneH).
		Render(
			titleStyle.Render("Commits") + "\n" +
				renderList(m.commits, m.commitCursor, m.commitOffset, vp, m.focus == focusCommits),
		)

	// ── Details pane ────────────────────────────────────────────
	dvp := m.detailViewport()
	detailContent := renderDetailLines(m.detailLines, m.errText, m.detailScrollOffset, dvp)
	detailPane := paneStyle(m.focus == focusDetails).
		Width(detailW).
		Height(paneH).
		Render(
			titleStyle.Render("Commit Details") + "\n" +
				detailContent,
		)

	// ── Compose ─────────────────────────────────────────────────
	ui := lipgloss.JoinHorizontal(lipgloss.Top, branchPane, commitPane, detailPane)
	help := dimStyle.Render(" [Tab] Next Pane  [↑↓] Scroll  [Enter] Select  [q] Quit")

	v := tea.NewView("\n" + ui + "\n" + help + "\n")
	v.AltScreen = true
	return v
}

// renderList renders exactly `viewport` rows of content (no more, no less —
// lipgloss clips at paneH anyway, but we control it explicitly).
//
// Row budget:
//   - If there are items above → 1 row for ↑ indicator
//   - If there are items below → 1 row for ↓ indicator
//   - Remaining rows → list items
func renderList(items []string, cursor, offset, viewport int, focused bool) string {
	if len(items) == 0 {
		return dimStyle.Render("  (none)")
	}

	showTop := offset > 0
	// Tentatively check if bottom indicator is needed
	showBottom := offset+viewport < len(items)

	// Reserve rows for indicators
	slots := viewport
	if showTop {
		slots--
	}
	if showBottom {
		slots--
	}
	if slots < 1 {
		slots = 1
	}

	end := offset + slots
	if end > len(items) {
		end = len(items)
	}
	// Recalculate bottom after final end
	showBottom = end < len(items)

	var b strings.Builder

	if showTop {
		b.WriteString(scrollStyle.Render(fmt.Sprintf("  ↑ %d above", offset)))
		b.WriteByte('\n')
	}

	for i := offset; i < end; i++ {
		item := truncate(items[i], viewport) // prevent wide entries from wrapping
		switch {
		case i == cursor && focused:
			b.WriteString(cursorStyle.Render("▶ " + item))
		case i == cursor:
			b.WriteString(lipgloss.NewStyle().Foreground(inactiveColor).Render("▶ " + item))
		default:
			b.WriteString("  " + item)
		}
		b.WriteByte('\n')
	}

	if showBottom {
		b.WriteString(scrollStyle.Render(fmt.Sprintf("  ↓ %d below", len(items)-end)))
	}

	return b.String()
}

// renderDetailLines renders a window of pre-wrapped detail lines,
// with scroll indicators.
func renderDetailLines(lines []string, errText string, offset, viewport int) string {
	if errText != "" {
		return dimStyle.Render("  ⚠ " + errText)
	}
	if len(lines) == 0 {
		return dimStyle.Render("  Press Enter on a commit")
	}

	showTop := offset > 0
	showBottom := offset+viewport < len(lines)

	slots := viewport
	if showTop {
		slots--
	}
	if showBottom {
		slots--
	}
	if slots < 1 {
		slots = 1
	}

	end := offset + slots
	if end > len(lines) {
		end = len(lines)
	}
	showBottom = end < len(lines)

	var b strings.Builder

	if showTop {
		b.WriteString(scrollStyle.Render(fmt.Sprintf("  ↑ %d above", offset)))
		b.WriteByte('\n')
	}

	for _, l := range lines[offset:end] {
		b.WriteString(detailStyle.Render(l))
		b.WriteByte('\n')
	}

	if showBottom {
		b.WriteString(scrollStyle.Render(fmt.Sprintf("  ↓ %d below", len(lines)-end)))
	}

	return b.String()
}

// wrapDetailLines word-wraps text into lines of maxWidth and returns the slice.
// Called once when new details arrive or the window is resized.
func wrapDetailLines(text string, maxWidth int) []string {
	if maxWidth < 4 {
		maxWidth = 4
	}
	var result []string
	for _, para := range strings.Split(text, "\n") {
		if strings.TrimSpace(para) == "" {
			result = append(result, "")
			continue
		}
		words := strings.Fields(para)
		line := words[0]
		for _, w := range words[1:] {
			if len(line)+1+len(w) > maxWidth {
				result = append(result, line)
				line = w
			} else {
				line += " " + w
			}
		}
		result = append(result, line)
	}
	return result
}

// truncate caps a string at maxLen runes so a long branch/commit
// name never wraps inside its pane column.
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen < 4 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-1]) + "…"
}
