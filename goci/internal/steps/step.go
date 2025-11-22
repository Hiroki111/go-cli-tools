package steps

import (
	"goci/internal/errors"
	"os/exec"
)

type Step struct {
	Name    string
	Exe     string
	Args    []string
	Message string
	Project string
}

func NewStep(name, exe, message, project string, args []string) Step {
	return Step{
		Name:    name,
		Exe:     exe,
		Message: message,
		Args:    args,
		Project: project,
	}
}

func (s Step) Execute() (string, error) {
	cmd := exec.Command(s.Exe, s.Args...)
	cmd.Dir = s.Project

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", &errors.StepErr{
			Step:    s.Name,
			Message: string(output),
			Cause:   err,
		}
	}

	return s.Message, nil
}
