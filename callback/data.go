// Package callback provides reflective CallbackData pack/unpack.
package callback

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	"pyrorhythm.dev/fn/res"
)

const maxCallbackDataLen = 64

// Callback packs and unpacks typed callback data via reflection.
type Callback[T any] struct {
	Prefix string
	Sep    byte
}

// New creates a Callback with default separator ':'.
func New[T any]() *Callback[T] {
	return &Callback[T]{Sep: ':'}
}

// Pack encodes v as callback_data (max 64 bytes).
func (c *Callback[T]) Pack(v T) (string, error) {
	prefix := c.prefixFor(v)
	parts, err := c.encodeFields(v)
	if err != nil {
		return "", err
	}
	data := prefix
	if len(parts) > 0 {
		if data != "" {
			data += string(c.sep())
		}
		data += strings.Join(parts, string(c.sep()))
	}
	if len(data) > maxCallbackDataLen {
		return "", fmt.Errorf("callback: data exceeds %d bytes", maxCallbackDataLen)
	}
	return data, nil
}

// Unpack decodes callback_data into T.
func (c *Callback[T]) Unpack(data string) res.Of[T] {
	sep := string(c.sep())
	prefix := c.Prefix
	if prefix == "" {
		var sample T
		prefix = c.prefixFor(sample)
	}
	if prefix != "" {
		if !strings.HasPrefix(data, prefix) {
			return res.Errn[T]("callback: bad prefix")
		}
		data = strings.TrimPrefix(data, prefix)
		if strings.HasPrefix(data, sep) {
			data = strings.TrimPrefix(data, sep)
		}
	}
	parts := []string{}
	if data != "" {
		parts = strings.Split(data, sep)
	}
	return c.decodeFields(parts)
}

// Filter returns a predicate that unpacks and runs pred.
func (c *Callback[T]) Filter(pred func(T) bool) th.Predicate {
	return func(_ context.Context, u telego.Update) bool {
		if u.CallbackQuery == nil {
			return false
		}
		r := c.Unpack(u.CallbackQuery.Data)
		if r.Err() != nil {
			return false
		}
		return pred(r.Val())
	}
}

// Equal returns a predicate matching packed value of v.
func (c *Callback[T]) Equal(v T) th.Predicate {
	data, err := c.Pack(v)
	if err != nil {
		return func(context.Context, telego.Update) bool { return false }
	}
	return th.CallbackDataEqual(data)
}

func (c *Callback[T]) sep() byte {
	if c.Sep == 0 {
		return ':'
	}
	return c.Sep
}

func (c *Callback[T]) prefixFor(v T) string {
	if c.Prefix != "" {
		return c.Prefix
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return ""
	}
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		if tag := rt.Field(i).Tag.Get("cb"); tag != "" {
			return tag
		}
	}
	return ""
}

func (c *Callback[T]) encodeFields(v T) ([]string, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("callback: need struct type")
	}
	prefix := c.prefixFor(v)
	parts := make([]string, 0, rv.NumField())
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		if tag := rv.Type().Field(i).Tag.Get("cb"); tag != "" && tag == prefix {
			continue
		}
		s, err := encodeValue(f)
		if err != nil {
			return nil, err
		}
		parts = append(parts, s)
	}
	return parts, nil
}

func (c *Callback[T]) decodeFields(parts []string) res.Of[T] {
	var out T
	rv := reflect.ValueOf(&out).Elem()
	if rv.Kind() != reflect.Struct {
		return res.Errn[T]("callback: need struct type")
	}
	idx := 0
	for i := 0; i < rv.NumField(); i++ {
		if tag := rv.Type().Field(i).Tag.Get("cb"); tag != "" {
			continue
		}
		if idx >= len(parts) {
			return res.Errn[T]("callback: field count mismatch")
		}
		if err := decodeValue(rv.Field(i), parts[idx]); err != nil {
			return res.Errw[T](err, "callback")
		}
		idx++
	}
	if idx != len(parts) {
		return res.Errn[T]("callback: extra fields")
	}
	return res.OKAny(out)
}

func encodeValue(f reflect.Value) (string, error) {
	switch f.Kind() {
	case reflect.String:
		s := f.String()
		if strings.ContainsRune(s, ':') {
			return "", fmt.Errorf("callback: colon in string field")
		}
		return s, nil
	case reflect.Int, reflect.Int64:
		return strconv.FormatInt(f.Int(), 10), nil
	case reflect.Bool:
		if f.Bool() {
			return "1", nil
		}
		return "0", nil
	default:
		return "", fmt.Errorf("callback: unsupported type %s", f.Kind())
	}
}

func decodeValue(f reflect.Value, s string) error {
	switch f.Kind() {
	case reflect.String:
		f.SetString(s)
		return nil
	case reflect.Int, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		f.SetInt(n)
		return nil
	case reflect.Bool:
		switch s {
		case "1", "true":
			f.SetBool(true)
		case "0", "false":
			f.SetBool(false)
		default:
			return fmt.Errorf("bad bool %q", s)
		}
		return nil
	default:
		return fmt.Errorf("unsupported type %s", f.Kind())
	}
}
