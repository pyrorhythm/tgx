package di_test

import (
	"context"
	"testing"

	"pyrorhythm.dev/tgx/di"
)

func TestPutGet(t *testing.T) {
	ctx := di.Put(context.Background(), 42)
	v := di.Get[int](ctx)
	if !v.Valid() || v.Val() != 42 {
		t.Fatalf("got %v", v)
	}
}
