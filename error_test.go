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
	err = h.Maskf(err, "annotation")
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
	err = h.Maskf(err, "annotation")
	err = h.Mask(err)
	s := err.(*Error).String()

	prefix := `{"cause":{"message":"test error"},"desc":"This is the`

	if !strings.HasPrefix(s, prefix) {
		t.Fatalf("expected %s to have prefix %s", s, prefix)
	}
}
