package microerror

import (
	"errors"
	"go/build"
	"reflect"
	"strconv"
	"testing"
)

func Test_Stack(t *testing.T) {
	testCases := []struct {
		name          string
		inputErr      error
		expectedStack string
	}{
		{
			name:          "case 0: annotated microerror error",
			inputErr:      Maskf(&Error{Kind: "testKind"}, "annotation"),
			expectedStack: "[{" + build.Default.GOPATH + "/src/github.com/giantswarm/microerror/funcs_test.go:19: annotation} {test kind}]",
		},
		{
			name:          "case 1: non annotated microerror error",
			inputErr:      &Error{Kind: "testKind"},
			expectedStack: "test kind",
		},
		{
			name:          "case 2: external error",
			inputErr:      errors.New("external error"),
			expectedStack: "external error",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			stack := Stack(tc.inputErr)
			if !reflect.DeepEqual(stack, tc.expectedStack) {
				t.Fatalf("stack = %q, want %q", stack, tc.expectedStack)
			}
		})
	}
}
