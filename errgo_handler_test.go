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
	var _ Handler = NewErrgoHandler(DefaultErrgoHandlerConfig())
}

func TestErrgoHandler_Mask_Nil(t *testing.T) {
	handler := NewErrgoHandler(DefaultErrgoHandlerConfig())
	err := handler.Mask(nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestErrgoHandler_Maskf_Nil(t *testing.T) {
	handler := NewErrgoHandler(DefaultErrgoHandlerConfig())
	err := handler.Maskf(nil, "test annotation")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestErrgoHandler_Stack(t *testing.T) {
	tests := []struct {
		desc     string
		depth    int
		maskFunc func(error) error
	}{
		{
			desc:  "Mask (1)",
			depth: 1,
			maskFunc: func(err error) error {
				h := NewErrgoHandler(DefaultErrgoHandlerConfig())
				return h.Mask(err)
			},
		},
		{
			desc:  "Mask (3)",
			depth: 3,
			maskFunc: func(err error) error {
				h := NewErrgoHandler(DefaultErrgoHandlerConfig())
				err = h.Mask(err)
				err = h.Mask(err)
				err = h.Mask(err)
				return err
			},
		},
		{
			desc:  "Maskf (3)",
			depth: 3,
			maskFunc: func(err error) error {
				h := NewErrgoHandler(DefaultErrgoHandlerConfig())
				err = h.Maskf(err, "1")
				err = h.Maskf(err, "2")
				err = h.Maskf(err, "3")
				return err
			},
		},
	}

	for i, tc := range tests {
		err := tc.maskFunc(fmt.Errorf("test"))

		var depth int
		for {
			// Check err location.
			if err, ok := err.(errgo.Locationer); ok {
				file := filepath.Base(err.Location().File)
				wfile := "errgo_handler_test.go"
				if file != wfile {
					t.Errorf("#%d %s: expected  %s, got %s", i, tc.desc, wfile, file)
				}
			}

			if cerr, ok := err.(errgo.Wrapper); ok {
				depth++
				err = cerr.Underlying()
			} else {
				break
			}
		}

		if tc.depth != depth {
			t.Fatalf("#%d %s: expected depth = %d, got %d", i, tc.desc, tc.depth, depth)
		}

	}
}
