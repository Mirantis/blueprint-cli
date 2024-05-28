package utils

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/ssh"
)

func RemotePowershellScript(user string, addr string, privateKey string, scriptPath string) (string, string, error) {

	pwd, err := ExecCommandWithReturn("pwd")
	log.Info().Msgf("pwd is %s", pwd)

	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		return "", "", err
	}

	log.Info().Msgf("Running command script on %s@%s", user, addr)

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
	session.Stdin = bytes.NewBuffer(scriptContent)

	err = session.Run("powershell -nologo -noprofile")
	if err != nil {
		return "", "", err
	}

	cleanStdOut := strings.TrimFunc(out.String(), func(r rune) bool {
		return !unicode.IsGraphic(r)
	})

	cleanStdErr := strings.TrimFunc(stderr.String(), func(r rune) bool {
		return !unicode.IsGraphic(r)
	})

	return cleanStdOut, cleanStdErr, err
}

//
//func RemoteScript(user string, addr string, privateKey string, scriptPath string) (string, string, error) {
//
//	pwd, err := ExecCommandWithReturn("pwd")
//	log.Info().Msgf("pwd is %s", pwd)
//
//	scriptContent, err := os.ReadFile(scriptPath)
//	if err != nil {
//		return "", "", err
//	}
//
//	// privateKey could be read from a file, or retrieved from another storage
//	// source, such as the Secret Service / GNOME Keyring
//	key, err := ssh.ParsePrivateKey([]byte(privateKey))
//	if err != nil {
//		return "", "", err
//	}
//	// Authentication
//	config := &ssh.ClientConfig{
//		User:            user,
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//		Auth: []ssh.AuthMethod{
//			ssh.PublicKeys(key),
//		},
//	}
//	// Connect
//	client, err := ssh.Dial("tcp", net.JoinHostPort(addr, "22"), config)
//	if err != nil {
//		return "", "", err
//	}
//	// Create a session. It is one session per command.
//	session, err := client.NewSession()
//	if err != nil {
//		return "", "", err
//	}
//	defer session.Close()
//	var out, stderr bytes.Buffer // import "bytes"
//	session.Stdout = &out        // get output
//	session.Stderr = &stderr
//	// you can also pass what gets input to the stdin, allowing you to pipe
//	// content from client to server
//	//      session.Stdin = bytes.NewBufferString("My input")
//
//	// Write the script content to the remote shell
//	stdin, err := session.StdinPipe()
//	if err != nil {
//		return "", "", err
//	}
//	defer stdin.Close()
//
//	// Start a shell
//	err = session.Shell()
//	if err != nil {
//		return "", "", err
//	}
//
//	_, err = io.Copy(stdin, strings.NewReader(string(scriptContent)))
//	if err != nil {
//		return "", "", err
//	}
//
//	// Close the shell
//	err = session.Wait()
//	if err != nil {
//		return "", "", err
//	}
//
//	// clean the output of non-printable characters
//	cleanStdOut := strings.TrimFunc(out.String(), func(r rune) bool {
//		return !unicode.IsGraphic(r)
//	})
//
//	cleanStdErr := strings.TrimFunc(stderr.String(), func(r rune) bool {
//		return !unicode.IsGraphic(r)
//	})
//
//	return cleanStdOut, cleanStdErr, err
//}

// RemoteCommand takes user, addr and privateKey and initiates an SSH session.
// It then runs the provided cmd and returns stdout, stderr output and error.
func RemoteCommand(user string, addr string, privateKey string, cmd string) (string, string, error) {
	log.Info().Msgf("Running command %s on %s@%s", cmd, user, addr)

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

func RetryWithExponentialBackoff(maxRetries int, task func() error) error {
	initialBackoff := time.Second
	maxBackoff := time.Minute

	for attempt := 1; ; attempt++ {
		var err error

		if err = task(); err == nil {
			return nil
		}
		log.Info().Msgf("Task failed with error %s, retrying attempt %d", err, attempt)

		if attempt >= maxRetries {
			return fmt.Errorf("task failed after %d attempts", maxRetries)
		}

		backoff := time.Duration(1<<uint(attempt)) * initialBackoff
		if backoff > maxBackoff {
			backoff = maxBackoff
		}

		time.Sleep(backoff)
	}
}
