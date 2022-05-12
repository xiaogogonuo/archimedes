package errors

import "github.com/pkg/errors"

func NewError(message string) error {
	return errors.New(message)
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

func A() error {
	return NewError("A error")
}

func B() error {
	e := A()
	if e != nil {
		return NewError("B call a error")
		//return Wrap(e, "B call a error")
	}
	return nil
}