package utils

import (
	"bytes"
	"net"
	"strings"
	"unicode"

	"golang.org/x/crypto/ssh"
)

// RemoteCommand takes user, addr and privateKey and initiates an SSH session.
// It then runs the provided cmd and returns stdout, stderr output and error.
func RemoteCommand(user string, addr string, privateKey string, cmd string) (string, string, error) {
	// privateKey could be read from a file, or retrieved from another storage
	// source, such as the Secret Service / GNOME Keyring
	key, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return "", "", err
	}
	// Authentication
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	// Connect
	client, err := ssh.Dial("tcp", net.JoinHostPort(addr, "22"), config)
	if err != nil {
		return "", "", err
	}
	// Create a session. It is one session per command.
	session, err := client.NewSession()
	if err != nil {
		return "", "", err
	}
	defer session.Close()
	var out, stderr bytes.Buffer // import "bytes"
	session.Stdout = &out        // get output
	session.Stderr = &stderr
	// you can also pass what gets input to the stdin, allowing you to pipe
	// content from client to server
	//      session.Stdin = bytes.NewBufferString("My input")

	// Finally, run the command
	err = session.Run(cmd)

	// clean the output of non-printable characters
	cleanStdOut := strings.TrimFunc(out.String(), func(r rune) bool {
		return !unicode.IsGraphic(r)
	})

	cleanStdErr := strings.TrimFunc(stderr.String(), func(r rune) bool {
		return !unicode.IsGraphic(r)
	})

	return cleanStdOut, cleanStdErr, err
}
