package appkit

import (
	"fmt"
	"strings"
	"testing"
)

func StringArraysEqual(t *testing.T, a, b []string) {
	arrayToString := func(arr []string) string {
		return fmt.Sprintf("Array: [%s] len: %d, cap: %d",
			strings.Join(arr, ", "), len(arr), cap(arr),
		)
	}
	as := arrayToString(a)
	bs := arrayToString(b)
	if as != bs {
		t.Errorf("Expecting arrays to be equal:\n[%s]\n\nGot:\n[%s]\n", as, bs)
	}
}

func Test_JoinSplitArguments(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		joined string
	}{
		{"Empty input", []string{}, ""},
		{"Empty string", []string{""}, "\000"},
		{"Two empty strings", []string{"", ""}, "\000\000"},
		{"One argument", []string{"a"}, "a\000"},
		{"Two arguments", []string{"a", "b"}, "a\000b\000"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			joined := JoinArguments(tt.input)
			if joined != tt.joined {
				t.Errorf("Expecting joined string:\n[%s]\n\nGot:\n[%s]\n", tt.joined, joined)
			}

			split := SplitArguments(joined)
			StringArraysEqual(t, tt.input, split)
		})
	}
}
