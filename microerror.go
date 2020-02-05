package microerror

import (
	"fmt"
	"runtime"

	"errors"
)

// Cause is here only for backward compatibility purposes and should not be used.
//
// NOTE: Use errors.Is/errors.As instead.
func Cause(err error) error {
	// If type of err is Error then this is the cause. This also covers all
	// calls that initiated with Maskf because Maskf takes only Error type.
	var microErr Error
	if errors.As(err, &microErr) {
		return microErr
	}

	// Now this is known that the masking was initiated with Mask so unwrap
	// all stackedError and return what's unwrapped from the one at the
	// bottom of the stack.
	var stackedErr stackedError
	for errors.As(err, &stackedErr) {
		err = stackedErr.Unwrap()
	}

	return err
}

func Maskf(err Error, f string, v ...interface{}) error {
	annotatedErr := annotatedError{
		annotation: fmt.Sprintf(f, v...),
		underlying: err,
	}

	return mask(annotatedErr)
}

func Mask(err error) error {
	if err == nil {
		return nil
	}

	return mask(err)
}

func mask(err error) error {
	_, file, line, _ := runtime.Caller(2)

	return stackedError{
		stackEntry: StackEntry{
			File: file,
			Line: line,
		},
		underlying: err,
	}
}
