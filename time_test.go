package main

import (
	"testing"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		input string
		want  string // expected time in UTC format
	}{
		{"Sun, 07 Jun 2020 02:00:00 +0200", "2020-06-07T00:00:00Z"},
		{"Wed, 26 Apr 2023 20:36:36 +0000", "2023-04-26T20:36:36Z"},
		{"2023-04-28T10:59:32.609-04:00", "2023-04-28T14:59:32Z"},
		{"Mon, 24 Apr 2023 04:00:00 -0000", "2023-04-24T04:00:00Z"},
		{"Sun, 16 Apr 2023 13:01:52 GMT", "2023-04-16T13:01:52Z"},
	}
	for _, tc := range tests {
		got, err := parseTime(tc.input)
		if err != nil {
			t.Errorf("parseTime(%q) returned an error: %v", tc.input, err)
		}

		if got != tc.want {
			t.Errorf("parseTime(%q) returned %v, want %v", tc.input, got, tc.want)
		}
	}
}
