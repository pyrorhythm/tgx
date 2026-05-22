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
		got := p.Parse(tc.in)
		if got.Valid() != tc.wantOK {
			t.Fatalf("%q: valid=%v want %v", tc.in, got.Valid(), tc.wantOK)
		}
		if !got.Valid() {
			continue
		}
		v := got.Val()
		if v.Command != tc.cmd || v.Args != tc.args {
			t.Fatalf("%q: got cmd=%q args=%q want cmd=%q args=%q", tc.in, v.Command, v.Args, tc.cmd, tc.args)
		}
	}
}
