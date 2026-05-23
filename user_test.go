package tgx_test

import (
	"testing"

	"github.com/mymmrac/telego"

	"pyrorhythm.dev/tgx"
)

func TestUserID(t *testing.T) {
	u := telego.Update{Message: &telego.Message{From: &telego.User{ID: 99}}}
	if id := tgx.UserID(&u); id.Val() != 99 {
		t.Fatalf("got %v", id.Val())
	}

	u = telego.Update{GuestMessage: &telego.Message{From: &telego.User{ID: 42}}}
	if id := tgx.UserID(&u); id.Val() != 42 {
		t.Fatalf("guest message: got %v", id.Val())
	}
}
