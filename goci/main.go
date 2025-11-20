package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

type executer interface {
	execute() (string, error)
}

func main() {
	project := flag.String("p", "", "Project directory")
	branch := flag.String("b", "", "Branch name to push into")
	flag.Parse()

	if err := run(*project, *branch, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(project string, branch string, out io.Writer) error {
	if project == "" {
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}
	if branch == "" {
		return fmt.Errorf("branch name is required: %w", ErrValidation)
	}

	pipeline, err := loadPipeline(project, branch)
	if err != nil {
		return err
	}
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
