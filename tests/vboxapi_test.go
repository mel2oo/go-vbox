package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/mel2oo/go-vbox/vboxapi"
	"github.com/mel2oo/go-vbox/vboxmange"
)

// vboxmanage setproperty websrvauthlibrary null
//vboxmanage setproperty vrdeauthlibrary “VBoxAuthSimple”

func Test_VboxApi(t *testing.T) {

	vbox := vboxapi.New("", "", "http://10.20.152.15:9999", false, "win7")
	if vbox == nil {
		fmt.Printf("can't open vbox api")
	}
	if err := vbox.Logon(); err != nil {
		fmt.Printf("Unable to log on to vboxweb: %v\n", err)
	}

	machines, err := vbox.GetMachines()
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

func Test_VboxApiManage(t *testing.T) {

	vmname := "win7_64_1"
	vbox := vboxmange.NewManage("", "", "http://10.20.152.15:9999", false)

	if err := vbox.VboxLogon(); err != nil {
		fmt.Printf("Unable to log on to vboxweb: %v\n", err)
	}

	// machines, err := vbox.VboxGetMachines()
	// if err != nil {
	// 	fmt.Printf("can;t get machine:%v", err)
	// } else {
	// 	fmt.Printf("machine:%v", machines)
	// }
	// for _, m_id := range machines {
	// 	id, _ := m_id.MachineGetID()
	// 	name, _ := m_id.MachineGetName()
	// 	fmt.Printf("id:%s,name:%s\n", id, name)

	// }
	win7_1, err := vbox.VboxGetMachine(vmname)
	if err != nil {
		fmt.Printf("can;t get machine:%v\n", err)
	} else {
		fmt.Printf("machine:%v\n", win7_1)
	}
	win7_Status, err := win7_1.MachineGetStatus()
	if err != nil {
		fmt.Printf("can;t get machine:%v\n", err)
	} else {
		fmt.Printf("machine:%v\n", win7_Status)
	}

	err = win7_1.MachineSnapRestoreSnap("control")
	if err != nil {
		fmt.Printf("can't restore snap control :%v\n", err)
	} else {
		fmt.Printf("success restore control\n")
	}

	err = win7_1.VmStart()
	if err != nil {
		// text := win7_1.VboxError(err.Error())
		fmt.Printf("can't start machine:%v\n", err)
	} else {
		fmt.Printf("success start machine:%v\n", vmname)
	}
	time.Sleep(time.Second * 10)
	err = win7_1.VmStop()
	if err != nil {
		fmt.Printf("can't stop machine:%v\n", err)
	} else {
		fmt.Printf("success stop machine:%v\n", vmname)
	}

}
