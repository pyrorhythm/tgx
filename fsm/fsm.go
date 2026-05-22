package fsm

import (
	"context"
	"maps"

	"github.com/mymmrac/telego"
	"github.com/pkg/errors"

	"pyrorhythm.dev/fn/res"
)

// StateID is a type for state identifier
type StateID string

// Callback is a function that will be called on state transition
type Callback func(ctx context.Context, b *telego.Bot, u telego.Update)

// FSM is a finite state machine
type FSM[K comparable, V any] struct {
	initialStateID StateID
	callbacks      map[StateID]Callback
	userStates     UserStateStorage
	storage        DataStorage[K, V]
}

// UserStateStorage is an interface for user state storage
type UserStateStorage interface {
	Set(userID int64, stateID StateID) error
	Exists(userID int64) res.Of[bool]
	Get(userID int64) res.Of[StateID]
}

// DataStorage is an interface for data storage
type DataStorage[K comparable, V any] interface {
	Set(userID int64, key K, value V) error
	Get(userID int64, key K) res.Of[V]
	Delete(userID int64, key K) error
}

// New creates a new FSM
func New[K comparable, V any](initialStateName StateID, callbacks map[StateID]Callback, opts ...Option[K, V]) *FSM[K, V] {
	s := &FSM[K, V]{
		initialStateID: initialStateName,
		callbacks:      callbacks,
		userStates:     initialUserStateStorage(),
		storage:        initialDataStorage[K, V](),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// AddCallback adds a callback for a state
func (f *FSM[K, V]) AddCallback(stateID StateID, callback Callback) {
	f.callbacks[stateID] = callback
}

// AddCallbacks adds callbacks for states
func (f *FSM[K, V]) AddCallbacks(cb map[StateID]Callback) {
	maps.Copy(f.callbacks, cb)
}

// Transition transitions the user to a new state
func (f *FSM[K, V]) Transition(
	ctx context.Context,
	userID int64,
	stateID StateID,
	b *telego.Bot,
	u telego.Update,
) error {
	if err := f.userStates.Set(userID, stateID); err != nil {
		return errors.Wrap(err, "failed to set user state")
	}

	if cb, ok := f.callbacks[stateID]; ok {
		cb(ctx, b, u)
	}

	return nil
}

// Current returns the current state of the user
func (f *FSM[K, V]) Current(userID int64) res.Of[StateID] {
	ok, err := f.userStates.Exists(userID).Unpack()
	if err != nil {
		return res.Errw[StateID](err, "failed to check user state")
	}
	if !ok {
		if err := f.userStates.Set(userID, f.initialStateID); err != nil {
			return res.Errw[StateID](err, "failed to set user state to initial")
		}
		return res.OKAny(f.initialStateID)
	}

	state, err := f.userStates.Get(userID).Unpack()
	if err != nil {
		return res.Errw[StateID](err, "failed to get user state")
	}

	return res.OKAny(state)
}

// Reset resets the state of the user to the initial state
func (f *FSM[K, V]) Reset(userID int64) error {
	return f.userStates.Set(userID, f.initialStateID)
}

// Set sets a value to data storage by userID and comparable
func (f *FSM[K, V]) Set(userID int64, key K, value V) error {
	return errors.Wrap(f.storage.Set(userID, key, value), "failed to set user data")
}

// Get gets a value from data storage by userID and comparable
func (f *FSM[K, V]) Get(userID int64, key K) res.Of[V] {
	r := f.storage.Get(userID, key)
	if r.Err() != nil {
		return res.Errw[V](r.Err(), "failed to get user data")
	}
	return r
}

// Delete deletes a value from data storage by userID and comparable
func (f *FSM[K, V]) Delete(userID int64, key K) error {
	return errors.Wrap(f.storage.Delete(userID, key), "failed to delete user data")
}
