package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	var testCases = []struct {
		name        string
		project     string
		out         string
		expectedErr error
	}{
		{
			name:        "success",
			project:     "./testdata/tool",
			out:         "Go Build: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCCESS\n",
			expectedErr: nil,
		},
		{
			name:        "fail",
			project:     "./testdata/toolErr",
			out:         "",
			expectedErr: &stepErr{step: "go build"},
		},
		{
			name:        "failFormat",
			project:     "./testdata/toolFmtErr",
			out:         "",
			expectedErr: &stepErr{step: "go fmt"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
