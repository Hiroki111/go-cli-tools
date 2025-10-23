package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// "l" means that it counts lines
	isCountingLines := flag.Bool("l", false, "Count lines")
	// "b" means that it counts the number of bytes
	isCountingBytes := flag.Bool("b", false, "Count bytes")
	flag.Parse()
	if *isCountingBytes && *isCountingLines {
		fmt.Fprintln(os.Stderr, "It's not allowed to use -l and -b flags at the same time.")
		return
	}

	fmt.Println(count(os.Stdin, *isCountingLines, *isCountingBytes))
}

func count(r io.Reader, isCountingLines bool, isCountingBytes bool) int {

	if isCountingBytes {
		data, err := io.ReadAll(r)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			return 0
		}
		return len(data)
	}

	scanner := bufio.NewScanner(r)
	if !isCountingLines {
		scanner.Split(bufio.ScanWords)
	}

	wordCount := 0

	for scanner.Scan() {
		wordCount++
	}

	return wordCount
}
