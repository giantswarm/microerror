package microerror

import (
	"testing"
)

func Test_MetaError(t *testing.T) {
	var err error

	errOne := invalidConfigError
	errTwo := invalidConfigError

	{
		_, err = FromMetaError(errOne, "foo")
		if !IsWrongTypeError(err) {
			t.Fatalf("expected %#v got %#v", true, false)
		}

		_, err = FromMetaError(errTwo, "foo")
		if !IsWrongTypeError(err) {
			t.Fatalf("expected %#v got %#v", true, false)
		}

		if !IsInvalidConfig(errOne) {
			t.Fatalf("expected %#v got %#v", true, false)
		}

		if !IsInvalidConfig(errTwo) {
			t.Fatalf("expected %#v got %#v", true, false)
		}
	}

	{
		errOne = NewMetaError(errOne, map[string]string{"foo": "barOne"})
		errTwo = NewMetaError(errTwo, map[string]string{"foo": "barTwo"})
	}

	{
		valOne, err := FromMetaError(errOne, "foo")
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
		if valOne != "barOne" {
			t.Fatalf("expected %#v got %#v", "barOne", valOne)
		}

		valTwo, err := FromMetaError(errTwo, "foo")
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
		if valTwo != "barTwo" {
			t.Fatalf("expected %#v got %#v", "barTwo", valTwo)
		}

		if !IsInvalidConfig(errOne) {
			t.Fatalf("expected %#v got %#v", true, false)
		}

		if !IsInvalidConfig(errTwo) {
			t.Fatalf("expected %#v got %#v", true, false)
		}
	}

	{
		errOne = Mask(errOne)
		errTwo = Mask(errTwo)
	}

	{
		valOne, err := FromMetaError(errOne, "foo")
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
		if valOne != "barOne" {
			t.Fatalf("expected %#v got %#v", "barOne", valOne)
		}

		valTwo, err := FromMetaError(errTwo, "foo")
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
		if valTwo != "barTwo" {
			t.Fatalf("expected %#v got %#v", "barTwo", valTwo)
		}

		if !IsInvalidConfig(errOne) {
			t.Fatalf("expected %#v got %#v", true, false)
		}

		if !IsInvalidConfig(errTwo) {
			t.Fatalf("expected %#v got %#v", true, false)
		}
	}
}
