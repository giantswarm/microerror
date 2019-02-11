package microerror

import "encoding/json"

type maskedError struct {
	Cause *annotatedError    `json:"cause"`
	Stack []maskedErrorStack `json:"stack,omitempty"`
}

type maskedErrorStack struct {
	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
}

func newMaskedError(err error) *maskedError {
	e, ok := err.(*maskedError)
	if ok {
		return e
	}

	e = &maskedError{
		Cause: withAnnotation(err, ""),
	}

	return e
}

func (e *maskedError) Error() string {
	return e.Cause.Error()
}

// TODO test with microerror
// TODO test with arbitrary error
func (e *maskedError) MarshalJSON() ([]byte, error) {
	microErr, ok := e.Cause.error.(*Error)
	if !ok {
		microErr = defaultError
	}

	type MaskedErrorClone maskedError

	v := struct {
		*MaskedErrorClone `json:",inline"`
		*Error            `json:",inline"`
	}{
		MaskedErrorClone: (*MaskedErrorClone)(e),
		Error:            microErr,
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return b, nil
}
