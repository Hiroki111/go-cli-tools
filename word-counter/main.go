package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// "l" means that it counts lines. The default value is false, so the default behaviour is to count words.
	lines := flag.Bool("l", false, "Count lines")
	flag.Parse()
	fmt.Println(count(os.Stdin, *lines))
}

func count(r io.Reader, countLines bool) int {
	scanner := bufio.NewScanner(r)

	// If the count lines flag is not set, it counts words so it defines
	// the scanner split type to words (the default behaviour is to split by lines)
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	wordCount := 0

	for scanner.Scan() {
		wordCount++
	}

	return wordCount
}
