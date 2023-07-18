package vboxmange

import (
	"github.com/mel2oo/go-vbox/vboxwebsrv"
)

type VSnapshot struct {
	*VMachine
	SnapshotId string
}

func (manager *VMachine) MachineGetSnapShot(snapname string) (*VSnapshot, error) {
	request := vboxwebsrv.IMachinefindSnapshot{This: manager.MachineId, NameOrId: snapname}
	response, err := manager.IMachinefindSnapshot(&request)
	if err != nil {
		return nil, err // TODO: Wrap the error
	}
	return &VSnapshot{manager, response.Returnval}, nil
}

// 返回的是内部管理id 如果需要获取到对应的名字 需要再转换
func (manager *VSnapshot) SnapShotGetOnline() (bool, error) {

	request := vboxwebsrv.ISnapshotgetOnline{This: manager.SnapshotId}

	response, err := manager.ISnapshotgetOnline(&request)
	if err != nil {
		return false, err // TODO: Wrap the error
	}
	return response.Returnval, nil
}
