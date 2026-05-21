package callback_test

import (
	"strings"
	"testing"

	"github.com/pyrorhythm/tgx/callback"
)

type adminAction struct {
	Action string `cb:"adm"`
	UserID int64
}

func TestCallbackRoundTrip(t *testing.T) {
	cb := callback.New[adminAction]()
	v := adminAction{Action: "ban", UserID: 123456789}
	data, err := cb.Pack(v)
	if err != nil {
		t.Fatal(err)
	}
	got := cb.Unpack(data)
	if got.Err() != nil {
		t.Fatal(got.Err())
	}
	if g := got.Val(); g.UserID != v.UserID {
		t.Fatalf("user id: got %d want %d", g.UserID, v.UserID)
	}
}

type longPayload struct {
	Note string `cb:"lng"`
	Pad  string
}

func TestCallbackOverflow(t *testing.T) {
	cb := callback.New[longPayload]()
	_, err := cb.Pack(longPayload{Pad: strings.Repeat("x", 70)})
	if err == nil {
		t.Fatal("expected overflow error")
	}
}

func TestCallbackBadPrefix(t *testing.T) {
	cb := callback.New[adminAction]()
	r := cb.Unpack("nope:1:2")
	if r.Err() == nil {
		t.Fatal("expected error")
	}
}
