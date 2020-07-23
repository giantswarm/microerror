package microerror

import (
	"errors"
	"testing"
)

func TestPretty(t *testing.T) {
	testCases := []struct {
		name            string
		errorFactory    func() error
		expectedMessage string
	}{
		{
			name: "case 0: simple error",
			errorFactory: func() error {
				err := errors.New("something went wrong")

				return err
			},
			expectedMessage: "Error: Something went wrong",
		},
		{
			name: "case 1: simple empty error",
			errorFactory: func() error {
				err := errors.New("")

				return err
			},
			expectedMessage: "",
		},
		{
			name: "case 2: simple error with 'error:' prefix in message",
			errorFactory: func() error {
				err := errors.New("error: something went wrong")

				return err
			},
			expectedMessage: "Error: Something went wrong",
		},
		{
			name: "case 3: masked simple error",
			errorFactory: func() error {
				err := errors.New("something went wrong")

				return Mask(err)
			},
			expectedMessage: "Error: Something went wrong",
		},
		{
			name: "case 4: microerror, 0 depth",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return err
			},
			expectedMessage: "Error: Something went wrong",
		},
		{
			name: "case 5: microerror, 1 depth",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return Mask(err)
			},
			expectedMessage: "Error: Something went wrong",
		},
		{
			name: "case 6: microerror, 1 depth, with annotation",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return Maskf(err, "something bad happened, and we had to crash")
			},
			expectedMessage: "Error: Something went wrong: Something bad happened, and we had to crash",
		},
		{
			name: "case 7: microerror, 1 depth, unknown kind, with annotation",
			errorFactory: func() error {
				err := &Error{
					Kind: kindUnknown,
				}

				return Maskf(err, "something bad happened, and we had to crash")
			},
			expectedMessage: "Error: Something bad happened, and we had to crash",
		},
		{
			name: "case 8: microerror, 1 depth, nil kind, with annotation",
			errorFactory: func() error {
				err := &Error{
					Kind: kindNil,
				}

				return Maskf(err, "something bad happened, and we had to crash")
			},
			expectedMessage: "Error: Something bad happened, and we had to crash",
		},
		{
			name: "case 9: microerror, 1 depth, with multiline annotation",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return Maskf(err, "something bad happened, and we had to crash\nthat's the first time it happened, really")
			},
			expectedMessage: "Error: Something went wrong: Something bad happened, and we had to crash\nthat's the first time it happened, really",
		},
		{
			name: "case 10: microerror, 10 depth",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				// Let's build up this stack trace.
				newErr := Mask(err)
				for i := 0; i < 11; i++ {
					newErr = Mask(newErr)
				}

				return newErr
			},
			expectedMessage: "Error: Something went wrong",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.errorFactory()
			message := Pretty(err, false)

			if message != tc.expectedMessage {
				t.Fatalf("expected %q got %q", tc.expectedMessage, message)
			}
		})
	}
}
