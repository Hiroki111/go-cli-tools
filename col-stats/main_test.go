package main

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name           string
		col            int
		op             string
		expectedResult string
		files          []string
		expectedErr    error
	}{
		{name: "RunAvg1File", col: 3, op: "avg", expectedResult: "227.6\n",
			files: []string{"./testdata/example.csv"}, expectedErr: nil},
		{name: "RunAvgMultiFiles", col: 3, op: "avg", expectedResult: "233.84\n",
			files: []string{"./testdata/example.csv", "./testdata/example2.csv"}, expectedErr: nil},
		{name: "RunFailRead", col: 2, op: "avg", expectedResult: "",
			files: []string{"./testdata/example.csv", "./testdata/fakefile.csv"}, expectedErr: os.ErrNotExist},
		{name: "RunFailColumn", col: 0, op: "avg", expectedResult: "",
			files: []string{"./testdata/example.csv"}, expectedErr: ErrInvalidColumn},
		{name: "RunFailNoFiles", col: 2, op: "avg", expectedResult: "",
			files: []string{}, expectedErr: ErrNoFiles},
		{name: "RunFailOperation", col: 2, op: "invalid", expectedResult: "",
			files: []string{"./testdata/example.csv"}, expectedErr: ErrInvalidOperation},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var result bytes.Buffer
			err := run(testCase.files, testCase.op, testCase.col, &result)

			if testCase.expectedErr != nil {
				if err == nil {
					t.Errorf("Expected error. Got nil instead")
				}

				if !errors.Is(err, testCase.expectedErr) {
					t.Errorf("Expected error %q, got %q instead", testCase.expectedErr, err)
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}

			if result.String() != testCase.expectedResult {
				t.Errorf("Expected %q, got %q instead", testCase.expectedResult, &result)
			}
		})
	}
}
