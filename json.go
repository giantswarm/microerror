package microerror

import (
	"encoding/json"
	"errors"
)

// JSON prints the error with enriched information in JSON format. Enriched
// information includes:
//
//	- All fields from Error type.
//	- Error stack.
//
// The rendered JSON can be unmarshalled with JSONError type.
func JSON(err error) string {
	// This is ugly but this is consequence of us using Error type as
	// pointer. For being backward compatible we have to use pointer to
	// a pointer here.
	microErr := &Error{}
	if !errors.As(err, &microErr) && !errors.As(err, &stackedError{}) {
		err = annotatedError{
			annotation: err.Error(),
			underlying: &Error{
				Kind: kindUnknown,
			},
		}
	}

	bytes, err := json.Marshal(err)
	if err != nil {
		panic(err.Error())
	}

	return string(bytes)
}
