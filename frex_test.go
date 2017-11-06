package main

import "testing"

func ParseArgsValidationTest(t *testing.T) {
	var testCases = []struct {
		args   []string
		errMsg string
	}{
		{nil, "nil"}

		{{},         "Not enought arguments"},
		{{"a"},      "Not enought arguments"},
		{{"a", "a"}, "Not enought arguments"},

		{{"(\d", "a", "a"}, "error"},

		// not safe test, because fakeDir and fakeFile can exist
		{{"\d", "a", "fakeDir\fakeFile"}, "not found"},
	}

	for _, tc := range testCases {
		args, err := parseArgs(tc.args)

		if got != tc.want {
			t.Errorf("FindChar(%q) == %q, want %q", tc.s, got, tc.want)
		}
	}
}
