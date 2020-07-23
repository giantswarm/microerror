package microerror

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"testing"
)

// This test uses golden files.
//
// Run this command to update the snapshots:
// go test . -run Test_formatStackTrace -update
//
func Test_formatStackTrace(t *testing.T) {
	testCases := []struct {
		name               string
		errorFactory       func() error
		expectedGoldenFile string
	}{
		{
			name: "case 0: format stack trace, 1 depth",
			errorFactory: func() error {
				err := &Error{
					Kind: "somethingWentWrongError",
				}

				return Mask(err)
			},
			expectedGoldenFile: "stack-trace-format-1-depth.golden",
		},
		{
			name: "case 1: format stack trace, 10 depth",
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
			expectedGoldenFile: "stack-trace-format-10-depth.golden",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.errorFactory()
			sErr, ok := err.(*stackedError)
			if !ok {
				t.Fatalf("expected error to be of type '%T', but instead got '%T'", &stackedError{}, err)
			}
			trace := createStackTrace(sErr)
			output := formatStackTrace(trace)

			// Change paths to avoid prefixes like
			// "/Users/username/go/src/" so this can test can be
			// executed on different machines.
			{
				r := regexp.MustCompile(`/.*(/.*\.go:\d+)`)
				output = r.ReplaceAllString(output, "--REPLACED--$1")
			}

			var expected string
			{
				golden := filepath.Join("testdata", tc.expectedGoldenFile)
				if *update {
					err := ioutil.WriteFile(golden, []byte(output), 0644)
					if err != nil {
						t.Fatal(err)
					}
				}

				bytes, err := ioutil.ReadFile(golden)
				if err != nil {
					t.Fatal(err)
				}

				expected = string(bytes)
			}

			if output != expected {
				t.Fatalf("stack trace not expected, got:\n%s", expected)
			}
		})
	}
}
