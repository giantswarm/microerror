package microerror

import (
	"fmt"
	"strings"
	"testing"
)

func Test_Error_String_1(t *testing.T) {
	c := ErrorHandlerConfig{}
	h := NewErrorHandler(c)

	e := fmt.Errorf("test error")

	err := h.Mask(e)
	err = h.Mask(err)
	err = h.Mask(err)
	s := err.(*Error).String()

	prefix := `{"cause":{"message":"test error"},"desc":"This is the`

	if !strings.HasPrefix(s, prefix) {
		t.Fatalf("expected %s to have prefix %s", s, prefix)
	}
}

func Test_Error_String_2(t *testing.T) {
	c := ErrorHandlerConfig{}
	h := NewErrorHandler(c)

	e := &Error{
		Kind: "testError",
	}

	err := h.Mask(e)
	err = h.Mask(err)
	err = h.Mask(err)
	s := err.(*Error).String()

	prefix := `{"cause":{"message":"test error"},"desc":"This is the`

	if !strings.HasPrefix(s, prefix) {
		t.Fatalf("expected %s to have prefix %s", s, prefix)
	}
}

func Test_Error_Error_1(t *testing.T) {
	c := ErrorHandlerConfig{}
	h := NewErrorHandler(c)

	e := fmt.Errorf("test error")

	err := h.Mask(e)
	err = h.Mask(err)
	err = h.Mask(err)

	m := err.Error()
	s := "test error"
	if m != s {
		t.Fatalf("expected %s to equal %s", m, s)
	}
}

func Test_Error_Error_2(t *testing.T) {
	c := ErrorHandlerConfig{}
	h := NewErrorHandler(c)

	e := &Error{
		Kind: "testError",
	}

	err := h.Mask(e)
	err = h.Mask(err)
	err = h.Mask(err)

	m := err.Error()
	s := "test error"
	if m != s {
		t.Fatalf("expected %s to equal %s", m, s)
	}
}

func Test_Error_Error_3(t *testing.T) {
	c := ErrorHandlerConfig{}
	h := NewErrorHandler(c)

	e := &Error{
		Desc: "test desc",
		Kind: "testError",
	}

	err := h.Mask(e)
	err = h.Mask(err)
	err = h.Mask(err)

	m := err.Error()
	s := "test error: test desc"
	if m != s {
		t.Fatalf("expected %s to equal %s", m, s)
	}
}

func Test_Error_SetDescf(t *testing.T) {
	c := ErrorHandlerConfig{}
	h := NewErrorHandler(c)

	e := &Error{
		Desc: "test desc",
		Kind: "testError",
	}

	err := h.Mask(e.SetDescf("test %#q", "format"))
	err = h.Mask(err)
	err = h.Mask(err)

	m := err.Error()
	s := "test error: test `format`"
	if m != s {
		t.Fatalf("expected %s to equal %s", m, s)
	}
}

func Test_Error_SetStack_1(t *testing.T) {
	c := ErrorHandlerConfig{}
	h := NewErrorHandler(c)

	e1 := &Error{
		Desc: "test desc 1",
		Kind: "testKind1",
	}
	e2 := &Error{
		Desc: "test desc 2",
		Kind: "testKind2",
		Stack: []Stack{
			{},
			{},
			{},
		},
	}

	err := h.Mask(e1.SetStack(e2).SetDescf(e2.Error()))
	err = h.Mask(err)
	err = h.Mask(err)

	m := err.Error()
	s := "test kind1: test kind2: test desc 2"
	if m != s {
		t.Fatalf("expected %s to equal %s", m, s)
	}

	l := len(err.(*Error).Stack)
	if l != 6 {
		t.Fatalf("expected %d got %d", 6, l)
	}
}

func Test_Error_SetStack_2(t *testing.T) {
	c := ErrorHandlerConfig{}
	h := NewErrorHandler(c)

	e1 := &Error{
		Desc: "test desc 1",
		Kind: "testKind",
	}
	e2 := &Error{
		Desc: "test desc 2",
		Kind: "testKind",
		Stack: []Stack{
			{},
			{},
			{},
		},
	}

	err := h.Mask(e1.SetStack(e2).SetDescf(e2.Error()))
	err = h.Mask(err)
	err = h.Mask(err)

	m := err.Error()
	s := "test kind: test desc 2"
	if m != s {
		t.Fatalf("expected %s to equal %s", m, s)
	}

	l := len(err.(*Error).Stack)
	if l != 6 {
		t.Fatalf("expected %d got %d", 6, l)
	}
}
