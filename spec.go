package microerror

type Handler interface {
	// Cause returns the cause of the given error. If the cause of the err can not
	// be found it returns the err itself.
	//
	// Cause is the usual way to diagnose errors that may have been wrapped by Mask
	// or Maskf.
	Cause(err error) error
	// Mask is a simple error masker. Masked errors act as tracers within the
	// source code. Inspecting an masked error shows where the error was passed
	// through within the code base. This is gold for debugging and bug hunting.
	Mask(err error) error
}
