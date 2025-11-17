package main

import (
	"context"
	"os/exec"
	"time"
)

type timeoutStep struct {
	step
	timeout time.Duration
}

func newTimeoutStep(name, exe, message, project string, args []string, timeout time.Duration) timeoutStep {
	s := timeoutStep{}
	s.step = newStep(name, exe, message, project, args)
	s.timeout = timeout

	if s.timeout == 0 {
		s.timeout = 30 * time.Second
	}

	return s
}

var command = exec.CommandContext

func (s timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := command(ctx, s.exe, s.args...)
	cmd.Dir = s.project

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", &stepErr{
				step:    s.name,
				message: "failed time out",
				cause:   context.DeadlineExceeded,
			}
		}

		return "", &stepErr{
			step:    s.name,
			message: "failed to execute",
			cause:   err,
		}
	}

	return s.message, nil
}
