package utils

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode"
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

// ExecCommandQuietly executes a command and returns an error if it fails without any stdout
func ExecCommandQuietly(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// ExecCommandWithReturn executes a command and returns the output as a string.
func ExecCommandWithReturn(name string) (string, error) {
	cmd := exec.Command("sh", "-c", name)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// clean the output of non-printable characters
	cleanStdOut := strings.TrimFunc(string(out), func(r rune) bool {
		return !unicode.IsGraphic(r)
	})

	return cleanStdOut, nil
}

// ReadLines reads lines from an io.Reader and returns them as a slice of strings.
func ReadLines(r io.Reader) ([]string, error) {
	var lines []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
