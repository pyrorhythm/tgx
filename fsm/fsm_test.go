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
	st := m.Current(42)
	if st.Err() != nil || st.Val() != "next" {
		t.Fatalf("state: %q err=%v", st.Val(), st.Err())
	}
}

func TestDataStorageMissingKey(t *testing.T) {
	m := fsm.New[string, string]("start", nil)
	if err := m.Set(1, "k", "v"); err != nil {
		t.Fatal(err)
	}
	got := m.Get(1, "missing")
	if got.Err() == nil {
		t.Fatal("expected error for missing key")
	}
}
