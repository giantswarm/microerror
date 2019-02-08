package microerror

import (
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

// Mask wraps an error to record its stack. All errors should be masked along
// the path of execution. At the root cause of a situation where an error occurs
// at runtime, a typed error can be returned, while masking it. This would look
// like the following.
//
//     return microerror.Mask(executionFailedError)
//
// The example above will result in detailed information where the actual
// misconfiguration happened. When debugging the error, it is then very easy to
// fix the actual problem, because the stack tracks the originating source code
// position. When forwarding errors received from arbitrary sources, errors
// should be wrapped as well. The following example shows how.
//
//     return microerror.Mask(err)
//
// Note that it is not necessary to annotate errors that are only received and
// forwarded to the next caller.
//
func (h *ErrorHandler) Mask(err error) error {
	if err == nil {
		return nil
	}

	e := newDefaultError()
	e.Cause = h.Cause(err)
	w, ok := err.(*Error)
	if ok {
		if w.Desc != "" {
			e.Desc = w.Desc
		}
		if w.Docs != "" {
			e.Docs = w.Docs
		}
		if w.Kind != "" {
			e.Kind = w.Kind
		}
		e.Stack = append(e.Stack, w.Stack...)
	}

	_, f, l, _ := runtime.Caller(h.callDepth)
	s := Stack{
		File: f,
		Line: l,
	}

	e.Stack = append(e.Stack, s)

	return e
}
