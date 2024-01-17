package utils

import (
	"os"
	"os/exec"
)

// ExecCommand executes a command and returns an error if it fails.
func ExecCommand(name string) error {
	cmd := exec.Command("sh", "-c", name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
