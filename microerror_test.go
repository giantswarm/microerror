package microerror

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func Test_errors_Is(t *testing.T) {
	var testOneError = &Error{
		Kind: "testOneError",
	}

	var testTwoError = &Error{
		Kind: "testTwoError",
	}

	if !reflect.DeepEqual(errors.Is(testOneError, testOneError), true) {
		t.Fatalf("errors.Is(testOneError, testOneError) = %v, want %v", errors.Is(testOneError, testOneError), true)
	}

	if !reflect.DeepEqual(errors.Is(testOneError, testTwoError), false) {
		t.Fatalf("errors.Is(testOneError, testTwoError) = %v, want %v", errors.Is(testOneError, testTwoError), false)
	}

	var childOneError = FromParent(testOneError)

	if !reflect.DeepEqual(errors.Is(childOneError, childOneError), true) {
		t.Fatalf("errors.Is(childOneError, childOneError) = %v, want %v", errors.Is(childOneError, childOneError), true)
	}

	if !reflect.DeepEqual(errors.Is(childOneError, testOneError), true) {
		t.Fatalf("errors.Is(childOneError, testOneError) = %v, want %v", errors.Is(childOneError, testOneError), true)
	}

	if !reflect.DeepEqual(errors.Is(testOneError, childOneError), false) {
		t.Fatalf("errors.Is(testOneError, childOneError) = %v, want %v", errors.Is(testOneError, childOneError), false)
	}

	maskedOneErr := Maskf(testOneError, "maskf one")

	if !reflect.DeepEqual(errors.Is(maskedOneErr, testOneError), true) {
		t.Fatalf("errors.Is(maskedOneErr, testOneError) = %v, want %v", errors.Is(maskedOneErr, testOneError), true)
	}

	if !reflect.DeepEqual(errors.Is(maskedOneErr, testTwoError), false) {
		t.Fatalf("errors.Is(maskedOneErr, testTwoError) = %v, want %v", errors.Is(maskedOneErr, testTwoError), false)
	}

	if !reflect.DeepEqual(errors.Is(maskedOneErr, childOneError), false) {
		t.Fatalf("errors.Is(maskedOneErr, childOneError) = %v, want %v", errors.Is(maskedOneErr, childOneError), false)
	}

	wrappedOneErr := Mask(maskedOneErr)

	if !reflect.DeepEqual(errors.Is(wrappedOneErr, testOneError), true) {
		t.Fatalf("errors.Is(wrappedOneErr, testOneError) = %v, want %v", errors.Is(wrappedOneErr, testOneError), true)
	}

	if !reflect.DeepEqual(errors.Is(wrappedOneErr, testTwoError), false) {
		t.Fatalf("errors.Is(wrappedOneErr, testTwoError) = %v, want %v", errors.Is(wrappedOneErr, testTwoError), false)
	}

	if !reflect.DeepEqual(errors.Is(wrappedOneErr, childOneError), false) {
		t.Fatalf("errors.Is(wrappedOneErr, childOneError) = %v, want %v", errors.Is(wrappedOneErr, childOneError), false)
	}

	maskedChildOneErr := Maskf(childOneError, "maskf child one")

	if !reflect.DeepEqual(errors.Is(maskedChildOneErr, testOneError), true) {
		t.Fatalf("errors.Is(maskedChildOneErr, testOneError) = %v, want %v", errors.Is(maskedChildOneErr, testOneError), true)
	}

	if !reflect.DeepEqual(errors.Is(maskedChildOneErr, testTwoError), false) {
		t.Fatalf("errors.Is(maskedChildOneErr, testTwoError) = %v, want %v", errors.Is(maskedChildOneErr, testTwoError), false)
	}

	if !reflect.DeepEqual(errors.Is(maskedChildOneErr, childOneError), true) {
		t.Fatalf("errors.Is(maskedChildOneErr, childOneError) = %v, want %v", errors.Is(maskedChildOneErr, childOneError), false)
	}

	wrappedChildOneErr := Mask(maskedChildOneErr)

	if !reflect.DeepEqual(errors.Is(wrappedChildOneErr, testOneError), true) {
		t.Fatalf("errors.Is(wrappedChildOneErr, testOneError) = %v, want %v", errors.Is(wrappedChildOneErr, testOneError), true)
	}

	if !reflect.DeepEqual(errors.Is(wrappedChildOneErr, testTwoError), false) {
		t.Fatalf("errors.Is(wrappedChildOneErr, testTwoError) = %v, want %v", errors.Is(wrappedChildOneErr, testTwoError), false)
	}

	if !reflect.DeepEqual(errors.Is(wrappedChildOneErr, childOneError), true) {
		t.Fatalf("errors.Is(wrappedChildOneErr, childOneError) = %v, want %v", errors.Is(wrappedChildOneErr, childOneError), false)
	}
}

func Test_Cause(t *testing.T) {
	var testCauseMicroError = &Error{
		Kind: "testCauseErrorB",
	}
	var testCauseErrorsNewError = errors.New("test cause error A")
	var testCauseErrorsNewWrappedError = fmt.Errorf("test cause error B: %w", errors.New("test cause error A"))

	testCases := []struct {
		name               string
		inputErrorFunc     func() error
		expectedCauseError error
	}{
		{
			name: "case 0: nil",
			inputErrorFunc: func() error {
				return nil
			},
			expectedCauseError: nil,
		},
		{
			name: "case 1: no masking error=microerror.Error",
			inputErrorFunc: func() error {
				err := testCauseMicroError
				return err
			},
			expectedCauseError: testCauseMicroError,
		},
		{
			name: "case 2: no masking error=errors.New",
			inputErrorFunc: func() error {
				err := testCauseErrorsNewError
				return err
			},
			expectedCauseError: testCauseErrorsNewError,
		},
		{
			name: "case 3: Maskf depth=1 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Maskf(testCauseMicroError, "test annotation")
				return err
			},
			expectedCauseError: testCauseMicroError,
		},
		{
			name: "case 4: Maskf depth=3 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Maskf(testCauseMicroError, "test annotation")
				err = Mask(err)
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedCauseError: testCauseMicroError,
		},
		{
			name: "case 5: Mask depth=1 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Mask(testCauseErrorsNewError)
				return err
			},
			expectedCauseError: testCauseErrorsNewError,
		},
		{
			name: "case 6: Mask depth=3 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Mask(testCauseErrorsNewError)
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedCauseError: testCauseErrorsNewError,
		},
		{
			name: "case 7: Mask depth=1 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Mask(testCauseErrorsNewError)
				return err
			},
			expectedCauseError: testCauseErrorsNewError,
		},
		{
			name: "case 8: Mask depth=3 error=errors.New",
			inputErrorFunc: func() error {
				err := Mask(testCauseErrorsNewError)
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedCauseError: testCauseErrorsNewError,
		},
		{
			name: "case 9: Mask depth=3 error=fmt.Printf",
			inputErrorFunc: func() error {
				err := Mask(testCauseErrorsNewWrappedError)
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedCauseError: testCauseErrorsNewWrappedError,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)
			err := tc.inputErrorFunc()
			cause := Cause(err)
			if cause != tc.expectedCauseError {
				t.Errorf("err = %#v, want %#v", cause, tc.expectedCauseError)
			}
		})
	}
}

func Test_Mask_Nil(t *testing.T) {
	err := Mask(nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

var testMicroErr = &Error{
	Desc: "test-desc",
	Docs: "test-docs",
	Kind: "testKind",
}

func Test_Mask_Error(t *testing.T) {
	testCases := []struct {
		name           string
		inputErrorFunc func() error
		expectedError  string
	}{
		{
			name: "case 0: Maskf depth=1 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Maskf(testMicroErr, "test annotation")
				return err
			},
			expectedError: "test kind: test annotation",
		},
		{
			name: "case 1: Maskf depth=3 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Maskf(testMicroErr, "test annotation")
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedError: "test kind: test annotation",
		},
		{
			name: "case 2: Maskf depth=3 error=microerror.Error & empty annotation",
			inputErrorFunc: func() error {
				err := Maskf(testMicroErr, "")
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedError: "test kind",
		},
		{
			name: "case 3: Mask depth=1 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Mask(testMicroErr)
				return err
			},
			expectedError: "test kind",
		},
		{
			name: "case 4: Mask depth=3 error=microerror.Error",
			inputErrorFunc: func() error {
				err := Mask(testMicroErr)
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedError: "test kind",
		},
		{
			name: "case 5: Mask depth=1 error=errors.New",
			inputErrorFunc: func() error {
				err := Mask(errors.New("test error"))
				return err
			},
			expectedError: "test error",
		},
		{
			name: "case 6: Mask depth=3 error=errors.New",
			inputErrorFunc: func() error {
				err := Mask(errors.New("test error"))
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedError: "test error",
		},
		{
			name: "case 9: Mask depth=3 error=fmt.Printf",
			inputErrorFunc: func() error {
				err := Mask(fmt.Errorf("test cause error B: %w", errors.New("test cause error A")))
				err = Mask(err)
				err = Mask(err)
				return err
			},
			expectedError: "test cause error B: test cause error A",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := tc.inputErrorFunc()
			if err.Error() != tc.expectedError {
				t.Errorf("err.Error() = %#q, want %#q", err.Error(), tc.expectedError)
			}
		})
	}
}
