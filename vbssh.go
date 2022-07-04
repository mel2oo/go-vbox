package vbox

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type sshcmd struct {
	client  *ssh.Client
	program string
	sudoer  bool
	sudo    bool
	guest   bool
}

func NewSSHCmd(user, password, host string, port int) (Command, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}

	manage = &sshcmd{
		client:  client,
		program: "VBoxManage",
		sudoer:  true,
		sudo:    true,
		guest:   false}

	return manage, nil
}

func (s *sshcmd) setOpts(opts ...option) Command {
	var cmd Command = s
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func (s *sshcmd) isGuest() bool {
	return s.guest
}

func (s *sshcmd) path() string {
	return s.program
}

func (s *sshcmd) run(args ...string) error {
	defer s.setOpts(sudo(false))

	session, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if Verbose {
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
	}

	cmdline := fmt.Sprintf("%s %s", s.program, strings.Join(args, " "))
	return session.Run(cmdline)
}

func (s *sshcmd) runOut(args ...string) (string, error) {
	defer s.setOpts(sudo(false))

	session, err := s.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	if Verbose {
		session.Stderr = os.Stderr
	}

	cmdline := fmt.Sprintf("%s %s", s.program, strings.Join(args, " "))
	b, err := session.Output(cmdline)
	if err != nil {
		return "", err
	}

	return string(b), err
}

func (s *sshcmd) runOutErr(args ...string) (string, string, error) {
	defer s.setOpts(sudo(false))

	session, err := s.client.NewSession()
	if err != nil {
		return "", "", err
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	cmdline := fmt.Sprintf("%s %s", s.program, strings.Join(args, " "))
	if err := session.Run(cmdline); err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), err
}
