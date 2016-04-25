package Utiles

import "runtime/debug"

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

var (
	OUTPUT_STACK_ON_ERROR = false
)

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
