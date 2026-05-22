package filters

import (
	"context"
	"slices"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	"pyrorhythm.dev/tgx"
	"pyrorhythm.dev/tgx/fsm"
)

// Command matches an exact bot command (without leading slash in name).
func Command(name string) tgx.Filter {
	p := th.CommandEqual(name)
	return func(u *telego.Update) bool {
		return p(context.Background(), *u)
	}
}

// CommandPrefix matches commands with a prefix.
func CommandPrefix(prefix string) tgx.Filter {
	p := th.CommandPrefix(prefix)
	return func(u *telego.Update) bool {
		return p(context.Background(), *u)
	}
}

// TextEquals matches message text.
func TextEquals(text string) tgx.Filter {
	return func(u *telego.Update) bool {
		if u.Message == nil {
			return false
		}
		return u.Message.Text == text
	}
}

// TextContains matches substring in message text.
func TextContains(sub string) tgx.Filter {
	return func(u *telego.Update) bool {
		if u.Message == nil {
			return false
		}
		return strings.Contains(u.Message.Text, sub)
	}
}

// CallbackDataEqual matches callback data exactly.
func CallbackDataEqual(data string) tgx.Filter {
	p := th.CallbackDataEqual(data)
	return func(u *telego.Update) bool {
		return p(context.Background(), *u)
	}
}

// CallbackDataPrefix matches callback data prefix.
func CallbackDataPrefix(prefix string) tgx.Filter {
	p := th.CallbackDataPrefix(prefix)
	return func(u *telego.Update) bool {
		return p(context.Background(), *u)
	}
}

// UserID matches updates from a specific user.
func UserID(id int64) tgx.Filter {
	return func(u *telego.Update) bool {
		uid := tgx.UserID(u)
		return uid.Valid() && uid.Val() == id
	}
}

// AllowedUsers matches updates from allowed user IDs.
func AllowedUsers(ids ...int64) tgx.Filter {
	return func(u *telego.Update) bool {
		uid := tgx.UserID(u)
		if !uid.Valid() {
			return false
		}
		return slices.Contains(ids, uid.Val())
	}
}

// State matches FSM state for the user in the update.
func State[K comparable, V any](machine *fsm.FSM[K, V], state fsm.StateID) tgx.Filter {
	return func(u *telego.Update) bool {
		uid := tgx.UserID(u)
		if !uid.Valid() {
			return false
		}
		cur := machine.Current(uid.Val())
		return cur.OK() && cur.Val() == state
	}
}
