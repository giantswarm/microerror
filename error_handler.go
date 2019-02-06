package microerror

import (
	"errors"
	"fmt"
	"runtime"
)

type ErrorHandlerConfig struct {
	// CallDepth is useful when creating a wrapper for ErrorHandler. Its
	// value is used to push stack location and skip wrapping function
	// location as an origin. The default value is 0.
	CallDepth int
}

// ErrorHandler implements the Handler interface.
type ErrorHandler struct {
	callDepth int
}

func NewErrorHandler(config ErrorHandlerConfig) *ErrorHandler {
	return &ErrorHandler{
		callDepth: config.CallDepth + 1, // +1 for ErrorHandler wrapping methods
	}
}

func (h *ErrorHandler) New(s string) error {
	return errors.New(s)
}

func (h *ErrorHandler) Newf(f string, v ...interface{}) error {
	return fmt.Errorf(f, v...)
}

func (h *ErrorHandler) Cause(err error) error {
	e, ok := err.(*Error)
	if !ok {
		return err
	}

	if e.Cause != nil {
		return e.Cause
	}

	return e
}

func (h *ErrorHandler) Mask(err error) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*Error)
	if !ok {
		e = newDefaultError()
		e.Cause = h.Cause(err)
	}

	_, f, l, _ := runtime.Caller(h.callDepth)
	s := Stack{
		File:    f,
		Line:    l,
		Message: "",
	}

	e.Stack = append(e.Stack, s)

	return e
}

func (h *ErrorHandler) Maskf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*Error)
	if !ok {
		e = newDefaultError()
		e.Cause = h.Cause(err)
	}

	_, f, l, _ := runtime.Caller(h.callDepth)
	s := Stack{
		File:    f,
		Line:    l,
		Message: fmt.Sprintf(format, args...),
	}

	e.Stack = append(e.Stack, s)

	return e
}
