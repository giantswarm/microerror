package microerror

var defaultError = &Error{
	Desc: "This is the default microerror error. It wraps an arbitrary third party error. See more information in the transported cause and stack.",
	Docs: "https://github.com/giantswarm/microerror",
	Kind: "defaultError",
}

// IsDefault asserts defaultError.
func IsDefault(err error) bool {
	return Cause(err) == defaultError
}

// Error is a predefined error structure whose purpose is to act as container
// for meta information associated to a specific error. The specific error type
// matching can be used as usual. The usual error masking and cause gathering
// should be used. Using Error might look as follows. In the beginning is
// a usual error defined, along with its matcher. This error is the root cause
// once emitted during runtime.
//
//	runtimevar notEnoughWorkersError = &microerror.Error{
//	    Docs: "https://github.com/giantswarm/ops-recipes/blob/master/349-not-enough-workers.md",
//	    Kind: "notEnoughWorkersError",
//	}
//
//	func IsNotEnoughWorkers(err error) bool {
//	    return microerror.Cause(err) == notEnoughWorkersError
//	}
//
type Error struct {
	Desc string `json:"desc,omitempty"`
	Docs string `json:"docs,omitempty"`
	Kind string `json:"kind,omitempty"`
}

func (e *Error) Error() string {
	return toStringCase(e.Kind)
}
