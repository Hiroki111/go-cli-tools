package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

func main() {
	op := flag.String("op", "sum", "Operation to be executed")
	column := flag.Int("col", 1, "CSV column on which to execute operation")

	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fileNames []string, op string, column int, out io.Writer) error {
	var opFunc statsFunc

	if len(fileNames) == 0 {
		return ErrNoFiles
	}

	if column < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, column)
	}

	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	consolidate := make([]float64, 0)

	// NOTE: See this file in the previous commit. Perhaps using filesCh doesn't increase the performance.
	// I did benchmarking with -benchmem, but ns/op became worse unlike what I saw in the book.
	// The book's author says this filesCh approach is used because this program's goroutines are CPU-bound.
	// However, this program doesn't do heavy computation (e.g., Parsing several GB of files, doing expensive math).
	// To me, the goroutines here seem IO-bound, rather than CPU-bound.
	filesCh := make(chan string)
	resultCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	go func() {
		defer close(filesCh)
		for _, fileName := range fileNames {
			filesCh <- fileName
		}
	}()

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for fileName := range filesCh {
				// os.Open -> async
				file, err := os.Open(fileName)
				if err != nil {
					errCh <- fmt.Errorf("Cannot open file: %w", err)
					return
				}

				// csv2float -> async (It reads a file)
				data, err := csv2float(file, column)
				if err != nil {
					errCh <- err
				}

				// file.Close -> async
				if err := file.Close(); err != nil {
					errCh <- err
				}

				resultCh <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resultCh:
			consolidate = append(consolidate, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, opFunc(consolidate))
			return err
		}
	}
}
