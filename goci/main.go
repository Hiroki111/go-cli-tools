package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type executer interface {
	execute() (string, error)
}

func main() {
	project := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*project, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(project string, out io.Writer) error {
	if project == "" {
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}

	gocycloUpperLimit := 20
	pipeline := make([]executer, 6)
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		project,
		[]string{"build", ".", "errors"},
	)
	pipeline[1] = newStep(
		"golangci-lint run",
		"golangci-lint",
		"Golangci-lint: SUCCESS",
		project,
		[]string{"run"},
	)
	pipeline[2] = newStep(
		fmt.Sprintf("gocyclo -over %d", gocycloUpperLimit),
		"gocyclo",
		"Gocyclo: SUCCESS",
		project,
		[]string{"-over", strconv.Itoa(gocycloUpperLimit), "."},
	)
	pipeline[3] = newStep(
		"go test",
		"go",
		"Go Test: SUCCESS",
		project,
		[]string{"test", "-v"},
	)
	pipeline[4] = newExceptionStep(
		"go fmt",
		"gofmt",
		"Gofmt: SUCCESS",
		project,
		[]string{"-l", "."},
	)
	pipeline[5] = newTimeoutStep(
		"git push",
		"git",
		"Git Push: SUCCESS",
		project,
		[]string{"push", "origin", "master"},
		10*time.Second,
	)
	signalCh := make(chan os.Signal, 1)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for _, s := range pipeline {
			message, err := s.execute()
			if err != nil {
				errCh <- err
				return
			}

			_, err = fmt.Fprintln(out, message)
			if err != nil {
				errCh <- err
				return
			}
		}
		close(doneCh)
	}()

	for {
		select {
		case receivedSignal := <-signalCh:
			signal.Stop(signalCh)
			return fmt.Errorf("%s: Exiting: %w", receivedSignal, ErrSignal)
		case err := <-errCh:
			return err
		case <-doneCh:
			return nil
		}
	}
}
