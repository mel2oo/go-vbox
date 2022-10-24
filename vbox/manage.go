package vbox

import (
	"regexp"
	"strings"
)

var (
	reVMNameUUID = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
)

type Option func(*Manage)

type Manage struct {
	cmd Command
	bin string
}

func NewManage(opts ...Option) (*Manage, error) {
	mgr, err := DefaultManage()
	if err != nil {
		return nil, err
	}

	for _, o := range opts {
		o(mgr)
	}

	return mgr, nil
}

func DefaultManage() (*Manage, error) {
	cmd, err := NewCmd()
	if err != nil {
		panic(err)
	}

	return &Manage{
		cmd: cmd,
		bin: "VBoxManage",
	}, nil
}

func WithCmd(cmd Command) Option {
	return func(m *Manage) {
		m.cmd = cmd
	}
}

func WithBin(bin string) Option {
	return func(m *Manage) {
		m.bin = bin
	}
}

// VBoxManage --version
func (m *Manage) Version() (string, error) {
	stdout, err := m.cmd.RunOutput(m.bin, "--version")
	return strings.TrimSpace(stdout), err
}

// VBoxManage registervm <filename> --password file
func (m *Manage) Register(file string) error {
	return m.cmd.Run(m.bin, file)
}

// VBoxManage unregistervm < uuid | vmname > [--delete]
func (m *Manage) Unregister(id string, delete bool) error {
	if !delete {
		return m.cmd.Run(m.bin, "unregistervm", id)
	}
	return m.cmd.Run(m.bin, "unregistervm", id, "--delete")
}
