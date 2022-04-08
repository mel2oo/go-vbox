package vbox

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	manage Command
	client *ssh.Client
)

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)
)

// Manage returns the Command to run VBoxManage/VBoxControl.
func Manage() Command {
	if manage != nil {
		return manage
	}

	sudoer, err := isSudoer()
	if err != nil {
		Debug("Error getting sudoer status: '%v'", err)
	}

	if vbprog, err := lookupVBoxProgram("VBoxManage"); err == nil {
		manage = command{program: vbprog, sudoer: sudoer, guest: false, remote: false}
	} else if vbprog, err := lookupVBoxProgram("VBoxControl"); err == nil {
		manage = command{program: vbprog, sudoer: sudoer, guest: true, remote: false}
	} else if client != nil {
		manage = command{program: "VBoxManage", sudoer: sudoer, guest: false, remote: true}
	} else {
		manage = command{program: "false", sudoer: false, guest: false, remote: false}
	}
	Debug("manage: '%+v'", manage)

	return manage
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

func Connect(user, password, host string, port int) error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		err          error
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

	client, err = ssh.Dial("tcp", addr, clientConfig)

	return err
}
