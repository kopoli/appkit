package appkit

import (
	"flag"
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

type Opt struct {
	flags []string
	value string
	desc  string
}

func Test_OptionsHelp(t *testing.T) {
	tests := []struct {
		name    string
		options []Opt
		out     string
	}{
		{"Empty", []Opt{}, ""},
		{"Single argument", []Opt{{[]string{"a"}, "", "a"}}, "  -a  a\n"},
		{"Two arguments",
			[]Opt{
				{[]string{"a"}, "", "a"},
				{[]string{"b"}, "", "b"},
			},
			"  -a  a\n  -b  b\n",
		},
		{"Two arguments, one with value",
			[]Opt{
				{[]string{"a"}, "val", "a"},
				{[]string{"b"}, "", "b"},
			},
			"  -a string  a (default \"val\")\n" +
				"  -b         b\n",
		},
		{"Long and short option",
			[]Opt{
				{[]string{"a"}, "", "a"},
				{[]string{"aa"}, "", "a"},
			},
			"  -a, -aa  a\n",
		},
		{"Multiline usage string",
			[]Opt{
				{[]string{"a"}, "", "a\nb"},
				{[]string{"aa"}, "", "a\nb"},
			},
			"  -a, -aa  a\n" +
				"           b\n",
		},

		{"Long and short options",
			[]Opt{
				{[]string{"a"}, "", "a"},
				{[]string{"aa"}, "", "a"},
				{[]string{"first"}, "", "b"},
			},
			"  -a, -aa  a\n" +
				"  -first   b\n",
		},
		{"Long and short options with longer usage",
			[]Opt{
				{[]string{"a"}, "", "a"},
				{[]string{"aa"}, "", "a"},
				{[]string{"first"}, "", "b\nc"},
			},
			"  -a, -aa  a\n" +
				"  -first   b\n" +
				"           c\n",
		},
		{"Long and short options with longer usage on both",
			[]Opt{
				{[]string{"a"}, "", "a\nb"},
				{[]string{"aa"}, "", "a\nb"},
				{[]string{"first"}, "", "b\nc"},
			},
			"  -a, -aa  a\n" +
				"           b\n" +
				"  -first   b\n" +
				"           c\n",
		},
		{"Two options with long and short",
			[]Opt{
				{[]string{"a"}, "", "a"},
				{[]string{"aa"}, "", "a"},
				{[]string{"first"}, "", "b"},
				{[]string{"f"}, "", "b"},
			},
			"  -a, -aa     a\n" +
				"  -f, -first  b\n",
		},
		{"Two options with long and short with value",
			[]Opt{
				{[]string{"a"}, "val", "a"},
				{[]string{"aa"}, "val", "a"},
				{[]string{"first"}, "", "b"},
				{[]string{"f"}, "", "b"},
			},
			"  -a, -aa string  a (default \"val\")\n" +
				"  -f, -first      b\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			flags := flag.NewFlagSet("test", flag.ContinueOnError)
			for _, opt := range tt.options {
				for _, f := range opt.flags {
					if opt.value == "" {
						_ = flags.Bool(f, false, opt.desc)
					} else {
						_ = flags.String(f, opt.value, opt.desc)
					}
				}
			}
			s := OptionsHelp(flags)
			if s != tt.out {
				t.Errorf("Expecting string:\n%s\nGot string:\n%s\n", tt.out, s)
			}
		})
	}
}
