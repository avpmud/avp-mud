package user

import "sync"

// User contains persistent user-related data.
type User struct {
	sync.RWMutex
	Name string `json:"name"`
}

// Transaction performs a concurrent-safe write-based transaction on State.
func (u *User) Transaction(f func(u *User)) {
	u.Lock()
	defer u.Unlock()
	f(u)
}
