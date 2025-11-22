package errors

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation failed")
	ErrSignal     = errors.New("received signal")
)

type StepErr struct {
	Step    string
	Message string
	Cause   error
}

func (s *StepErr) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", s.Step, s.Message, s.Cause)
}

func (s *StepErr) Is(target error) bool {
	t, ok := target.(*StepErr)
	if !ok {
		return false
	}

	return t.Step == s.Step
}

func (s *StepErr) Unwrap() error {
	return s.Cause
}
