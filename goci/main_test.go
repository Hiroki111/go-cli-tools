package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	_, err := exec.LookPath("git")
	if err != nil {
		t.Skip("Git not installed. Skipping test.")
	}

	var testCases = []struct {
		name        string
		project     string
		out         string
		expectedErr error
		setupGit    bool
	}{
		{
			name:        "success",
			project:     "./testdata/tool",
			out:         "Go Build: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCCESS\nGit Push: SUCCESS\n",
			expectedErr: nil,
			setupGit:    true,
		},
		{
			name:        "fail",
			project:     "./testdata/toolErr",
			out:         "",
			expectedErr: &stepErr{step: "go build"},
			setupGit:    false,
		},
		{
			name:        "failFormat",
			project:     "./testdata/toolFmtErr",
			out:         "",
			expectedErr: &stepErr{step: "go fmt"},
			setupGit:    false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setupGit {
				cleanup := setupGit(t, testCase.project)
				defer cleanup()
			}

			var out bytes.Buffer
			err := run(testCase.project, &out)

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
		os.RemoveAll(tempDir)
		os.RemoveAll(filepath.Join(projectPath, ".git"))
	}
}
