package Utiles

import (
	"runtime/debug"
	"errors"
)

var (
	OUTPUT_STACK_ON_ERROR = false
)

type PanicError struct {
	error
	stack []byte
}

func (p *PanicError) Error() string {
	if p.stack == nil {
		return p.error.Error()
	}
	return p.error.Error() + "\n" + (string)(p.stack)
}

func NewPanicError(err error) error {
	if !OUTPUT_STACK_ON_ERROR {
		return &PanicError{
			err,
			nil,
		}
	}
	return &PanicError{
		err,
		debug.Stack(),
	}
}

func NewError(err string) error {
	if !OUTPUT_STACK_ON_ERROR {
		return &PanicError{
			errors.New(err),
			nil,
		}
	}
	return &PanicError{
		errors.New(err),
		debug.Stack(),
	}
}