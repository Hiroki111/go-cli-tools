package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	var testCases = []struct {
		name        string
		project     string
		out         string
		expectedErr error
		setupGit    bool
		mockCmd     func(ctx context.Context, name string, arg ...string) *exec.Cmd
	}{
		{
			name:        "success",
			project:     "./testdata/tool",
			out:         "Go Build: SUCCESS\nGolangci-lint: SUCCESS\nGocyclo: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCCESS\nGit Push: SUCCESS\n",
			expectedErr: nil,
			setupGit:    true,
			mockCmd:     nil,
		},
		{
			name:        "successMock",
			project:     "./testdata/tool",
			out:         "Go Build: SUCCESS\nGolangci-lint: SUCCESS\nGocyclo: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCCESS\nGit Push: SUCCESS\n",
			expectedErr: nil,
			setupGit:    false,
			mockCmd:     mockCmdContext,
		},
		{
			name:        "fail",
			project:     "./testdata/toolErr",
			out:         "",
			expectedErr: &stepErr{step: "go build"},
			setupGit:    false,
			mockCmd:     nil,
		},
		{
			name:        "failFormat",
			project:     "./testdata/toolFmtErr",
			out:         "",
			expectedErr: &stepErr{step: "go fmt"},
			setupGit:    false,
			mockCmd:     nil,
		},
		{
			name:        "failTimeout",
			project:     "./testdata/tool",
			out:         "",
			expectedErr: context.DeadlineExceeded,
			setupGit:    false,
			mockCmd:     mockCmdTimeout,
		},
	}
	var branch = "master"

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setupGit {
				_, err := exec.LookPath("git")
				if err != nil {
					t.Skip("Git not installed. Skipping test.")
				}
				cleanup := setupGit(t, testCase.project)
				defer cleanup()
			}

			if testCase.mockCmd != nil {
				command = testCase.mockCmd
			}

			var out bytes.Buffer
			err := run(testCase.project, branch, &out)

			if testCase.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error: %q. Got 'nil' instead.", testCase.expectedErr)
					return
				}

				if !errors.Is(err, testCase.expectedErr) {
					t.Errorf("Expected error: %q. Got %q", testCase.expectedErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}

			if out.String() != testCase.out {
				t.Errorf("Expected output: %q. Got %q", testCase.out, out.String())
			}
		})
	}
}

func TestRunKill(t *testing.T) {
	var testCases = []struct {
		name        string
		project     string
		signal      syscall.Signal
		expectedErr error
	}{
		{"SIGINT", "./testdata/tool", syscall.SIGINT, ErrSignal},
		{"SIGTERM", "./testdata/tool", syscall.SIGTERM, ErrSignal},
		{"SIGQUIT", "./testdata/tool", syscall.SIGQUIT, nil},
	}
	var branch = "master"

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			command = mockCmdTimeout
		})

		errCh := make(chan error)
		irrelevantSignalCh := make(chan os.Signal, 1)
		expectedSignalCh := make(chan os.Signal, 1)

		signal.Notify(irrelevantSignalCh, syscall.SIGQUIT)
		defer signal.Stop(irrelevantSignalCh)

		signal.Notify(expectedSignalCh, testCase.signal)
		defer signal.Stop(expectedSignalCh)

		go func() {
			errCh <- run(testCase.project, branch, io.Discard)
		}()
		go func() {
			time.Sleep(2 * time.Second)
			_ = syscall.Kill(syscall.Getpid(), testCase.signal)
		}()

		select {
		case err := <-errCh:
			if err == nil {
				t.Errorf("Expected error. Got 'nil' instead.")
				return
			}

			if !errors.Is(err, testCase.expectedErr) {
				t.Errorf("Expected error: %q, got %q", testCase.expectedErr, err)
			}

			select {
			case receivedSignal := <-expectedSignalCh:
				if receivedSignal != testCase.signal {
					t.Errorf("Expected signal %q, got %q", testCase.signal, receivedSignal)
				}
			default:
				t.Errorf("Signal not received")
			}
		case <-irrelevantSignalCh:
		}
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("GO_HELPER_TIMEOUT") == "1" {
		time.Sleep(15 * time.Second)
	}

	if os.Args[2] == "git" {
		_, _ = fmt.Fprintln(os.Stdout, "Everything up-to-date")
		os.Exit(0)
	}

	os.Exit(1)
}

func setupGit(t *testing.T, project string) func() {
	t.Helper()

	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	tempDir, err := os.MkdirTemp("", "gocitest")
	if err != nil {
		t.Fatal(err)
	}

	projectPath, err := filepath.Abs(project)
	if err != nil {
		t.Fatal(err)
	}

	remoteURI := fmt.Sprintf("file://%s", tempDir)

	var gitCmdList = []struct {
		args []string
		dir  string
		env  []string
	}{
		{[]string{"init", "--bare"}, tempDir, nil},
		{[]string{"init"}, projectPath, nil},
		{[]string{"remote", "add", "origin", remoteURI}, projectPath, nil},
		{[]string{"add", "."}, projectPath, nil},
		{[]string{"commit", "-m", "test"}, projectPath,
			[]string{
				"GIT_COMMITTER_NAME=test",
				"GIT_COMMITTER_EMAIL=test@example.com",
				"GIT_AUTHOR_NAME=test",
				"GIT_AUTHOR_EMAIL=test@example.com",
			}},
	}

	for _, g := range gitCmdList {
		gitCmd := exec.Command(gitExec, g.args...)
		gitCmd.Dir = g.dir

		if g.env != nil {
			gitCmd.Env = append(os.Environ(), g.env...)
		}

		if err := gitCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}

	return func() {
		_ = os.RemoveAll(tempDir)
		_ = os.RemoveAll(filepath.Join(projectPath, ".git"))
	}
}

func mockCmdContext(ctx context.Context, exe string, args ...string) *exec.Cmd {
	arguments := []string{"-test.run=TestHelperProcess"}
	arguments = append(arguments, exe)
	arguments = append(arguments, args...)

	cmd := exec.CommandContext(ctx, os.Args[0], arguments...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func mockCmdTimeout(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cmd := mockCmdContext(ctx, exe, args...)
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")
	return cmd
}
