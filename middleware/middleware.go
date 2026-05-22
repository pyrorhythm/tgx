// Package middleware provides common tgx middleware.
package middleware

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	"pyrorhythm.dev/fn/opt"
	"pyrorhythm.dev/tgx"
)

const maxSerializeKeys = 4096

// Recover wraps panics into errors.
func Recover() tgx.Middleware {
	return func(ctx *tgx.Ctx, next func() error) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("tgx: panic: %v", r)
			}
		}()
		return next()
	}
}

// Logger logs handled updates at debug level.
func Logger() tgx.Middleware {
	return func(ctx *tgx.Ctx, next func() error) error {
		slog.Debug("tgx update", "id", ctx.Update.UpdateID)
		return next()
	}
}

// ACL allows only listed user IDs.
func ACL(allowed map[int64]struct{}) tgx.Middleware {
	return func(ctx *tgx.Ctx, next func() error) error {
		uid := tgx.UserID(&ctx.Update)
		if !uid.Valid() {
			return nil
		}
		if _, ok := allowed[uid.Val()]; !ok {
			return nil
		}
		return next()
	}
}

// StorageKey identifies a serialization bucket.
type StorageKey struct {
	BotID    int64
	ChatID   int64
	UserID   int64
	ThreadID int64
}

// KeyFromUpdate extracts a storage key from an update.
func KeyFromUpdate(u telego.Update) opt.Of[StorageKey] {
	var k StorageKey
	switch {
	case u.Message != nil:
		k.ChatID = u.Message.Chat.ID
		if u.Message.From != nil {
			k.UserID = u.Message.From.ID
		}
		k.ThreadID = int64(u.Message.MessageThreadID)
		return opt.SomeAny(k)
	case u.CallbackQuery != nil:
		if u.CallbackQuery.Message != nil {
			k.ChatID = u.CallbackQuery.Message.GetChat().ID
		}
		k.UserID = u.CallbackQuery.From.ID
		return opt.SomeAny(k)
	case u.InlineQuery != nil:
		k.UserID = u.InlineQuery.From.ID
		return opt.SomeAny(k)
	case u.ChosenInlineResult != nil:
		k.UserID = u.ChosenInlineResult.From.ID
		return opt.SomeAny(k)
	default:
		return opt.Nil[StorageKey]()
	}
}

// PerKeySerialize serializes handlers per storage key (for FSM safety).
func PerKeySerialize(keyFn func(telego.Update) opt.Of[StorageKey]) tgx.Middleware {
	var (
		mu    sync.Mutex
		locks = make(map[StorageKey]*sync.Mutex)
		order []StorageKey
	)
	return func(ctx *tgx.Ctx, next func() error) error {
		key := keyFn(ctx.Update)
		if !key.Valid() {
			return next()
		}
		k := key.Val()

		mu.Lock()
		lk, exists := locks[k]
		if !exists {
			if len(order) >= maxSerializeKeys {
				evict := order[0]
				order = order[1:]
				delete(locks, evict)
			}
			lk = &sync.Mutex{}
			locks[k] = lk
			order = append(order, k)
		}
		mu.Unlock()

		lk.Lock()
		defer lk.Unlock()
		return next()
	}
}

// TelegoRecover applies telego's panic recovery on the dispatcher root.
func TelegoRecover() th.Handler {
	return th.PanicRecovery()
}
