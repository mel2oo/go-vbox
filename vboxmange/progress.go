package vboxmange

import "github.com/mel2oo/go-vbox/vboxwebsrv"

type Progress struct {
	*VboxManage

	ProgressId string
}

func (p *Progress) ProgressWaitForCompletion(timeout int32) error {
	request := vboxwebsrv.IProgresswaitForCompletion{This: p.managedObjectId}
	request.Timeout = timeout

	_, err := p.IProgresswaitForCompletion(&request)
	if err != nil {
		return err // TODO: Wrap the error
	}
	p.ProgressRelease()

	// TODO: See if we need to do anything with the response
	return nil
}

func (p *Progress) ProgressGetPercent() (uint32, error) {
	request := vboxwebsrv.IProgressgetPercent{This: p.managedObjectId}

	response, err := p.IProgressgetPercent(&request)
	if err != nil {
		return 0, err // TODO: Wrap the error
	}
	p.ProgressRelease()

	// TODO: See if we need to do anything with the response
	return response.Returnval, nil
}

func (p *Progress) ProgressRelease() error {

	return p.Release(p.ProgressId)
}
