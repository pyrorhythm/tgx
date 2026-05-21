package filters

import (
	"slices"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	"github.com/pyrorhythm/tgx"
	"github.com/pyrorhythm/tgx/fsm"
)

// Command matches an exact bot command (without leading slash in name).
func Command(name string) tgx.Filter {
	p := th.CommandEqual(name)
	return func(u *telego.Update) bool { return p(nil, *u) }
}

// CommandPrefix matches commands with a prefix.
func CommandPrefix(prefix string) tgx.Filter {
	p := th.CommandPrefix(prefix)
	return func(u *telego.Update) bool { return p(nil, *u) }
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
	return func(u *telego.Update) bool { return p(nil, *u) }
}

// CallbackDataPrefix matches callback data prefix.
func CallbackDataPrefix(prefix string) tgx.Filter {
	p := th.CallbackDataPrefix(prefix)
	return func(u *telego.Update) bool { return p(nil, *u) }
}

// UserID matches updates from a specific user.
func UserID(id int64) tgx.Filter {
	return func(u *telego.Update) bool {
		uid, ok := userIDFromUpdate(u)
		return ok && uid == id
	}
}

// AllowedUsers matches updates from allowed user IDs.
func AllowedUsers(ids ...int64) tgx.Filter {
	return func(u *telego.Update) bool {
		uid, ok := userIDFromUpdate(u)
		if !ok {
			return false
		}
		return slices.Contains(ids, uid)
	}
}

// State matches FSM state for the user in the update.
func State[K comparable, V any](machine *fsm.FSM[K, V], state fsm.StateID) tgx.Filter {
	return func(u *telego.Update) bool {
		uid, ok := userIDFromUpdate(u)
		if !ok {
			return false
		}
		cur, err := machine.Current(uid)
		return err == nil && cur == state
	}
}

func userIDFromUpdate(u *telego.Update) (int64, bool) {
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
