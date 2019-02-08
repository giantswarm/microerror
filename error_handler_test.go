package microerror

import (
	"fmt"
	"path/filepath"
	"testing"
)

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
				err = h.Mask(err)

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
				err = h.Mask(err)

				return err.Error()
			},
			ExpectedMessage: "test error",
		},
		{
			Name: "Case 2",
			ErrorFunc: func() string {
				err := &Error{
					Kind: "testError",
				}

				return err.Error()
			},
			ExpectedMessage: "test error",
		},
		{
			Name: "Case 3",
			ErrorFunc: func() string {
				err := &Error{}

				return err.Error()
			},
			ExpectedMessage: "",
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
				96,
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
				115,
				116,
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
				137,
				138,
				141,
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
				166,
				167,
				168,
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
				191,
				192,
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
				213,
				214,
				217,
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

func Test_ErrorHandler_Cause_1(t *testing.T) {
	h := NewErrorHandler(ErrorHandlerConfig{})

	e := fmt.Errorf("test error")

	err := h.Mask(e)
	err = h.Mask(err)

	c := h.Cause(err)

	if c != e {
		t.Fatalf("expected %T to equal %T", c, e)
	}
}

func Test_ErrorHandler_Cause_2(t *testing.T) {
	h := NewErrorHandler(ErrorHandlerConfig{})

	e := &Error{
		Kind: "testError",
	}

	err := h.Mask(e)
	err = h.Mask(err)

	c := h.Cause(err)

	if c != e {
		t.Fatalf("expected %T to equal %T", c, e)
	}
}

func Test_ErrorHandler_Desc_1(t *testing.T) {
	h := NewErrorHandler(ErrorHandlerConfig{})

	e := fmt.Errorf("test error")

	err := h.Mask(e)
	err = h.Mask(err)

	d := err.(*Error).Desc
	s := "This is the default microerror error. It wraps an arbitrary third party error. See more information in the transported cause and stack."
	if d != s {
		t.Fatalf("expected %s to equal %s", d, s)
	}
}

func Test_ErrorHandler_Desc_2(t *testing.T) {
	h := NewErrorHandler(ErrorHandlerConfig{})

	e := &Error{
		Desc: "test desc",
		Kind: "testError",
	}

	err := h.Mask(e)
	err = h.Mask(err)

	d := err.(*Error).Desc
	s := "test desc"
	if d != s {
		t.Fatalf("expected %s to equal %s", d, s)
	}
}

func Test_ErrorHandler_Docs_1(t *testing.T) {
	h := NewErrorHandler(ErrorHandlerConfig{})

	e := fmt.Errorf("test error")

	err := h.Mask(e)
	err = h.Mask(err)

	d := err.(*Error).Docs
	s := "https://github.com/giantswarm/microerror"
	if d != s {
		t.Fatalf("expected %s to equal %s", d, s)
	}
}

func Test_ErrorHandler_Docs_2(t *testing.T) {
	h := NewErrorHandler(ErrorHandlerConfig{})

	e := &Error{
		Docs: "test desc",
		Kind: "testError",
	}

	err := h.Mask(e)
	err = h.Mask(err)

	d := err.(*Error).Docs
	s := "test desc"
	if d != s {
		t.Fatalf("expected %s to equal %s", d, s)
	}
}

func Test_ErrorHandler_Kind_1(t *testing.T) {
	h := NewErrorHandler(ErrorHandlerConfig{})

	e := fmt.Errorf("test error")

	err := h.Mask(e)
	err = h.Mask(err)

	d := err.(*Error).Kind
	s := "defaultMicroError"
	if d != s {
		t.Fatalf("expected %s to equal %s", d, s)
	}
}

func Test_ErrorHandler_Kind_2(t *testing.T) {
	h := NewErrorHandler(ErrorHandlerConfig{})

	e := &Error{
		Kind: "testKind",
	}

	err := h.Mask(e)
	err = h.Mask(err)

	d := err.(*Error).Kind
	s := "testKind"
	if d != s {
		t.Fatalf("expected %s to equal %s", d, s)
	}
}
