package vbox

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func manage() *Manage {
	cmd, err := NewSSHCmd("root", "Dbapp@2121", "10.20.152.15", 22, time.Second*3)
	if err != nil {
		panic(err)
	}

	mgr, err := NewManage(WithCmd(cmd))
	if err != nil {
		panic(err)
	}

	return mgr
}

func TestListMachines(t *testing.T) {
	vms, err := manage().ListMachines()
	if err != nil {
		t.Fail()
		return
	}

	data, err := json.MarshalIndent(vms, "", "  ")
	if err != nil {
		t.Fail()
		return
	}

	fmt.Println(string(data))
}

func TestGetMachine(t *testing.T) {
	vm, err := manage().GetMachine("win7_64_2")
	if err != nil {
		t.Fail()
		return
	}

	data, err := json.MarshalIndent(vm, "", "  ")
	if err != nil {
		t.Fail()
		return
	}

	fmt.Println(string(data))
}

func TestStart(t *testing.T) {
	vm, err := manage().GetMachine("win7_64_2")
	if err != nil {
		t.Fail()
		return
	}

	if err := vm.Start(); err != nil {
		t.Fail()
	}
}

func TestPause(t *testing.T) {
	vm, err := manage().GetMachine("win7_64_2")
	if err != nil {
		t.Fail()
		return
	}

	if err := vm.Pause(); err != nil {
		t.Fail()
	}
}

func TestResume(t *testing.T) {
	vm, err := manage().GetMachine("win7_64_2")
	if err != nil {
		t.Fail()
		return
	}

	if err := vm.Resume(); err != nil {
		t.Fail()
	}
}

func TestReset(t *testing.T) {
	vm, err := manage().GetMachine("win7_64_2")
	if err != nil {
		t.Fail()
		return
	}

	if err := vm.Reset(); err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestPoweroff(t *testing.T) {
	vm, err := manage().GetMachine("win7_64_2")
	if err != nil {
		t.Fail()
		return
	}

	if err := vm.Poweroff(); err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestSave(t *testing.T) {
	vm, err := manage().GetMachine("win7_64_2")
	if err != nil {
		t.Fail()
		return
	}

	if err := vm.Save(); err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestAcpiPowerButton(t *testing.T) {
	vm, err := manage().GetMachine("win7_64_2")
	if err != nil {
		t.Fail()
		return
	}

	if err := vm.AcpiPowerButton(); err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestAcpiSleepButton(t *testing.T) {
	vm, err := manage().GetMachine("win7_64_2")
	if err != nil {
		t.Fail()
		return
	}

	if err := vm.AcpiSleepButton(); err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

// snapshot
func TestListSnapshot(t *testing.T) {
	vm, err := manage().GetMachine("win7_64_2")
	if err != nil {
		t.Fail()
		return
	}

	out, err := vm.ListSnapshot()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fail()
		return
	}

	fmt.Println(string(data))
}

func TestBindCPU(t *testing.T) {
	vm, err := manage().GetMachine("win7_64_2")
	if err != nil {
		t.Fail()
		return
	}

	if err := vm.BindCpu("12"); err != nil {
		t.Fail()
		fmt.Println(err)
	}
}
