package middleware_test

import (
	"testing"

	"github.com/mymmrac/telego"

	"pyrorhythm.dev/tgx"
	"pyrorhythm.dev/tgx/middleware"
)

func TestACL(t *testing.T) {
	allowed := map[int64]struct{}{1: {}}
	mw := middleware.ACL(allowed)

	var called bool
	ctx := &tgx.Ctx{Update: telego.Update{Message: &telego.Message{From: &telego.User{ID: 1}}}}
	err := mw(ctx, func() error { called = true; return nil })
	if err != nil || !called {
		t.Fatalf("allowed user: err=%v called=%v", err, called)
	}

	called = false
	ctx.Update = telego.Update{Message: &telego.Message{From: &telego.User{ID: 2}}}
	_ = mw(ctx, func() error { called = true; return nil })
	if called {
		t.Fatal("denied user should not run handler")
	}
}
