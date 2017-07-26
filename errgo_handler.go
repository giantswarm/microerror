package microerror

import (
	"fmt"

	"github.com/juju/errgo"
)

// ErrgoHandler implements Handler interface.
type ErrgoHandler struct {
	maskFunc func(err error, allow ...func(error) bool) error
}

func NewErrgoHandler() *ErrgoHandler {
	return &ErrgoHandler{
		maskFunc: errgo.MaskFunc(errgo.Any),
	}
}

func (h *ErrgoHandler) Cause(err error) error {
	return errgo.Cause(err)
}

func (h *ErrgoHandler) Mask(err error) error {
	return h.maskFunc(err)
}

func (h *ErrgoHandler) Maskf(err error, f string, v ...interface{}) error {
	if err == nil {
		return nil
	}

	f = fmt.Sprintf("%s: %s", err.Error(), f)
	newErr := errgo.WithCausef(nil, errgo.Cause(err), f, v...)
	newErr.(*errgo.Err).SetLocation(1)

	return newErr
}
