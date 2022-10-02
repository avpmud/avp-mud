package state

import "sync"

const (
	STATE_LOGIN = iota
	STATE_MAIN
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
	return !(s.State == STATE_LOGIN)
}

// Transaction performs a concurrent-safe write-based transaction on State.
func (s *State) Transaction(f func(s *State)) {
	s.Lock()
	defer s.Unlock()
	f(s)
}
