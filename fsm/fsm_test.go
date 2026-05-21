package fsm_test

import (
	"context"
	"testing"

	"github.com/mymmrac/telego"

	"pyrorhythm.dev/tgx/fsm"
)

func TestFSMTransition(t *testing.T) {
	called := false
	m := fsm.New[string, string]("start", map[fsm.StateID]fsm.Callback{
		"next": func(ctx context.Context, b *telego.Bot, u telego.Update) {
			called = true
		},
	})
	if err := m.Transition(context.Background(), 42, "next", nil, telego.Update{}); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("callback not called")
	}
	st, err := m.Current(42)
	if err != nil || st != "next" {
		t.Fatalf("state: %q err=%v", st, err)
	}
}
