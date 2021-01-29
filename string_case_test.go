package microerror

import (
	"reflect"
	"strconv"
	"testing"
)

func Test_toCamelCase(t *testing.T) {
	testCases := []struct {
		name      string
		inputStr  string
		expectStr string
	}{
		{
			name:      "case 0",
			inputStr:  "not found error",
			expectStr: "notFoundError",
		},
		{
			name:      "case 1",
			inputStr:  "URL not found error",
			expectStr: "urlNotFoundError",
		},
		{
			name:      "case 2",
			inputStr:  "NOT FOUND ERROR",
			expectStr: "notFOUNDERROR",
		},
		{
			name:      "case 3",
			inputStr:  "UPPER",
			expectStr: "upper",
		},
		{
			name:      "case 4",
			inputStr:  "lower",
			expectStr: "lower",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			str := toCamelCase(tc.inputStr)
			if !reflect.DeepEqual(str, tc.expectStr) {
				t.Fatalf("str = %v, want %v", str, tc.expectStr)
			}
		})
	}
}

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
		{
			Name:           "case 6: camel case with abbreviation at the start",
			InputString:    "APINotAvailableError",
			ExpectedString: "api not available error",
		},
		{
			Name:           "case 7: camel case with abbreviation in the middle",
			InputString:    "invalidHTTPStatusError",
			ExpectedString: "invalid http status error",
		},
		{
			Name:           "case 8: camel case with abbreviation at the end",
			InputString:    "fooBarBAZ",
			ExpectedString: "foo bar baz",
		},
		{
			Name:           "case 9: with version numbers at the start",
			InputString:    "v2RouteNotReachable",
			ExpectedString: "v2 route not reachable",
		},
		{
			Name:           "case 10: with version numbers in the middle",
			InputString:    "oldV2RouteNotReachable",
			ExpectedString: "old v2 route not reachable",
		},
		{
			Name:           "case 11: with version numbers in the middle",
			InputString:    "oldV2RouteNotReachable",
			ExpectedString: "old v2 route not reachable",
		},
		{
			Name:           "case 12: with version numbers at the end does not work",
			InputString:    "statusCode200",
			ExpectedString: "status code200",
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
