package microerror

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/juju/errgo"
)

func TestErrgoHandler_Interface(t *testing.T) {
	// This will not complie if ErrgoHandler does not fulfill Handler
	// interface.
	var _ Handler = NewErrgoHandler()
}

func TestErrgoHandler_Mask_Nil(t *testing.T) {
	handler := NewErrgoHandler()
	err := handler.Mask(nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestErrgoHandler_Maskf_Nil(t *testing.T) {
	handler := NewErrgoHandler()
	err := handler.Maskf(nil, "test annotation")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestErrgoHandler_Mask_Location(t *testing.T) {
	handler := NewErrgoHandler()

	err := fmt.Errorf("test")

	err = handler.Mask(err)
	err = handler.Mask(err)
	err = handler.Mask(err)

	errgoErr, ok := err.(*errgo.Err)
	if !ok {
		t.Fatalf("expected type *errgo.Err, got %T", err)
	}

	file := filepath.Base(errgoErr.Location().File)
	wfile := "errgo_handler_test.go"
	if file != wfile {
		t.Fatalf("expected  %s, got %s", wfile, file)
	}
}

func TestErrgoHandler_Maskf_Location(t *testing.T) {
	handler := NewErrgoHandler()

	err := fmt.Errorf("test")

	err = handler.Maskf(err, "1")
	err = handler.Maskf(err, "2")
	err = handler.Maskf(err, "3")

	errgoErr, ok := err.(*errgo.Err)
	if !ok {
		t.Fatalf("expected type *errgo.Err, got %T", err)
	}

	file := filepath.Base(errgoErr.Location().File)
	wfile := "errgo_handler_test.go"
	if file != wfile {
		t.Fatalf("expected  %s, got %s", wfile, file)
	}
}
