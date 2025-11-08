package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{name: "NoFilter", root: "testdata",
			cfg:      config{ext: "", size: 0, list: true},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", root: "testdata",
			cfg:      config{ext: ".log", size: 0, list: true},
			expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeMatch", root: "testdata",
			cfg:      config{ext: ".log", size: 10, list: true},
			expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeNoMatch", root: "testdata",
			cfg:      config{ext: ".log", size: 20, list: true},
			expected: ""},
		{name: "FilterExtensionNoMatch", root: "testdata",
			cfg:      config{ext: ".gz", size: 0, list: true},
			expected: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(testCase.root, &buffer, testCase.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if testCase.expected != res {
				t.Errorf("Expected %q, got %q instead\n", testCase.expected, res)
			}
		})
	}
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		name               string
		cfg                config
		extToKeep          string
		numOfFilesToDelete int
		numOfFilesToKeep   int
		expected           string
	}{
		{
			name:               "DeleteExtensionNoMatch",
			cfg:                config{ext: ".log", del: true},
			extToKeep:          ".gz",
			numOfFilesToDelete: 0,
			numOfFilesToKeep:   10,
			expected:           "",
		},
		{
			name:               "DeleteExtensionMatch",
			cfg:                config{ext: ".log", del: true},
			extToKeep:          "",
			numOfFilesToDelete: 10,
			numOfFilesToKeep:   0,
			expected:           "",
		},
		{
			name:               "DeleteExtensionMixed",
			cfg:                config{ext: ".log", del: true},
			extToKeep:          ".gz",
			numOfFilesToDelete: 5,
			numOfFilesToKeep:   5,
			expected:           "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)

			testCase.cfg.wLog = &logBuffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				testCase.cfg.ext:   testCase.numOfFilesToDelete,
				testCase.extToKeep: testCase.numOfFilesToKeep,
			})
			defer cleanup()

			if err := run(tempDir, &buffer, testCase.cfg); err != nil {
				t.Fatal(err)
			}

			result := buffer.String()

			if testCase.expected != result {
				t.Errorf("Expected %q, got %q instead\n", testCase.expected, result)
			}

			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != testCase.numOfFilesToKeep {
				t.Errorf("Expected %d files left, got %d instead\n", testCase.numOfFilesToKeep, len(filesLeft))
			}

			expectedLogLines := testCase.numOfFilesToDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expectedLogLines {
				t.Errorf("Expected %d log lines, got %d instead\n", expectedLogLines, len(lines))
			}
		})
	}
}

func createTempDir(t *testing.T, files map[string]int) (dirName string, cleanup func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fileName := fmt.Sprintf("file%d%s", j, k)
			filePath := filepath.Join(tempDir, fileName)
			if err := os.WriteFile(filePath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	return tempDir, func() { os.RemoveAll(tempDir) }
}
