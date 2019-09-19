package microerror

import (
	"encoding/json"
	"fmt"
	"runtime"

	"errors"
)

const kindUnknown = "unknown"

type Error struct {
	Desc string `json:"desc,omitempty"`
	Docs string `json:"docs,omitempty"`
	Kind string `json:"kind"`
}

func (e Error) Error() string {
	return toStringCase(e.Kind)
}

type annotatedError struct {
	annotation string
	underlying Error
}

func (e annotatedError) Error() string {
	if e.annotation == "" {
		return e.underlying.Error()
	}
	return e.underlying.Error() + ": " + e.annotation
}

func (e annotatedError) MarshalJSON() ([]byte, error) {
	o := JSONError{
		Error: e.underlying,

		Annotation: e.annotation,
	}

	bytes, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("microerror.annotatedError.MarshalJSON: %w object=%#v", err, o)
	}

	return bytes, nil
}

func (e annotatedError) Unwrap() error {
	return e.underlying
}

type stackedError struct {
	stackEntry StackEntry
	underlying error
}

func (e stackedError) Error() string {
	return e.underlying.Error()
}

func (e stackedError) MarshalJSON() ([]byte, error) {
	var stack = []StackEntry{
		e.stackEntry,
	}
	{
		underlying := e.underlying
		var stacked stackedError
		for errors.As(underlying, &stacked) {
			stack = append([]StackEntry{stacked.stackEntry}, stack...)
			underlying = stacked.underlying
		}
	}

	var microErr Error
	var annotation string
	{
		if errors.As(e, &microErr) {
			var annotatedErr annotatedError
			if errors.As(e, &annotatedErr) {
				annotation = annotatedErr.annotation
			}
		} else {
			microErr = Error{
				Kind: kindUnknown,
			}

			annotation = e.Error()
		}
	}

	o := JSONError{
		Error: microErr,

		Annotation: annotation,
		Stack:      stack,
	}

	bytes, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("microerror.stackedError.MarshalJSON: %w object=%#v", err, o)
	}

	return bytes, nil
}

func (e stackedError) Unwrap() error {
	return e.underlying
}

type JSONError struct {
	Error `json:",inline"`

	Annotation string       `json:"annotation,omitempty"`
	Stack      []StackEntry `json:"stack,omitempty"`
}

type StackEntry struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

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
