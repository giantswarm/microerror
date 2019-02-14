package microerror

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
)

func Test_ErrorHandler_Stack(t *testing.T) {
	testCases := []struct {
		Name          string
		ErrorFunc     func() error
		ExpectedFiles []string
		ExpectedLines []int
	}{
		{
			Name: "Case 0",
			ErrorFunc: func() error {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				e := fmt.Errorf("test error")

				err := h.Mask(e)

				return err
			},
			ExpectedFiles: []string{
				"error_handler_test.go",
			},
			ExpectedLines: []int{
				25,
			},
		},
		{
			Name: "Case 1",
			ErrorFunc: func() error {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				e := fmt.Errorf("test error")

				err := h.Mask(e)
				err = h.Mask(err)

				return err
			},
			ExpectedFiles: []string{
				"error_handler_test.go",
				"error_handler_test.go",
			},
			ExpectedLines: []int{
				44,
				45,
			},
		},
		{
			Name: "Case 2",
			ErrorFunc: func() error {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				e := fmt.Errorf("test error")

				err := h.Mask(e)
				err = h.Mask(err)
				// Note this comment is kept for the test to ensure a bump in the line
				// number recording.
				err = h.Mask(err)

				return err
			},
			ExpectedFiles: []string{
				"error_handler_test.go",
				"error_handler_test.go",
				"error_handler_test.go",
			},
			ExpectedLines: []int{
				66,
				67,
				70,
			},
		},
		{
			Name: "Case 3",
			ErrorFunc: func() error {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				e := &Error{
					Kind: "testError",
				}

				err := h.Mask(e)
				err = h.Mask(err)
				err = h.Mask(err)

				return err
			},
			ExpectedFiles: []string{
				"error_handler_test.go",
				"error_handler_test.go",
				"error_handler_test.go",
			},
			ExpectedLines: []int{
				95,
				96,
				97,
			},
		},
		{
			Name: "Case 4",
			ErrorFunc: func() error {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				e := fmt.Errorf("test error")

				err := h.Mask(e)
				err = h.Mask(err)

				return err
			},
			ExpectedFiles: []string{
				"error_handler_test.go",
				"error_handler_test.go",
			},
			ExpectedLines: []int{
				120,
				121,
			},
		},
		{
			Name: "Case 5",
			ErrorFunc: func() error {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				e := fmt.Errorf("test error")

				err := h.Mask(e)
				err = h.Mask(err)
				// Note this comment is kept for the test to ensure a bump in the line
				// number recording.
				err = h.Mask(err)

				return err
			},
			ExpectedFiles: []string{
				"error_handler_test.go",
				"error_handler_test.go",
				"error_handler_test.go",
			},
			ExpectedLines: []int{
				142,
				143,
				146,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.ErrorFunc()

			e, ok := err.(*maskedError)
			if !ok {
				t.Fatalf("expected %T got %T", &Error{}, err)
			}

			if len(e.Stack) != len(tc.ExpectedLines) {
				t.Fatalf("expected %d got %d", len(tc.ExpectedLines), len(e.Stack))
			}
			if len(e.Stack) != len(tc.ExpectedFiles) {
				t.Fatalf("expected %d got %d", len(tc.ExpectedFiles), len(e.Stack))
			}

			for i, _ := range e.Stack {
				if filepath.Base(e.Stack[i].File) != tc.ExpectedFiles[i] {
					t.Fatalf("expected %s got %s", tc.ExpectedFiles[i], filepath.Base(e.Stack[i].File))
				}
			}
			for i, _ := range e.Stack {
				if e.Stack[i].Line != tc.ExpectedLines[i] {
					t.Fatalf("expected %d got %d", tc.ExpectedLines[i], e.Stack[i].Line)
				}
			}
		})
	}
}
func Test_ErrorHandler_Interface(t *testing.T) {
	// This will not complie if ErrorHandler does not fulfill Handler
	// interface.
	c := ErrorHandlerConfig{}
	var _ Handler = NewErrorHandler(c)
}

func Test_ErrorHandler_Mask_Nil(t *testing.T) {
	c := ErrorHandlerConfig{}
	handler := NewErrorHandler(c)
	err := handler.Mask(nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func Test_ErrorHandler_Maskf(t *testing.T) {
	c := ErrorHandlerConfig{}
	h := NewErrorHandler(c)

	e := &Error{
		Desc: "test desc",
		Kind: "testError",
	}

	err := h.Maskf(e, "test %#q", "format")
	err = h.Mask(err)
	err = h.Mask(err)

	m := err.Error()
	s := "test error: test `format`"
	if m != s {
		t.Fatalf("expected %s to equal %s", m, s)
	}
}

func Test_ErrorHandler_Error(t *testing.T) {
	testCases := []struct {
		name            string
		inputError      func(h *ErrorHandler) error
		expectedMessage string
	}{
		{
			name: "case 0",
			inputError: func(h *ErrorHandler) error {
				e := fmt.Errorf("test error")

				err := h.Mask(e)
				err = h.Mask(err)
				err = h.Mask(err)

				return err
			},
			expectedMessage: "test error",
		},
		{
			name: "case 1",
			inputError: func(h *ErrorHandler) error {
				e := &Error{
					Kind: "testError",
				}

				err := h.Mask(e)
				err = h.Mask(err)
				err = h.Mask(err)

				return err
			},
			expectedMessage: "test error",
		},
		{
			name: "case 2: Maskf and Mask",
			inputError: func(h *ErrorHandler) error {
				e := &Error{
					Kind: "testError",
				}

				err := h.Maskf(e, "additional runtime info")
				err = h.Mask(err)
				err = h.Mask(err)

				return err
			},
			expectedMessage: "test error: additional runtime info",
		},
		{
			name: "case 3: only Maskf",
			inputError: func(h *ErrorHandler) error {
				e := &Error{
					Kind: "testError",
				}

				err := h.Maskf(e, "additional runtime info")

				return err
			},
			expectedMessage: "test error: additional runtime info",
		},
		{
			name: "case 3: no masking on *Error",
			inputError: func(h *ErrorHandler) error {
				e := &Error{
					Kind: "testError",
				}

				return e
			},
			expectedMessage: "test error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewErrorHandler(ErrorHandlerConfig{})

			err := tc.inputError(h)

			message := err.Error()
			if message != tc.expectedMessage {
				t.Fatalf("expected %q to equal %q", message, tc.expectedMessage)
			}
		})
	}
}

func Test_ErrorHandler_JSON(t *testing.T) {
	testCases := []struct {
		name         string
		inputError   error
		expectedKind string
		expectedDesc string
		expectedDocs string
	}{
		{
			name:         "case 0: non *Error",
			inputError:   fmt.Errorf("test error"),
			expectedDesc: defaultError.Desc,
			expectedDocs: defaultError.Docs,
			expectedKind: defaultError.Kind,
		},
		{
			name: "case 1: *Error",
			inputError: &Error{
				Desc: "Test description.",
				Docs: "https://giantswarm.io",
				Kind: "testError",
			},
			expectedDesc: "Test description.",
			expectedDocs: "https://giantswarm.io",
			expectedKind: "testError",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewErrorHandler(ErrorHandlerConfig{})

			err := h.Mask(tc.inputError)
			err = h.Mask(err)

			bytes, err := json.Marshal(err)
			if err != nil {
				t.Fatalf("expected error %v to be non nil", err)
			}

			// TODO Create a separate type for that.
			v := struct {
				Desc string `json:"desc"`
				Docs string `json:"docs"`
				Kind string `json:"kind"`
			}{}

			err = json.Unmarshal(bytes, &v)
			if err != nil {
				t.Fatalf("expected error %v to be non nil", err)
			}

			if v.Desc != tc.expectedDesc {
				t.Fatalf("expected %q to equal %q", v.Desc, tc.expectedDesc)
			}
			if v.Docs != tc.expectedDocs {
				t.Fatalf("expected %q to equal %q", v.Docs, tc.expectedDocs)
			}
			if v.Kind != tc.expectedKind {
				t.Fatalf("expected %q to equal %q", v.Kind, tc.expectedKind)
			}
		})
	}
}
