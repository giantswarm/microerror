package microerror

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// This test uses golden files.
//
// Run this command to update the snapshots:
// go test . -run TestPretty -update
func TestPretty(t *testing.T) {
	testCases := []struct {
		name               string
		errorFactory       func() error
		stackTrace         bool
		expectedGoldenFile string
	}{
		{
			name: "case 0: simple error",
			errorFactory: func() error {
				err := errors.New("something went wrong")

				return err
			},
			expectedGoldenFile: "pretty-simple-error.golden",
		},
		{
			name: "case 1: simple empty error",
			errorFactory: func() error {
				err := errors.New("")

				return err
			},
			expectedGoldenFile: "pretty-simple-empty-error.golden",
		},
		{
			name: "case 2: simple error with 'error:' prefix in message",
			errorFactory: func() error {
				err := errors.New("error: something went wrong")

				return err
			},
			expectedGoldenFile: "pretty-simple-error-with-prefix.golden",
		},
		{
			name: "case 3: masked simple error",
			errorFactory: func() error {
				err := errors.New("something went wrong")

				return Mask(err)
			},
			expectedGoldenFile: "pretty-masked-simple-error.golden",
		},
		{
			name: "case 4: microerror, 0 depth",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return err
			},
			expectedGoldenFile: "pretty-microerror-0-depth.golden",
		},
		{
			name: "case 5: microerror, 1 depth",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return Mask(err)
			},
			expectedGoldenFile: "pretty-microerror-1-depth.golden",
		},
		{
			name: "case 6: microerror, 1 depth, with annotation",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return Maskf(err, "something bad happened, and we had to crash")
			},
			expectedGoldenFile: "pretty-microerror-1-depth-annotation.golden",
		},
		{
			name: "case 7: microerror, 1 depth, unknown kind, with annotation",
			errorFactory: func() error {
				err := &Error{
					Kind: kindUnknown,
				}

				return Maskf(err, "something bad happened, and we had to crash")
			},
			expectedGoldenFile: "pretty-microerror-1-depth-unknown-kind-annotation.golden",
		},
		{
			name: "case 8: microerror, 1 depth, nil kind, with annotation",
			errorFactory: func() error {
				err := &Error{
					Kind: kindNil,
				}

				return Maskf(err, "something bad happened, and we had to crash")
			},
			expectedGoldenFile: "pretty-microerror-1-depth-unknown-kind-annotation.golden",
		},
		{
			name: "case 9: microerror, 1 depth, with multiline annotation",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return Maskf(err, "something bad happened, and we had to crash\nthat's the first time it happened, really")
			},
			expectedGoldenFile: "pretty-microerror-1-depth-multiline-annotation.golden",
		},
		{
			name: "case 10: microerror, 10 depth",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				// Let's build up this stack trace.
				newErr := Mask(err)
				for i := 0; i < 10; i++ {
					newErr = Mask(newErr)
				}

				return newErr
			},
			expectedGoldenFile: "pretty-microerror-10-depth.golden",
		},
		{
			name: "case 11: microerror, 1 depth, with stack trace",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return Mask(err)
			},
			stackTrace:         true,
			expectedGoldenFile: "pretty-microerror-1-depth-stack-trace.golden",
		},
		{
			name: "case 12: microerror, 1 depth, with annotation, with stack trace",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return Maskf(err, "something bad happened, and we had to crash")
			},
			stackTrace:         true,
			expectedGoldenFile: "pretty-microerror-1-depth-annotation-stack-trace.golden",
		},
		{
			name: "case 13: microerror, 1 depth, with multiline annotation, with stack trace",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return Maskf(err, "something bad happened, and we had to crash\nthat's the first time it happened, really")
			},
			stackTrace:         true,
			expectedGoldenFile: "pretty-microerror-1-depth-multiline-annotation-stack-trace.golden",
		},
		{
			name: "case 14: microerror, 10 depth, with stack trace",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				// Let's build up this stack trace.
				newErr := Mask(err)
				for i := 0; i < 10; i++ {
					newErr = Mask(newErr)
				}

				return newErr
			},
			stackTrace:         true,
			expectedGoldenFile: "pretty-microerror-10-depth-stack-trace.golden",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.errorFactory()
			message := Pretty(err, tc.stackTrace)

			// Change paths to avoid prefixes like
			// "/Users/username/go/src/" so this can test can be
			// executed on different machines.
			{
				r := regexp.MustCompile(`/.*(/.*\.go:\d+)`)
				message = r.ReplaceAllString(message, "--REPLACED--$1")
			}

			var expected string
			{
				golden := filepath.Join("testdata", tc.expectedGoldenFile)
				if *update {
					err := os.WriteFile(golden, []byte(message), 0644) //nolint:gosec
					if err != nil {
						t.Fatal(err)
					}
				}

				bytes, err := os.ReadFile(golden)
				if err != nil {
					t.Fatal(err)
				}

				expected = string(bytes)
			}

			if message != expected {
				t.Fatalf("expected %q got %q", expected, message)
			}
		})
	}
}
