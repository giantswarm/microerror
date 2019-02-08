// Package microerror provides project wide helper functions for a more
// convenient and efficient error handling.
package microerror

var (
	handler Handler = NewErrorHandler(ErrorHandlerConfig{
		CallDepth: 1,
	})
)

// Cause returns the cause of the given error. If the cause of the err can not
// be found it returns the err itself. Cause is the usual way to diagnose errors
// that may have been wrapped by Mask or Maskf.
func Cause(err error) error {
	return handler.Cause(err)
}

// Mask is a simple error masker. Masked errors act as tracers within the
// source code. Inspecting an masked error shows where the error was passed
// through within the code base. This is gold for debugging and bug hunting.
func Mask(err error) error {
	return handler.Mask(err)
}
