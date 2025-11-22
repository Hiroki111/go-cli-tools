package steps

import (
	"context"
	"goci/internal/errors"
	"os/exec"
	"time"
)

type TimeoutStep struct {
	Step
	Timeout time.Duration
}

func NewTimeoutStep(name, exe, message, project string, args []string, timeout time.Duration) TimeoutStep {
	s := TimeoutStep{}
	s.Step = NewStep(name, exe, message, project, args)
	s.Timeout = timeout

	if s.Timeout == 0 {
		s.Timeout = 30 * time.Second
	}

	return s
}

var Command = exec.CommandContext

func (s TimeoutStep) Execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	cmd := Command(ctx, s.Exe, s.Args...)
	cmd.Dir = s.Project

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", &errors.StepErr{
				Step:    s.Name,
				Message: "failed time out",
				Cause:   context.DeadlineExceeded,
			}
		}

		return "", &errors.StepErr{
			Step:    s.Name,
			Message: "failed to execute",
			Cause:   err,
		}
	}

	return s.Message, nil
}
