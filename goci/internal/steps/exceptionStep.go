package steps

import (
	"bytes"
	"fmt"
	"goci/internal/errors"
	"os/exec"
)

type ExceptionStep struct {
	Step
}

func NewExceptionStep(name, exe, message, project string, args []string) ExceptionStep {
	s := ExceptionStep{}
	s.Step = NewStep(name, exe, message, project, args)
	return s
}

func (s ExceptionStep) Execute() (string, error) {
	cmd := exec.Command(s.Exe, s.Args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = s.Project

	if err := cmd.Run(); err != nil {
		return "", &errors.StepErr{
			Step:    s.Name,
			Message: "failed to executre",
			Cause:   err,
		}
	}

	if out.Len() > 0 {
		return "", &errors.StepErr{
			Step:    s.Name,
			Message: fmt.Sprintf("invalid format: %s", out.String()),
			Cause:   nil,
		}
	}

	return s.Message, nil
}
