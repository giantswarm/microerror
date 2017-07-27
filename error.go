// Package microerror provides project wide helper functions for a more convenient
// and efficient error handling.
package microerror

var (
	Default Handler = NewErrgoHandler(DefaultErrgoHandlerConfig())
)

// New returns a new error with the given error message. It is a drop-in
// replacement for errors.New from the standard library.
func New(s string) error {
	err := Default.New(s)
	err = hackErrgoCallDepth(err)
	return err
}

// Newf returns a new error with the given printf-formatted error message.
func Newf(f string, v ...interface{}) error {
	err := Default.Newf(f, v...)
	err = hackErrgoCallDepth(err)
	return err
}

// Cause returns the cause of the given error. If the cause of the err can not
// be found it returns the err itself.
//
// Cause is the usual way to diagnose errors that may have been wrapped by Mask
// or Maskf.
func Cause(err error) error {
	return Default.Cause(err)
}

// Mask is a simple error masker. Masked errors act as tracers within the
// source code. Inspecting an masked error shows where the error was passed
// through within the code base. This is gold for debugging and bug hunting.
func Mask(err error) error {
	err = Default.Mask(err)
	err = hackErrgoCallDepth(err)
	return err
}

// Maskf is like Mask. In addition to that it takes a format string and
// variadic arguments like fmt.Sprintf. The format string and variadic
// arguments are used to annotate the given errgo error.
func Maskf(err error, f string, v ...interface{}) error {
	err = Default.Maskf(err, f, v...)
	err = hackErrgoCallDepth(err)
	return err
}

type locationer interface {
	SetLocation(int)
}

func hackErrgoCallDepth(err error) error {
	if err == nil {
		return nil
	}
	_, ok := Default.(*ErrgoHandler)
	if ok {
		err.(locationer).SetLocation(2)
	}
	return err
}
