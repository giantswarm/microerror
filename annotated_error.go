package microerror

import "encoding/json"

type annotatedError struct {
	error

	Annotation string
}

func withAnnotation(err error, annotation string) *annotatedError {
	e, ok := err.(*annotatedError)
	if ok {
		e.Annotation = annotation
		return e
	}

	return &annotatedError{
		error: err,

		Annotation: annotation,
	}
}

func (e *annotatedError) Error() string {
	s := e.error.Error()
	if e.Annotation != "" {
		s += ": "
		s += e.Annotation
	}

	return s
}

func (e *annotatedError) MarshalJSON() ([]byte, error) {
	m := e.Error()
	if e.Annotation != "" {
		m += ": "
		m += e.Annotation
	}

	v := struct {
		Message string `json:"message"`
	}{
		Message: m,
	}

	return json.Marshal(v)
}
