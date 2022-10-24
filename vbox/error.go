package vbox

import "errors"

var (
	ErrCommandTimeout  = errors.New("command timeout")
	ErrCommandNull     = errors.New("command can't be nil")
	ErrCommandNotFound = errors.New("command not found")
	ErrMachineState    = errors.New("machine state error")
)
