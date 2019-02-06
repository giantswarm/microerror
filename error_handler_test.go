package microerror

import (
	"fmt"
	"path/filepath"
	"testing"
)

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

func Test_ErrorHandler_Maskf_Nil(t *testing.T) {
	c := ErrorHandlerConfig{}
	handler := NewErrorHandler(c)
	err := handler.Maskf(nil, "test annotation")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func Test_ErrorHandler_Cause(t *testing.T) {
	testCases := []struct {
		Name             string
		ErrorFunc        func() error
		ExpectedFiles    []string
		ExpectedLines    []int
		ExpectedMessages []string
	}{
		{
			Name: "Case 0",
			ErrorFunc: func() error {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				err := fmt.Errorf("test error")

				err = h.Mask(err)

				return err
			},
			ExpectedFiles: []string{
				"error_handler_test.go",
			},
			ExpectedLines: []int{
				50,
			},
			ExpectedMessages: []string{
				"test error",
			},
		},
		{
			Name: "Case 1",
			ErrorFunc: func() error {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				err := fmt.Errorf("test error")

				err = h.Mask(err)
				err = h.Mask(err)

				return err
			},
			ExpectedFiles: []string{
				"error_handler_test.go",
				"error_handler_test.go",
			},
			ExpectedLines: []int{
				72,
				73,
			},
			ExpectedMessages: []string{
				"test error",
				"",
			},
		},
		{
			Name: "Case 2",
			ErrorFunc: func() error {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				err := fmt.Errorf("test error")

				err = h.Mask(err)
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
				98,
				99,
				102,
			},
			ExpectedMessages: []string{
				"test error",
				"",
				"",
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
				err = h.Maskf(err, "annotation")

				return err
			},
			ExpectedFiles: []string{
				"error_handler_test.go",
				"error_handler_test.go",
				"error_handler_test.go",
			},
			ExpectedLines: []int{
				132,
				133,
				134,
			},
			ExpectedMessages: []string{
				"test error",
				"",
				"annotation",
			},
		},
		{
			Name: "Case 4",
			ErrorFunc: func() error {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				err := fmt.Errorf("test error")

				err = h.Mask(err)
				err = h.Mask(err)

				return err
			},
			ExpectedFiles: []string{
				"error_handler_test.go",
				"error_handler_test.go",
			},
			ExpectedLines: []int{
				162,
				163,
			},
			ExpectedMessages: []string{
				"test error",
				"",
			},
		},
		{
			Name: "Case 5",
			ErrorFunc: func() error {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				err := fmt.Errorf("test error")

				err = h.Mask(err)
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
				188,
				189,
				192,
			},
			ExpectedMessages: []string{
				"test error",
				"",
				"",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.ErrorFunc()

			e, ok := err.(*Error)
			if !ok {
				t.Fatalf("expected %T got %T", &Error{}, err)
			}

			if len(e.Stack) != len(tc.ExpectedLines) {
				t.Fatalf("expected %d got %d", len(tc.ExpectedLines), len(e.Stack))
			}
			if len(e.Stack) != len(tc.ExpectedFiles) {
				t.Fatalf("expected %d got %d", len(tc.ExpectedFiles), len(e.Stack))
			}
			if len(e.Stack) != len(tc.ExpectedMessages) {
				t.Fatalf("expected %d got %d", len(tc.ExpectedMessages), len(e.Stack))
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
			for i, _ := range e.Stack {
				if e.Stack[i].Message != tc.ExpectedMessages[i] {
					t.Fatalf("expected %s got %s", tc.ExpectedMessages[i], e.Stack[i].Message)
				}
			}
		})
	}
}

func Test_ErrorHandler_Error(t *testing.T) {
	testCases := []struct {
		Name            string
		ErrorFunc       func() string
		ExpectedMessage string
	}{
		{
			Name: "Case 0",
			ErrorFunc: func() string {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				e := fmt.Errorf("test error")

				err := h.Mask(e)
				err = h.Maskf(err, "annotation")

				return err.Error()
			},
			ExpectedMessage: "test error",
		},
		{
			Name: "Case 1",
			ErrorFunc: func() string {
				c := ErrorHandlerConfig{}
				h := NewErrorHandler(c)

				e := &Error{
					Kind: "testError",
				}

				err := h.Mask(e)
				err = h.Mask(e)
				err = h.Maskf(e, "annotation")

				return err.Error()
			},
			ExpectedMessage: "test error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			message := tc.ErrorFunc()

			if message != tc.ExpectedMessage {
				t.Fatalf("expected %s got %s", tc.ExpectedMessage, message)
			}
		})
	}
}
