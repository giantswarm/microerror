package microerror

import (
	"encoding/json"

	"github.com/giantswarm/microerror"
)

type metaError struct {
	Message  string
	Metadata map[string]string
}

func NewMetaError(message string) error {
	m := &metaError{
		Message:  message,
		Metadata: map[string]string{},
	}

	return m
}

func (m *metaError) Error() string {
	return m.Message
}

func (m *metaError) UnmarshalJSON(b []byte) error {
	var c metaError
	err := json.Unmarshal(b, &c)
	if err != nil {
		return microerror.Mask(err)
	}

	*m = c

	return nil
}

func (m *metaError) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return microerror.Mask(err)
	}

	return b, nil
}

func FromError(err error, key string) (string, error) {
	m, ok := err.(*metaError)
	if !ok {
		return "", microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &metaError{}, err)
	}

	v, ok := m.Metadata[key]
	if !ok {
		return "", microerror.Maskf(metadataNotFoundError, "no value for key '%s'", key)
	}

	return v, nil
}

var wrongTypeError = microerror.New("wrong type")

// IsWrongTypeError asserts wrongTypeError.
func IsWrongTypeError(err error) bool {
	return microerror.Cause(err) == wrongTypeError
}
