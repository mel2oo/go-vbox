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

func (manager *VMachine) VmStart(vmname string) error {
	vmstatus, err := manager.VmStatus(vmname)
	if err != nil {
		return err
	}
	switch vmstatus {
	case string(vboxwebsrv.MachineStatePoweredOff):
		break
	default:
		return fmt.Errorf("machine not is poweroff,current status is %s", vmstatus)
	}
	machineid, err := manager.GetMachine(vmname)
	if err != nil {
		return err
	}
	request := vboxwebsrv.IConsolepowerUp{This: machineid}
	response, err := manager.IConsolepowerUp(&request)
	if err != nil {
		return err // TODO: Wrap the error
	}
	waitrequest := vboxwebsrv.IProgresswaitForCompletion{This: response.Returnval, Timeout: -1}

	_, err = manager.IProgresswaitForCompletion(&waitrequest)
	if err != nil {
		return err // TODO: Wrap the error
	}
	// TODO: See if we need to do anything with the response
	return nil

}
func (manager *VMachine) VmStop(vmname string) error {
	vmstatus, err := manager.VmStatus(vmname)
	if err != nil {
		return err
	}
	switch vmstatus {
	case string(vboxwebsrv.MachineStateRunning):
		break
	default:
		return fmt.Errorf("machine not is running,current status is %s", vmstatus)
	}
	machineid, err := manager.GetMachine(vmname)
	if err != nil {
		return err
	}
	request := vboxwebsrv.IConsolepowerDown{This: machineid}
	response, err := manager.IConsolepowerDown(&request)
	if err != nil {
		return err // TODO: Wrap the error
	}
	waitrequest := vboxwebsrv.IProgresswaitForCompletion{This: response.Returnval, Timeout: -1}

	_, err = manager.IProgresswaitForCompletion(&waitrequest)
	if err != nil {
		return err // TODO: Wrap the error
	}
	// TODO: See if we need to do anything with the response
	return nil

}

func (manager *VMachine) VmTakeSnap(snapname string,desc string,pause bool) (string, error) {
	status, err := manager.MachineGetStatus()
	switch status{
		case string(vboxwebsrv.MachineStateRunning):
	}

	request := vboxwebsrv.IMachinegetName{This: vm_manager_id}

	response, err := manager.IMachinegetName(&request)
	if err != nil {
		return "", err // TODO: Wrap the error
	}
	// TODO: See if we need to do anything with the response
	return response.Returnval, nil
	return "", nil
}
func (manager *VMachine) VmDeleteSnap(vmname string) (string, error) {

	return "", nil
}
func (manager *VMachine) VmSnapRestoreSnap(vmname string) (string, error) {
	return "", nil
}
