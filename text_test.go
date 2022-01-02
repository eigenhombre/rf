package main

import (
	"testing"
)

func TestWrapping(t *testing.T) {
	var tests = []struct {
		input  string
		column int
		want   string
	}{
		{"", 80, ""},
		{"a", 80, "a"},
		{"123", 3, "123"},
		{"1234", 3, "\n1234"}, // Improve this?
		{"123 567", 3, "123\n567"},
		{"123 567 901", 7, "123 567\n901"},
		{"123  67", 3, "123\n67"},
		{"1 3 5 7", 1, "1\n3\n5\n7"},
	}
	for _, test := range tests {
		if got := wrapText(test.input, test.column); got != test.want {
			t.Errorf("fill(%q, %d) = %q", test.input, test.column, got)
		}
	}
}
