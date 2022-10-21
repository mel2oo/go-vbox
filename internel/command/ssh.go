package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var ErrSshTimeout = errors.New("ssh cmd timeout")

type sshcmd struct {
	client  *ssh.Client
	timeout time.Duration
}

func NewSSHCmd(user, password, host string, port int,
	timeout time.Duration) (Command, error) {
	client, err := newClient(user, []ssh.AuthMethod{ssh.Password(password)}, host, port)
	if err != nil {
		return nil, err
	}

	return &sshcmd{
		client:  client,
		timeout: timeout,
	}, nil
}

func NewSSHCmdWithPrivateKey(user, keyfile, host string, port int,
	timeout time.Duration) (Command, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	client, err := newClient(user, []ssh.AuthMethod{ssh.PublicKeys(signer)}, host, port)
	if err != nil {
		return nil, err
	}

	return &sshcmd{
		client:  client,
		timeout: timeout,
	}, nil
}

func newClient(user string, auth []ssh.AuthMethod, host string, port int) (*ssh.Client, error) {
	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
}

func (s *sshcmd) Run(name string, args ...string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if err := session.Start(
		fmt.Sprintf("%s %s", name, strings.Join(args, " "))); err != nil {
		return err
	}

	exit := make(chan struct{}, 1)
	go func() {
		session.Wait()
		exit <- struct{}{}
	}()

	select {
	case <-exit:
	case <-time.After(s.timeout):
		return ErrSshTimeout
	}

	return nil
}

func (s *sshcmd) RunOutput(name string, args ...string) (string, string, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
		exit   = make(chan struct{}, 1)
	)

	session, err := s.client.NewSession()
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	defer session.Close()

	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Start(fmt.Sprintf("%s %s", name, strings.Join(args, " ")))
	go func() {
		session.Wait()
		exit <- struct{}{}
	}()

	select {
	case <-exit:
	case <-time.After(s.timeout):
		return stdout.String(), stderr.String(), ErrSshTimeout
	}

	return stdout.String(), stderr.String(), err
}
