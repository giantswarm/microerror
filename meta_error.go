package microerror

import (
	"encoding/json"
)

// MetaError is a predefined error structure which purpose is to act as
// container for meta information associated to a specific error. The specific
// error type matching can be used as usual. The usual error masking and cause
// gathering can be used as usual. Using MetaError might look as follows. In the
// beginning is a usual error defined, along with its matcher. This error is the
// root cause once emitted during runtime.
//
//     var notEnoughWorkersError = microerror.MetaError{
//         Desc: "The amount of requested guest cluster workers exceeds the available number of host cluster nodes.",
//         Docs: "https://github.com/giantswarm/ops-recipes/blob/master/349-not-enough-workers.md",
//         Kind: "notEnoughWorkersError",
//     }
//
//     func IsNotEnoughWorkers(err error) bool {
//         return microerror.Cause(err) == notEnoughWorkersError
//     }
//
type MetaError struct {
	Desc string `json:"desc"`
	Docs string `json:"docs"`
	Kind string `json:"kind"`
}

func (m *MetaError) Error() string {
	return toStringCase(m.Kind)
}

func (m *MetaError) GoString() string {
	return m.String()
}

func (m *MetaError) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, Mask(err)
	}

	return b, nil
}

func (m *MetaError) String() string {
	b, err := m.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (m *MetaError) UnmarshalJSON(b []byte) error {
	var c MetaError
	err := json.Unmarshal(b, &c)
	if err != nil {
		return Mask(err)
	}

	*m = c

	return nil
}
