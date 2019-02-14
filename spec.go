package microerror

type Handler interface {
	// Cause returns the cause of the given error. If the cause of the err can not
	// be found it returns the err itself.
	//
	// Cause is the usual way to diagnose errors that may have been wrapped by Mask
	// or Maskf.
	Cause(err error) error
	// Mask wraps an error to record its stack. All errors should be masked
	// along the path of execution. At the root cause of a situation where
	// an error occurs at runtime, a typed error can be returned, while
	// masking it. This would look like the following.
	//
	//     return microerror.Mask(executionFailedError)
	//
	// The example above will result in detailed information where the
	// actual misconfiguration happened. When debugging the error, it is
	// then very easy to fix the actual problem, because the stack tracks
	// the originating source code position. When forwarding errors
	// received from arbitrary sources, errors should be wrapped as well.
	// The following example shows how.
	//
	//     return microerror.Mask(err)
	//
	// Note that it is not necessary to annotate errors that are only
	// received and forwarded to the next caller.
	Mask(err error) error
	// Maskf is like Mask but it works only with *Error. In addition to
	// that it takes a format string and variadic arguments like
	// fmt.Sprintf. The format string and variadic arguments are used to
	// annotate the given root error. Maskf should be used only when the error is created.
	Maskf(err *Error, f string, v ...interface{}) error
}
