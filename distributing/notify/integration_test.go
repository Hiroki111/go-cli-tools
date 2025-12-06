//go:build integration
// +build integration

// NOTE: This doesn't work in wsl2. `notify-send` is a program to send desktop notifications. See https://manpages.ubuntu.com/manpages/resolute/en/man1/notify-send.1.html

package notify_test

import (
	"fmt"
	"testing"

	"notify"
)

func TestSend(t *testing.T) {
	n := notify.New("test title", "test msg", notify.SeverityNormal)

	err := n.Send()

	if err != nil {
		fmt.Printf("error: %s\n", err)
		t.Error(err)
	}
}
