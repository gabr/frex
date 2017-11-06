package main

import (
	"testing"
	"strings"
)

type argsSlice []string

func (s argsSlice) String() string {
	if s == nil {
		return "[<<nil>>]"
	}

	res := "["
	last := len(s) - 1
	for i, v := range s {
		res += "\"" + v + "\""
		if i != last {
			res += ", "
		}
	}

	return res + "]"
}

func TestParseArgsValidation(t *testing.T) {
	var testCases = []struct {
		args   argsSlice
		errMsg string
	}{
		{nil, "nil"},

		{argsSlice{},         "Not enought arguments"},
		{argsSlice{"a"},      "Not enought arguments"},
		{argsSlice{"a", "a"}, "Not enought arguments"},

		{argsSlice{"(\\d", "a", "a"}, "error"},

		// not safe test, because fakeDir and fakeFile can exist
		{argsSlice{"\\d", "a", "fakeDir\\fakeFile"}, "not found"},
	}

	for _, tc := range testCases {
		_, err := parseArgs(tc.args)

		if err == nil {
			t.Errorf("parseArgs(%v) returned nil error, expected error with: %q",
				tc.args, tc.errMsg)
		} else if false == strings.Contains(err.Error(), tc.errMsg) {
			t.Errorf("parseArgs(%v) returnet wrong error message\n" +
				"Got: %s\nExpected to contain:%s\n",
				tc.args, err, tc.errMsg)
		}
	}
}
