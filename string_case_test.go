package microerror

import (
	"testing"
)

func Test_toStringCase(t *testing.T) {
	testCases := []struct {
		Name           string
		InputString    string
		ExpectedString string
	}{
		{
			Name:           "case 0: camel case to string case with lower start",
			InputString:    "fooBar",
			ExpectedString: "foo bar",
		},
		{
			Name:           "case 1: camel case to string case with lower start and longer input",
			InputString:    "fooBarBazupKick",
			ExpectedString: "foo bar bazup kick",
		},
		{
			Name:           "case 2: camel case to string case with upper start",
			InputString:    "FooBar",
			ExpectedString: "foo bar",
		},
		{
			Name:           "case 3: camel case to string case with upper start and longer input",
			InputString:    "FooBarBazupKick",
			ExpectedString: "foo bar bazup kick",
		},
		{
			Name:           "case 4: real private error kind",
			InputString:    "authenticationError",
			ExpectedString: "authentication error",
		},
		{
			Name:           "case 5: real public error kind",
			InputString:    "AuthenticationError",
			ExpectedString: "authentication error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			output := toStringCase(tc.InputString)
			if output != tc.ExpectedString {
				t.Fatalf("expected %#v got %#v", tc.ExpectedString, output)
			}
		})
	}
}
