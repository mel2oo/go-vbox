package vboxmange

import (
	"sync"

	"github.com/mel2oo/go-vbox/vboxwebsrv"
)

type VboxManage struct {
	*vboxwebsrv.VboxPortType
	managedObjectId string                //登录后的管理id
	SessionId       string                //登录后的管理id
	basicAuth       *vboxwebsrv.BasicAuth //需要认证时使用
	LockConsole     sync.Mutex
}

func NewManage(username, password, url string, tls bool) *VboxManage {
	basicAuth := &vboxwebsrv.BasicAuth{
		Login:    username,
		Password: password,
	}
	return &VboxManage{
		VboxPortType: vboxwebsrv.NewVboxPortType(url, tls, basicAuth),
		basicAuth:    basicAuth,
	}
}

func (manager *VboxManage) VboxLogon() error {
	request := vboxwebsrv.IWebsessionManagerlogon{
		Username: manager.basicAuth.Login,
		Password: manager.basicAuth.Password,
	}

	response, err := manager.IWebsessionManagerlogon(&request)
	if err != nil {
		return err // TODO: Wrap the error
	}

	manager.managedObjectId = response.Returnval
	return nil
}

func (manager *VboxManage) VboxLogoff() error {
	request := vboxwebsrv.IWebsessionManagerlogoff{
		RefIVirtualBox: manager.managedObjectId,
	}

	_, err := manager.IWebsessionManagerlogoff(&request)
	return err

}

func (manager *VboxManage) VboxCreateSession() error {
	request := vboxwebsrv.IWebsessionManagergetSessionObject{
		RefIVirtualBox: manager.managedObjectId,
	}

	response, err := manager.IWebsessionManagergetSessionObject(&request)
	if err != nil {
		return err
	}
	manager.SessionId = response.Returnval
	return nil
}

//Unlocks a machine that was previously locked for the current session.
// Calling this method is required every time a machine has been locked for a particular session using the IMachine::launchVMProcess or IMachine::lockMachine calls. Otherwise the state of the machine will be set to MachineState_Aborted on the server, and changes made to the machine settings will be lost.
// Generally, it is recommended to unlock all machines explicitly before terminating the application (regardless of the reason for the termination).

func (manager *VboxManage) VboxUnlockMachine() error {
	request := vboxwebsrv.ISessionunlockMachine{
		This: manager.SessionId,
	}

	response, err := manager.ISessionunlockMachine(&request)
	_ = response
	return err
}

func (manager *VboxManage) Release(managedObjectId string) error {
	request := vboxwebsrv.IManagedObjectRefrelease{This: managedObjectId}

	_, err := manager.IManagedObjectRefrelease(&request)
	return err
}

func (manager *VboxManage) VboxError(errorId string) (string, error) {

	request := vboxwebsrv.IVirtualBoxErrorInfogetText{This: errorId}

	response, err := manager.IVirtualBoxErrorInfogetText(&request)
	if err != nil {
		return "", err
	}
	return response.Returnval, nil
}

// func (manager *VboxManage) Vbox(managedObjectId string) error {
// 	request := vboxwebsrv.IManagedObjectRefrelease{This: managedObjectId}

// 	_, err := manager.IManagedObjectRefrelease(&request)
// 	return err
// }
