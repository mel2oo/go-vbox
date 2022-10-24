package vbox

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type Command interface {
	Run(name string, args ...string) error
	RunOutput(name string, args ...string) (string, error)
}

// local commandline
type comcmd struct {
	mutex *sync.Mutex
}

func NewCmd() (Command, error) {
	return &comcmd{mutex: new(sync.Mutex)}, nil
}

func (c *comcmd) Run(name string, args ...string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var (
		stderr bytes.Buffer
	)

	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr

	err := cmd.Run()
	if len(stderr.String()) > 0 {
		return errors.New(stderr.String())
	}

	return err
}

func (c *comcmd) RunOutput(name string, args ...string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if len(stderr.String()) > 0 {
		return stdout.String(), errors.New(stderr.String())
	}

	return stdout.String(), err
}

// ssh command line
type sshcmd struct {
	mutex   *sync.Mutex
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
		mutex:   new(sync.Mutex),
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
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var (
		stderr bytes.Buffer
		exit   = make(chan struct{}, 1)
	)

	session, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stderr = &stderr

	if err := session.Start(
		fmt.Sprintf("%s %s", name, strings.Join(args, " "))); err != nil {
		return err
	}

	go func() {
		session.Wait()
		exit <- struct{}{}
	}()

	select {
	case <-exit:
	case <-time.After(s.timeout):
		return ErrCommandTimeout
	}

	if len(stderr.String()) > 0 {
		return errors.New(stderr.String())
	}

	return err
}

func (s *sshcmd) RunOutput(name string, args ...string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
		exit   = make(chan struct{}, 1)
	)

	session, err := s.client.NewSession()
	if err != nil {
		return stdout.String(), err
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
		return stdout.String(), ErrCommandTimeout
	}

	if len(stderr.String()) > 0 {
		return stdout.String(), errors.New(stderr.String())
	}

	return stdout.String(), err
}
