package vbox

import (
	"encoding/json"
	"fmt"
	"testing"
)

func manage() *Manage {
	mgr, err := NewManage()
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
