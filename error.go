package microerror

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}

	k := toStringCase(e.Kind)

	if e.Desc != "" {
		return fmt.Sprintf("%s: %s", k, e.Desc)
	}

	return k
}

func (e *Error) MarshalJSON() ([]byte, error) {
	type ErrorClone Error

	b, err := json.Marshal(&struct {
		*ErrorClone
	}{
		ErrorClone: (*ErrorClone)(e),
	})
	if err != nil {
		return nil, err
	}

	return b, nil
}

// SetDescf fills the error's Desc property at runtime. This is useful when
// returning the route cause of an error at runtime. We usually reuse certain
// error types for different situations. Then contextual information need to be
// transported to make the error more useful. The example below shows how this
// could look like.
//
//     return microerror.Mask(executionFailedError.SetDescf("NS record %#q for HostedZone %#q not found", name, id))
//
// Note that modifying the description of a received error is considered an
// anti-pattern. This is not easily possible since the usually used error
// interface requires type assertions for *Error in order to change its
// description again. When receiving errors one should only use Mask and forward
// the received error to the next caller. When extending the description of one
// error with the message of another and in case the error kind is the same for
// both of these errors, SetDescf removes the prefix in f equal to Kind of e. If
// we would not prune the prefix, we would have issues with repetitive error
// kinds reflected in resulting error messages. One famous example looks as
// follows.
//
//     scaling.min and scaling.max must be equal on azure: invalid request error: invalid request error: invalid request error
//
func (e *Error) SetDescf(f string, args ...interface{}) *Error {
	k := toStringCase(e.Kind)
	if strings.HasPrefix(f, k) {
		f = f[len(k)+2 : len(f)]
	}

	e.Desc = fmt.Sprintf(f, args...)

	return e
}

// SetStack overwrites the stack of the current error with the stack of the
// error given. In case the given error does not assert to *Error, the current
// stack is not overwritten. This functionality is considered to be temporary
// and should only be used where absolutely necessary. One example we care about
// right now is the redirection of service errors in endpoints. In the example
// scenario a masked error is received, and for legacy reasons the endpoint does
// not want to forward the received error type, but rather its own. The downside
// of this technique is that we would lose the stack carried with the received
// error. Using SetStack the desired error type can be used and filled with the
// stack transported by the received error. The example below also shows how to
// preserve the original error message in such legacy situations.
//
//     return microerror.Mask(notFoundError.SetStack(err).SetDescf(err.Error()))
//
func (e *Error) SetStack(err error) *Error {
	s, ok := err.(*Error)
	if !ok {
		return e
	}

	e.Stack = s.Stack

	return e
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
