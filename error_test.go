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

func TestStack(t *testing.T) {
	tests := []struct {
		desc     string
		depth    int
		maskFunc func(error) error
	}{
		{
			desc:  "Mask (1)",
			depth: 1,
			maskFunc: func(err error) error {
				return Mask(err)
			},
		},
		{
			desc:  "Mask (3)",
			depth: 3,
			maskFunc: func(err error) error {
				err = Mask(err)
				err = Mask(err)
				err = Mask(err)
				return err
			},
		},
		{
			desc:  "Maskf (3)",
			depth: 3,
			maskFunc: func(err error) error {
				err = Maskf(err, "1")
				err = Maskf(err, "2")
				err = Maskf(err, "3")
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
				wfile := "error_test.go"
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
