package vbox

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
)

var (
	// Verbose toggles the library in verbose execution mode.
	Verbose bool
	// ErrMachineExist holds the error message when the machine already exists.
	ErrMachineExist = errors.New("machine already exists")
	// ErrMachineNotExist holds the error message when the machine does not exist.
	ErrMachineNotExist = errors.New("machine does not exist")
	// ErrCommandNotFound holds the error message when the VBoxManage commands was not found.
	ErrCommandNotFound = errors.New("command not found")
)

type command struct {
	program string
	sudoer  bool // Is current user a sudoer?
	sudo    bool // Is current command expected to be run under sudo?
	guest   bool
}

func NewCommand() (Command, error) {
	sudoer, err := isSudoer()
	if err != nil {
		return nil, err
	}

	if vbprog, err := lookupVBoxProgram("VBoxManage"); err == nil {
		manage = command{program: vbprog, sudoer: sudoer, guest: false}
	} else if vbprog, err := lookupVBoxProgram("VBoxControl"); err == nil {
		manage = command{program: vbprog, sudoer: sudoer, guest: true}
	} else {
		manage = command{program: "false", sudoer: false, guest: false}
	}

	return manage, nil
}

func lookupVBoxProgram(vbprog string) (string, error) {
	if runtime.GOOS == osWindows {
		if p := os.Getenv("VBOX_INSTALL_PATH"); p != "" {
			vbprog = filepath.Join(p, vbprog+".exe")
		} else {
			vbprog = filepath.Join("C:\\", "Program Files", "Oracle", "VirtualBox", vbprog+".exe")
		}
	}

	return exec.LookPath(vbprog)
}

func isSudoer() (bool, error) {
	me, err := user.Current()
	if err != nil {
		return false, err
	}
	Debug("User: '%+v'", me)
	if groupIDs, err := me.GroupIds(); runtime.GOOS == "linux" {
		if err != nil {
			return false, err
		}
		Debug("groupIDs: '%+v'", groupIDs)
		for _, groupID := range groupIDs {
			group, err := user.LookupGroupId(groupID)
			if err != nil {
				return false, err
			}
			Debug("group: '%+v'", group)
			if group.Name == "sudo" {
				return true, nil
			}
		}
	}
	return false, nil
}

func (vbcmd command) setOpts(opts ...option) Command {
	var cmd Command = &vbcmd
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func (vbcmd command) isGuest() bool {
	return vbcmd.guest
}

func (vbcmd command) path() string {
	return vbcmd.program
}

func (vbcmd command) prepare(args []string) *exec.Cmd {
	program := vbcmd.program
	argv := []string{}
	Debug("Command: '%+v', runtime.GOOS: '%s'", vbcmd, runtime.GOOS)
	if vbcmd.sudoer && vbcmd.sudo && runtime.GOOS != osWindows {
		program = "sudo"
		argv = append(argv, vbcmd.program)
	}
	argv = append(argv, args...)
	Debug("executing: %v %v", program, argv)
	return exec.Command(program, argv...) // #nosec
}

func (vbcmd command) run(args ...string) error {
	defer vbcmd.setOpts(sudo(false))

	cmd := vbcmd.prepare(args)
	if Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			return ErrCommandNotFound
		}
		return err
	}

	return nil
}

func (vbcmd command) runOut(args ...string) (string, error) {
	defer vbcmd.setOpts(sudo(false))

	cmd := vbcmd.prepare(args)
	if Verbose {
		cmd.Stderr = os.Stderr
	}

	b, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrCommandNotFound
		}
	}
	return string(b), err
}

func (vbcmd command) runOutErr(args ...string) (string, string, error) {
	defer vbcmd.setOpts(sudo(false))

	cmd := vbcmd.prepare(args)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrCommandNotFound
		}
	}
	return stdout.String(), stderr.String(), err
}
