package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
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
		return fmt.Errorf("Project directory is required: %w", ErrValidation)
	}

	pipeline := make([]executer, 4)
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		project,
		[]string{"build", ".", "errors"},
	)
	pipeline[1] = newStep(
		"go test",
		"go",
		"Go Test: SUCCESS",
		project,
		[]string{"test", "-v"},
	)
	pipeline[2] = newExceptionStep(
		"go fmt",
		"gofmt",
		"Gofmt: SUCCESS",
		project,
		[]string{"-l", "."},
	)
	pipeline[3] = newTimeoutStep(
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
