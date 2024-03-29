package vbox

import (
	"regexp"
)

var manage Command

// Command is the mock-able interface to run VirtualBox commands
// such as VBoxManage (host side) or VBoxControl (guest side)
type Command interface {
	setSudo(bool) Command
	isGuest() bool
	path() string
	run(args ...string) error
	runOut(args ...string) (string, error)
	rrunOut(cmdline string) (string, error)
	runOutErr(args ...string) (string, string, error)
}

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)
)
