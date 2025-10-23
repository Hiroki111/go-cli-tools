package main

import (
	"bytes"
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
