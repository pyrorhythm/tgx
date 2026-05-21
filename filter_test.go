package tgx_test

import (
	"testing"

	"github.com/mymmrac/telego"

	"pyrorhythm.dev/tgx"
)

func TestFilterCombinators(t *testing.T) {
	msg := "hello"
	u := telego.Update{Message: &telego.Message{Text: msg}}

	hasHello := func(u *telego.Update) bool {
		return u.Message != nil && u.Message.Text == "hello"
	}
	hasBye := func(u *telego.Update) bool {
		return u.Message != nil && u.Message.Text == "bye"
	}

	if !tgx.And(hasHello)(&u) {
		t.Fatal("And should match")
	}
	if tgx.And(hasHello, hasBye)(&u) {
		t.Fatal("And should not match")
	}
	if !tgx.Or(hasHello, hasBye)(&u) {
		t.Fatal("Or should match")
	}
	if tgx.Not(hasHello)(&u) {
		t.Fatal("Not should invert")
	}
}
