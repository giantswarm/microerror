package microerror

import "testing"

func Test_DeprecatedPrependStack_1(t *testing.T) {
	c := ErrorHandlerConfig{}
	h := NewErrorHandler(c)

	e1 := &Error{
		Desc: "test desc 1",
		Kind: "testKind1",
	}
	e2 := &maskedError{
		Cause: withAnnotation(&Error{
			Desc: "test desc 2",
			Kind: "testKind2",
		}, ""),
		Stack: []maskedErrorStack{
			{},
			{},
			{},
		},
	}

	err := h.Maskf(e1, e2.Error())
	err = DeprecatedPrependStack(err, e2)
	err = h.Mask(err)
	err = h.Mask(err)

	m := err.Error()
	s := "test kind1: test kind2"
	if m != s {
		t.Fatalf("expected %q to equal %q", m, s)
	}

	l := len(err.(*maskedError).Stack)
	if l != 6 {
		t.Fatalf("expected %d got %d", 6, l)
	}
}

func Test_DeprecatedPrependStack_2(t *testing.T) {
	c := ErrorHandlerConfig{}
	h := NewErrorHandler(c)

	e1 := &Error{
		Desc: "Test desc 1.",
		Kind: "testKind1",
	}
	e2 := &maskedError{
		Cause: withAnnotation(&Error{
			Desc: "Test desc 2.",
			Kind: "testKind2",
		}, ""),
		Stack: []maskedErrorStack{
			{},
			{},
			{},
		},
	}

	err := h.Maskf(e1, e2.Error())
	err = DeprecatedPrependStack(err, e2)
	err = h.Mask(err)
	err = h.Mask(err)

	m := err.Error()
	s := "test kind1: test kind2"
	if m != s {
		t.Fatalf("expected %q to equal %q", m, s)
	}

	l := len(err.(*maskedError).Stack)
	if l != 6 {
		t.Fatalf("expected %d got %d", 6, l)
	}
}
