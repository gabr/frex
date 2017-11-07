package main

import (
	"strings"
	"testing"
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

		{argsSlice{}, "Not enought arguments"},
		{argsSlice{"a"}, "Not enought arguments"},
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
			t.Errorf("parseArgs(%v) returnet wrong error message\n"+
				"Got: %s\nExpected to contain:%s\n",
				tc.args, err, tc.errMsg)
		}
	}
}

func TestFindOutNewLinePos(t *testing.T) {
	var testCases = []struct {
		s      string
		offset int

		expectedPos int
		expectedSeq string
	}{
		{"", 0, -1, ""},
		{"", 1, -2, ""},
		{"", 10, -2, ""},
		{"", -1, -3, ""},
		{"", -100, -3, ""},

		{"a", 0, -1, ""},
		{"a", 1, -1, ""},
		{"a", 2, -2, ""},
		{"a", -1, -3, ""},
		{"a", -2, -3, ""},

		{"abcd", 0, -1, ""},
		{"abcd", 4, -1, ""},
		{"abcd", 5, -2, ""},
		{"abcd", -1, -3, ""},
		{"abcd", -2, -3, ""},

		{"a\x0A", 0, 1, "\x0A"},
		{"ab\x0A", 0, 2, "\x0A"},
		{"abc\x0A", 0, 3, "\x0A"},

		{"a\x0A", 1, 0, "\x0A"},
		{"ab\x0A", 3, -1, ""},
		{"abc\x0A", 5, -2, ""},

		{"abc\x0Adef\x0D\x0Agh", 0, 3, "\x0A"},
		{"abc\x0Adef\x0D\x0Agh", 3, 0, "\x0A"},
		{"abc\x0Adef\x0D\x0Agh", 4, 3, "\x0D\x0A"},
	}

	for _, tc := range testCases {
		pos, seq := findOutNewLinePos(tc.s, tc.offset)

		if pos != tc.expectedPos {
			t.Errorf("findOutNewLinePos(%q, %d) "+
				"returned pos: %d, expected: %d",
				tc.s, tc.offset, pos, tc.expectedPos)
		}

		if seq != tc.expectedSeq {
			t.Errorf("findOutNewLinePos(%q, %d) "+
				"returned seq: %q, expected: %q",
				tc.s, tc.offset, seq, tc.expectedSeq)
		}
	}
}
