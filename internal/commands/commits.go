package commands

import (
	"fmt"
	"strings"
)

func (s *Service) GetCommitHashes() ([]byte, error) {
	return s.git(
		"log",
		"--format=%h",
	)
}

// GetCommitDetails returns a formatted summary of a single commit including
// hash, author, date, and message.
func (s *Service) GetCommitDetails(commitHash string) (string, error) {
	out, err := s.git(
		"show",
		"-s",
		"--format=%H%n%an%n%ar%n%s%n%n%b",
		commitHash,
	)
	if err != nil {
		return "", err
	}

	raw := strings.TrimSpace(string(out))
	parts := strings.SplitN(raw, "\n", 5)

	// Graceful fallback if format is unexpected
	if len(parts) < 4 {
		return raw, nil
	}

	hash := parts[0]
	author := parts[1]
	date := parts[2]
	subject := parts[3]

	var body string
	if len(parts) == 5 {
		body = strings.TrimSpace(parts[4])
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Commit:  %s\n", hash)
	fmt.Fprintf(&b, "Author:  %s\n", author)
	fmt.Fprintf(&b, "Date:    %s\n", date)
	b.WriteString("\n")
	b.WriteString(subject)
	if body != "" {
		b.WriteString("\n\n")
		b.WriteString(body)
	}

	return b.String(), nil
}

func (s *Service) GetRecentCommitHash() (string, error) {
	out, err := s.git(
		"rev-parse",
		"--short",
		"HEAD",
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
