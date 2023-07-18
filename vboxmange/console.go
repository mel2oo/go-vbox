package vboxmange

import "github.com/mel2oo/go-vbox/vboxwebsrv"

type Console struct {
	*VboxManage

	ConsoleID string
}

func (manager *VboxManage) GetConsole() (*Console, error) {

	request := vboxwebsrv.ISessiongetConsole{This: manager.SessionId}

	response, err := manager.ISessiongetConsole(&request)
	if err != nil {
		return nil, err // TODO: Wrap the error
	}

	// TODO: See if we need to do anything with the response
	return &Console{manager, response.Returnval}, nil
}

// PowerDown starts forcibly powering off the controlled VM.
// It returns a Progress and any error encountered.
func (c *Console) PowerDown() error {
	request := vboxwebsrv.IConsolepowerDown{This: c.ConsoleID}

	response, err := c.IConsolepowerDown(&request)
	if err != nil {
		return err //   Wrap the error
	}
	progress := &Progress{c.VboxManage, response.Returnval}
	return progress.ProgressWaitForCompletion(-1)
}

// PowerUp starts powering on the controlled VM.
// It returns a Progress and any error encountered.
func (c *Console) PowerUp() error {

	request := vboxwebsrv.IConsolepowerUp{This: c.ConsoleID}

	response, err := c.IConsolepowerUp(&request)
	if err != nil {
		return err //   Wrap the error
	}
	progress := &Progress{c.VboxManage, response.Returnval}
	return progress.ProgressWaitForCompletion(-1)
}
