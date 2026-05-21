// Package middleware provides common tgx middleware.
package middleware

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	"pyrorhythm.dev/tgx"
)

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
		uid, ok := userID(&ctx.Update)
		if !ok {
			return nil
		}
		if _, ok := allowed[uid]; !ok {
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
func KeyFromUpdate(u telego.Update) (StorageKey, bool) {
	var k StorageKey
	switch {
	case u.Message != nil:
		k.ChatID = u.Message.Chat.ID
		if u.Message.From != nil {
			k.UserID = u.Message.From.ID
		}
		k.ThreadID = int64(u.Message.MessageThreadID)
		return k, true
	case u.CallbackQuery != nil:
		if u.CallbackQuery.Message != nil {
			k.ChatID = u.CallbackQuery.Message.GetChat().ID
		}
		k.UserID = u.CallbackQuery.From.ID
		return k, true
	case u.InlineQuery != nil:
		k.UserID = u.InlineQuery.From.ID
		return k, true
	case u.ChosenInlineResult != nil:
		k.UserID = u.ChosenInlineResult.From.ID
		return k, true
	default:
		return k, false
	}
}

// PerKeySerialize serializes handlers per storage key (for FSM safety).
func PerKeySerialize(keyFn func(telego.Update) (StorageKey, bool)) tgx.Middleware {
	var (
		mu    sync.Mutex
		locks = make(map[StorageKey]*sync.Mutex)
	)
	return func(ctx *tgx.Ctx, next func() error) error {
		key, ok := keyFn(ctx.Update)
		if !ok {
			return next()
		}
		mu.Lock()
		lk, exists := locks[key]
		if !exists {
			lk = &sync.Mutex{}
			locks[key] = lk
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

func userID(u *telego.Update) (int64, bool) {
	switch {
	case u.Message != nil && u.Message.From != nil:
		return u.Message.From.ID, true
	case u.CallbackQuery != nil:
		return u.CallbackQuery.From.ID, true
	case u.InlineQuery != nil:
		return u.InlineQuery.From.ID, true
	case u.ChosenInlineResult != nil:
		return u.ChosenInlineResult.From.ID, true
	default:
		return 0, false
	}
}
