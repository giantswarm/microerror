package microerror

import "encoding/json"

// Error is a predefined error structure whose purpose is to act as container
// for meta information associated to a specific error. The specific error type
// matching can be used as usual. The usual error masking and cause gathering
// can be used as usual. Using Error might look as follows. In the beginning is
// a usual error defined, along with its matcher. This error is the root cause
// once emitted during runtime.
//
//     var notEnoughWorkersError = &microerror.Error{
//         Desc: "The amount of requested tenant cluster workers exceeds the available number of control plane nodes.",
//         Docs: "https://github.com/giantswarm/ops-recipes/blob/master/349-not-enough-workers.md",
//         Kind: "notEnoughWorkersError",
//     }
//
//     func IsNotEnoughWorkers(err error) bool {
//         return microerror.Cause(err) == notEnoughWorkersError
//     }
//
type Error struct {
	Cause error   `json:"cause,omitempty"`
	Desc  string  `json:"desc,omitempty"`
	Docs  string  `json:"docs,omitempty"`
	Kind  string  `json:"kind,omitempty"`
	Stack []Stack `json:"stack,omitempty"`
}

type Stack struct {
	File    string `json:"file,omitempty"`
	Line    int    `json:"line,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}

	return toStringCase(e.Kind)
}

func (e *Error) GoString() string {
	return e.String()
}

func (e *Error) String() string {
	// When the current cause is not of type *Error, we want to ensure the
	// transported cause is marshaled properly. Arbitrary error types are stupid
	// interfaces and marshal to nil otherwise. Here we lose all kinds of custom
	// information of all kinds of error implementations out there. We try to
	// preserve the error message though.
	if e.Cause != nil {
		e.Cause = &MarshalableError{
			Message: e.Cause.Error(),
		}
	}

	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func newDefaultError() *Error {
	return &Error{
		Cause: nil,
		Desc:  "This is the default microerror error. It wraps an arbitrary third party error. See more information in the transported cause and stack.",
		Docs:  "https://github.com/giantswarm/microerror",
		Kind:  "defaultMicroError",
		Stack: nil,
	}
}
