package vboxmange

import (
	"github.com/mel2oo/go-vbox/vboxwebsrv"
)

type VSnapshot struct {
	*VboxManage
	SnapshotId string
}

func (manager *VMachine) MachineGetSnapShot(snapname string) (*VSnapshot, error) {
	request := vboxwebsrv.IMachinefindSnapshot{This: manager.managedObjectId, NameOrId: snapname}
	response, err := manager.IMachinefindSnapshot(&request)
	if err != nil {
		return nil, err // TODO: Wrap the error
	}
	return &VSnapshot{manager, response.Returnval}, nil
}

// 返回的是内部管理id 如果需要获取到对应的名字 需要再转换
func (manager *VMachine) VboxGetMachines() ([]*VMachine, error) {

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
