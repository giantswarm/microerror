package microerror

import (
	"github.com/juju/errgo"
)

// ErrgoHandler implements Handler interface.
type ErrgoHandler struct {
	callDepth int
	maskFunc  func(err error, allow ...func(error) bool) error
}

func NewErrgoHandler() *ErrgoHandler {
	return NewErrgoHandlerCallDepth(0)
}

// NewErrgoHandlerCallDepth is useful when creating a wrapper for ErrgoHandler.
// The callDepth parameter is used to push stack location and skip wrapping
// function location as an origin. The default value is 0.
func NewErrgoHandlerCallDepth(callDepth int) *ErrgoHandler {
	return &ErrgoHandler{
		callDepth: callDepth + 1, // +1 for ErrgoHandler wrapping methods
		maskFunc:  errgo.MaskFunc(errgo.Any),
	}
}

func (h *ErrgoHandler) Cause(err error) error {
	return errgo.Cause(err)
}

func (h *ErrgoHandler) Mask(err error) error {
	if err == nil {
		return nil
	}

	newErr := h.maskFunc(err)
	newErr.(*errgo.Err).SetLocation(h.callDepth)
	return newErr
}

func (h *ErrgoHandler) Maskf(err error, f string, v ...interface{}) error {
	if err == nil {
		return nil
	}

	newErr := errgo.WithCausef(err, errgo.Cause(err), f, v...)
	newErr.(*errgo.Err).SetLocation(h.callDepth)
	return newErr
}
