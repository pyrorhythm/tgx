package fsm

import (
	"fmt"
	"sync"

	"pyrorhythm.dev/fn/res"
)

// userStateStorage is a type for default user's state storage
type userStateStorage struct {
	mu      sync.RWMutex
	Storage map[int64]StateID
}

// initialUserStateStorage creates in memory storage for user's state
func initialUserStateStorage() *userStateStorage {
	return &userStateStorage{
		Storage: make(map[int64]StateID),
	}
}

// Set sets user's state to state storage
func (u *userStateStorage) Set(userID int64, stateID StateID) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.Storage[userID] = stateID

	return nil
}

// Exists checks whether any user's state exist in state storage
func (u *userStateStorage) Exists(userID int64) res.Of[bool] {
	u.mu.RLock()
	defer u.mu.RUnlock()

	_, ok := u.Storage[userID]

	return res.OKAny(ok)
}

// Get gets user's state from state storage
func (u *userStateStorage) Get(userID int64) res.Of[StateID] {
	u.mu.RLock()
	defer u.mu.RUnlock()

	s, ok := u.Storage[userID]
	if !ok {
		return res.Errn[StateID](fmt.Sprintf("%v: userID: %d", ErrNoUserState, userID))
	}

	return res.OKAny(s)
}
