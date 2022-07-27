package vbox

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	return newSSHCmd(user, []ssh.AuthMethod{ssh.Password(password)}, host, port)
}

func NewSSHCmdWithPrivateKey(user, privateKeyFile, host string, port int) (Command, error) {
	key, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return newSSHCmd(user, []ssh.AuthMethod{ssh.PublicKeys(signer)}, host, port)
}

func newSSHCmd(user string, auth []ssh.AuthMethod, host string, port int) (Command, error) {
	var (
		addr         string
		clientConfig *ssh.ClientConfig
	)

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

func (s *sshcmd) setSudo(sudo bool) Command {
	s.sudo = sudo
	return s
}

func (s *sshcmd) isGuest() bool {
	return s.guest
}

func (s *sshcmd) path() string {
	return s.program
}

func (s *sshcmd) run(args ...string) error {
	defer s.setSudo(false)

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
	defer s.setSudo(false)

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

func (s *sshcmd) rrunOut(cmdline string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	if Verbose {
		session.Stderr = os.Stderr
	}

	b, err := session.Output(cmdline)
	if err != nil {
		return "", err
	}

	return string(b), err
}

func (s *sshcmd) runOutErr(args ...string) (string, string, error) {
	defer s.setSudo(false)

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
