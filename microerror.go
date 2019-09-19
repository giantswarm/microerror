package microerror

import (
	"fmt"
	"runtime"

	"errors"
)

func Is(err, target error) bool {
	return errors.Is(err, target)
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
