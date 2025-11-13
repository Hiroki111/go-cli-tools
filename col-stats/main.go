package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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

	for _, fileName := range fileNames {
		file, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("Cannot open file: %w", err)
		}

		data, err := csv2float(file, column)
		if err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}

		consolidate = append(consolidate, data...)
	}

	_, err := fmt.Fprintln(out, opFunc(consolidate))
	return err
}
