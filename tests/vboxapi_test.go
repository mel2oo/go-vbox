package tests

import (
	"fmt"
	"testing"

	"github.com/mel2oo/go-vbox/vboxapi"
)

// vboxmanage setproperty websrvauthlibrary null
//vboxmanage setproperty vrdeauthlibrary “VBoxAuthSimple”

func Test_VboxApi(t *testing.T) {

	vboxmanage := vboxapi.New("", "", "http://10.20.152.15:9999", false, "win7")
	if vboxmanage == nil {
		fmt.Printf("can't open vbox api")
	}
	if err := vboxmanage.Logon(); err != nil {
		fmt.Printf("Unable to log on to vboxweb: %v\n", err)
	}

	machines, err := vboxmanage.GetMachines()
	if err != nil {
		fmt.Printf("can;t get machine:%v", err)
	} else {
		fmt.Printf("machine:%v", machines)
	}
	for _, machine := range machines {
		machine.Refresh()
		id, _ := machine.GetID()
		name, _ := machine.GetName()

		fmt.Printf("id:%s,name:%s\n", id, name)

	}

}
