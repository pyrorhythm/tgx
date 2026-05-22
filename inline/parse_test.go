package inline_test

import (
	"testing"

	"pyrorhythm.dev/tgx/inline"
)

func TestParser_Parse(t *testing.T) {
	p := inline.NewParser(
		inline.CommandRule{Name: "search", MinArgsLen: 2},
		inline.CommandRule{Name: "search", MinArgsLen: 2},
		inline.CommandRule{Name: "now-playing", MinArgsLen: 0},
	)

	tests := []struct {
		in     string
		cmd    string
		args   string
		wantOK bool
	}{
		{"search drake", "search", "drake", true},
		{"search a", "", "", false},
		{"now-playing", "now-playing", "", true},
		{"drake gods plan", "search", "drake gods plan", true},
		{"", "", "", false},
		{"a", "", "", false},
	}
	for _, tc := range tests {
		got, ok := p.Parse(tc.in)
		if ok != tc.wantOK {
			t.Fatalf("%q: ok=%v want %v", tc.in, ok, tc.wantOK)
		}
		if !ok {
			continue
		}
		if got.Command != tc.cmd || got.Args != tc.args {
			t.Fatalf("%q: got cmd=%q args=%q want cmd=%q args=%q", tc.in, got.Command, got.Args, tc.cmd, tc.args)
		}
	}
}
