package microerror

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/juju/errgo"
)

func TestMask_Nil(t *testing.T) {
	err := Mask(nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestMaskf_Nil(t *testing.T) {
	err := Maskf(nil, "test annotation")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestMask_Location(t *testing.T) {
	err := fmt.Errorf("test")

	err = Mask(err)
	err = Mask(err)
	err = Mask(err)

	errgoErr, ok := err.(*errgo.Err)
	if !ok {
		t.Fatalf("expected type *errgo.Err, got %T", err)
	}

	file := filepath.Base(errgoErr.Location().File)
	wfile := "error_test.go"
	if file != wfile {
		t.Fatalf("expected  %s, got %s", wfile, file)
	}
}

func TestMaskf_Location(t *testing.T) {
	err := fmt.Errorf("test")

	err = Maskf(err, "1")
	err = Maskf(err, "2")
	err = Maskf(err, "3")

	errgoErr, ok := err.(*errgo.Err)
	if !ok {
		t.Fatalf("expected type *errgo.Err, got %T", err)
	}

	file := filepath.Base(errgoErr.Location().File)
	wfile := "error_test.go"
	if file != wfile {
		t.Fatalf("expected  %s, got %s", wfile, file)
	}
}
