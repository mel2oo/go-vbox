package vboxmange

import (
	"fmt"

	"github.com/mel2oo/go-vbox/vboxwebsrv"
)

type VMachine struct {
	*VboxManage
	MachineId string
}

func (manager *VboxManage) VboxGetMachine(vmname string) (*VMachine, error) {
	request := vboxwebsrv.IVirtualBoxfindMachine{This: manager.managedObjectId, NameOrId: vmname}
	response, err := manager.IVirtualBoxfindMachine(&request)
	if err != nil {
		return nil, err // TODO: Wrap the error
	}
	return &VMachine{manager, response.Returnval}, nil
}

// 返回的是内部管理id 如果需要获取到对应的名字 需要再转换
func (manager *VboxManage) VboxGetMachines() ([]*VMachine, error) {

	request := vboxwebsrv.IVirtualBoxgetMachines{This: manager.managedObjectId}

	response, err := manager.IVirtualBoxgetMachines(&request)
	if err != nil {
		return nil, err // TODO: Wrap the error
	}

	machines := make([]*VMachine, len(response.Returnval))
	for i, machineid := range response.Returnval {
		machines[i] = &VMachine{manager, machineid}
	}

	return machines, nil
}

func (manager *VMachine) MachineGetStatus() (string, error) {
	request := vboxwebsrv.IMachinegetState{This: manager.MachineId}
	response, err := manager.IMachinegetState(&request)
	if err != nil {
		return "", err // TODO: Wrap the error
	}
	return response.Returnval, nil
}

func (manager *VMachine) MachineGetID() (string, error) {
	request := vboxwebsrv.IMachinegetId{This: manager.MachineId}

	response, err := manager.IMachinegetId(&request)
	if err != nil {
		return "", err // TODO: Wrap the error
	}

	// TODO: See if we need to do anything with the response
	return response.Returnval, nil

}
func (manager *VMachine) MachineGetName() (string, error) {
	request := vboxwebsrv.IMachinegetName{This: manager.MachineId}
	response, err := manager.IMachinegetName(&request)
	if err != nil {
		return "", err // TODO: Wrap the error
	}
	// TODO: See if we need to do anything with the response
	return response.Returnval, nil
}

// return snapshot id
func (manager *VMachine) MachineTakeSnap(snapname string, desc string, pause bool) (string, error) {
	status, err := manager.MachineGetStatus()
	switch status {
	case string(vboxwebsrv.MachineStatePoweredOff):
	case string(vboxwebsrv.MachineStateSaved):
	case string(vboxwebsrv.MachineStateAborted):
	case string(vboxwebsrv.MachineStateRunning):
	case string(vboxwebsrv.MachineStatePaused):
	default:
		return "", fmt.Errorf("machine not is poweroff,current status is %s", status)
	}

	request := vboxwebsrv.IMachinetakeSnapshot{This: manager.MachineId, Name: snapname, Description: desc, Pause: pause}

	response, err := manager.IMachinetakeSnapshot(&request)
	if err != nil {
		return "", err // TODO: Wrap the error
	}
	progress := Progress{manager.VboxManage, response.Returnval}
	err = progress.ProgressWaitForCompletion(-1)
	if err != nil {
		return "", err
	}

	// TODO: See if we need to do anything with the response
	return response.Id, nil
}

func (manager *VMachine) MachineDeleteSnap(snapname string) error {
	snapid, err := manager.MachineGetSnapShot(snapname)
	if err != nil {
		return err
	}

	request := vboxwebsrv.IMachinedeleteSnapshot{This: manager.MachineId, Id: snapid.SnapshotId}
	response, err := manager.IMachinedeleteSnapshot(&request)
	if err != nil {
		return err // TODO: Wrap the error
	}

	progress := Progress{manager.VboxManage, response.Returnval}
	err = progress.ProgressWaitForCompletion(-1)
	if err != nil {
		return err
	}
	return nil
}

func (manager *VMachine) MachineSnapRestoreSnap(snapname string) error {
	snapid, err := manager.MachineGetSnapShot(snapname)
	if err != nil {
		return err
	}
	request := vboxwebsrv.IMachinerestoreSnapshot{This: manager.MachineId, Snapshot: snapid.SnapshotId}
	response, err := manager.IMachinerestoreSnapshot(&request)
	if err != nil {
		return err // TODO: Wrap the error
	}
	progress := Progress{manager.VboxManage, response.Returnval}
	err = progress.ProgressWaitForCompletion(-1)
	if err != nil {
		return err
	}
	return nil
}

func (manager *VMachine) MachineLock() error {
	// This     string    `xml:"_this,omitempty"`
	// Session  string    `xml:"session,omitempty"`
	// LockType *LockType `xml:"lockType,omitempty"`
	manager.VboxUnlockMachine()
	// lockType := vboxwebsrv.LockTypeShared
	request := vboxwebsrv.IMachinelockMachine{This: manager.MachineId, Session: manager.SessionId, LockType: vboxwebsrv.LockTypeVM}
	_, err := manager.IMachinelockMachine(&request)
	if err != nil {
		return err // TODO: Wrap the error
	}
	// TODO: See if we need to do anything with the response
	return nil
}

func (manager *VMachine) VmStart() error {

	vmstatus, err := manager.MachineGetStatus()

	switch vmstatus {
	case string(vboxwebsrv.MachineStateSaved):
	case string(vboxwebsrv.MachineStatePoweredOff):
		break
	default:
		return fmt.Errorf("machine not is poweroff,current status is %s", vmstatus)
	}

	err = manager.MachineLock()
	if err != nil {
		return err
	}
	defer manager.VboxUnlockMachine()
	VmConsole, err := manager.GetConsole()
	if err != nil {
		return err
	}
	return VmConsole.PowerUp()

}
func (manager *VMachine) VmStop() error {

	vmstatus, err := manager.MachineGetStatus()

	switch vmstatus {
	case string(vboxwebsrv.MachineStatePoweredOff):
		break
	default:
		return fmt.Errorf("machine not is poweroff,current status is %s", vmstatus)
	}

	err = manager.MachineLock()
	if err != nil {
		return err
	}
	defer manager.VboxUnlockMachine()
	VmConsole, err := manager.GetConsole()
	if err != nil {
		return err
	}
	return VmConsole.PowerDown()

}
