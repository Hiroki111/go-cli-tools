package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3 word4")

	exp := 4

	res := count(b, false, false)

	if res != exp {
		t.Errorf("Expected %d, got %d instead. \n", exp, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1 word 2 word3\nline2\nline3 word1")

	exp := 3

	res := count(b, true, false)

	if res != exp {
		t.Errorf("Expected %d, got %d instead. \n", exp, res)
	}
}

func TestCountBytes(t *testing.T) {
	// a -> 1 byte, ü -> 2 bytes, あ -> 3 bytes, 😀 -> 4 bytes
	b := bytes.NewBufferString("a ü あ 😀")

	exp := 13

	res := count(b, false, true)

	if res != exp {
		t.Errorf("Expected %d, got %d instead. \n", exp, res)
	}
}

func TestCountWithConflictingFlags(t *testing.T) {
	b := bytes.NewBufferString("dummy input")

	// Create a pipe to capture stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Replace os.Stderr with the writable end of the pipe
	oldStderr := os.Stderr
	os.Stderr = w

	// Run the function under test
	got := count(b, true, true)

	// Close the writer and restore stderr
	w.Close()
	os.Stderr = oldStderr

	// Read what was written to the pipe
	var stderr bytes.Buffer
	_, err = io.Copy(&stderr, r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	// Assertions
	if got != 0 {
		t.Errorf("Expected return value 0 for conflicting flags, got %d", got)
	}

	want := "It's not allowed to use -l and -b flags at the same time."
	if !bytes.Contains(stderr.Bytes(), []byte(want)) {
		t.Errorf("Expected stderr to contain %q, got %q", want, stderr.String())
	}
}
