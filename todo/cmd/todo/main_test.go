package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var binName = "todo"

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	testBin := filepath.Join(os.TempDir(), binName)
	build := exec.Command("go", "build", "-o", testBin)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Clean up
	os.Remove(testBin)
	os.Exit(code)
}

func TestTodoCLI(t *testing.T) {
	cmdPath := filepath.Join(os.TempDir(), binName)

	t.Run("AddNewTask", func(t *testing.T) {
		t.Setenv("TODO_FILENAME", filepath.Join(t.TempDir(), ".todo.json"))

		task := "test task number 1"
		cmd := exec.Command(cmdPath, "-add", task)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf(" 1: %s\n", task)
		if string(out) != expected {
			t.Errorf("Expected %q, got %q", expected, string(out))
		}
	})

	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		t.Setenv("TODO_FILENAME", filepath.Join(t.TempDir(), ".todo.json"))

		task := "task from stdin"
		cmd := exec.Command(cmdPath, "-add")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(stdin, task)
		stdin.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf(" 1: %s\n", task)
		if string(out) != expected {
			t.Errorf("Expected %q, got %q", expected, string(out))
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		t.Setenv("TODO_FILENAME", filepath.Join(t.TempDir(), ".todo.json"))

		task := "complete me"
		cmd := exec.Command(cmdPath, "-add", task)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		cmd = exec.Command(cmdPath, "-complete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("X 1: %s\n", task)
		if string(out) != expected {
			t.Errorf("Expected %q, got %q", expected, string(out))
		}
	})

	t.Run("DeleteTask", func(t *testing.T) {
		t.Setenv("TODO_FILENAME", filepath.Join(t.TempDir(), ".todo.json"))

		task1 := "task one"
		task2 := "task two"
		exec.Command(cmdPath, "-add", task1).Run()
		exec.Command(cmdPath, "-add", task2).Run()

		// Delete the first task
		cmd := exec.Command(cmdPath, "-del", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf(" 1: %s\n", task2)
		if string(out) != expected {
			t.Errorf("Expected %q, got %q", expected, string(out))
		}
	})
}
