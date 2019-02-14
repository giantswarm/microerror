package microerror

import (
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
		callDepth: config.CallDepth + 2, // +2 for ErrorHandler wrapping methods
	}
}

func (h *ErrorHandler) Cause(err error) error {
	e, ok := err.(*maskedError)
	if ok {
		return e.Cause.error
	}

	return nil
}

func (h *ErrorHandler) Mask(err error) error {
	return h.mask(err, "")
}

func (h *ErrorHandler) Maskf(err *Error, f string, v ...interface{}) error {
	msg := fmt.Sprintf(f, v...)

	return h.mask(err, msg)
}

func (h *ErrorHandler) mask(err error, msg string) error {
	if err == nil {
		return nil
	}

	e := newMaskedError(err)

	if len(msg) > 0 {
		e.Cause.Annotation = msg
	}

	_, f, l, _ := runtime.Caller(h.callDepth)
	s := maskedErrorStack{
		File: f,
		Line: l,
	}

	e.Stack = append(e.Stack, s)

	return e
}
