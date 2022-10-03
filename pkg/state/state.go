package state

import "sync"

const (
	STATE_LOGIN_NAME = iota
	STATE_LOGIN_PASSWORD
	STATE_MAIN
	STATE_MAIN_COMBAT
	STATE_CREATE_CONFIRM
	STATE_CREATE_PASSWORD
	STATE_CREATE_RACE
)

// State represents runtime user-related data.
type State struct {
	sync.RWMutex
	State int
}

// IsLoggedIn returns whether the user is logged in or not.
func (s *State) IsLoggedIn() bool {
	s.RLock()
	defer s.RUnlock()
	return (s.State == STATE_MAIN || s.State == STATE_MAIN_COMBAT)
}

// Transaction performs a concurrent-safe write-based transaction on State.
func (s *State) Transaction(f func(s *State)) {
	s.Lock()
	defer s.Unlock()
	f(s)
}
