package main

import "os/exec"

type step struct {
	name    string
	exe     string
	args    []string
	message string
	project string
}

func newStep(name, exe, message, project string, args []string) step {
	return step{
		name:    name,
		exe:     exe,
		message: message,
		args:    args,
		project: project,
	}
}

func (s step) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)
	cmd.Dir = s.project

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", &stepErr{
			step:    s.name,
			message: string(output),
			cause:   err,
		}
	}

	return s.message, nil
}
