package cmd

import (
	"bytes"
	"os/exec"
)

type Command interface {
	Run(name string, args ...string) error
	RunOutput(name string, args ...string) (string, string, error)
}

type comcmd struct{}

func NewCmd() (Command, error) {
	return &comcmd{}, nil
}

func (c *comcmd) Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func (c *comcmd) RunOutput(name string, args ...string) (string, string, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
