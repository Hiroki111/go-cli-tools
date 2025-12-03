package pomodoro_test

import (
	"pomo/pomodoro"
	"pomo/pomodoro/repository"
	"testing"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	return repository.NewInMemoryRepo(), func() {}
}
