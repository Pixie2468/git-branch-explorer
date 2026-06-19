package commands

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type Service struct {
	Path    string
	Timeout time.Duration
}

func (s *Service) git(args ...string) ([]byte, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	// Only apply timeout if it is greater than 0
	if s.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), s.Timeout)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = s.Path

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("git execution failed: %w (stderr: %s)", err, stderr.String())
	}

	return stdout.Bytes(), nil
}
