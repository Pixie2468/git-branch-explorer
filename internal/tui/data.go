package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/Pixie2468/git-branch-explorer/internal/commands"
)

// BranchesDataMsg is sent when the branch list has been fetched.
type BranchesDataMsg struct{ Data []string }

// CommitsDataMsg is sent when the commit hash list has been fetched.
type CommitsDataMsg struct{ Data []string }

// DetailsDataMsg is sent when a commit's details have been fetched.
type DetailsDataMsg struct{ Data string }

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func loadBranches(s *commands.Service) tea.Cmd {
	return func() tea.Msg {
		raw, err := s.List()
		if err != nil {
			return errMsg{err: err}
		}
		cleaned := strings.TrimSpace(string(raw))
		if cleaned == "" {
			return BranchesDataMsg{Data: nil}
		}
		return BranchesDataMsg{Data: strings.Split(cleaned, "\n")}
	}
}

func loadCommits(s *commands.Service) tea.Cmd {
	return func() tea.Msg {
		raw, err := s.GetCommitHashes()
		if err != nil {
			return errMsg{err: err}
		}
		cleaned := strings.TrimSpace(string(raw))
		if cleaned == "" {
			return CommitsDataMsg{Data: nil}
		}
		return CommitsDataMsg{Data: strings.Split(cleaned, "\n")}
	}
}

func loadCommitDetails(s *commands.Service, commitHash string) tea.Cmd {
	return func() tea.Msg {
		details, err := s.GetCommitDetails(commitHash)
		if err != nil {
			return errMsg{err: err}
		}
		return DetailsDataMsg{Data: details}
	}
}
