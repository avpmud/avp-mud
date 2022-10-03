package user

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"sync"
)

var CWD string

func init() {
	var err error
	CWD, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

// User contains persistent user-related data.
type User struct {
	sync.RWMutex
	Name     string `json:"name"`
	Password string `json:"password"`
	Race     int    `json:"race"`
}

// Exists returns whether a User exists.
func (u *User) Exists() bool {
	u.RLock()
	defer u.RUnlock()
	if _, err := os.Stat(path.Join(CWD, "lib", "usr", u.Name)); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// Load loads a User.
func (u *User) Load() (err error) {
	u.Lock()
	defer u.Unlock()

	// Open the user file
	f, err := os.Open(path.Join(CWD, "lib", "usr", u.Name))
	if err != nil {
		return
	}
	defer f.Close()

	// Read the file
	data, err := io.ReadAll(f)
	if err != nil {
		return
	}

	// Unmarshal the data to JSON
	return json.Unmarshal(data, &u)
}

// Save saves a User.
func (u *User) Save() (err error) {
	u.RLock()
	defer u.RUnlock()

	// Open the user file
	f, err := os.Create(path.Join(CWD, "lib", "usr", u.Name))
	if err != nil {
		return
	}
	defer f.Close()

	// Marshal JSON
	data, err := json.Marshal(u)
	if err != nil {
		return
	}

	// Write data
	_, err = f.Write(data)
	return
}

// Transaction performs a concurrent-safe write-based transaction on State.
func (u *User) Transaction(f func(u *User)) {
	u.Lock()
	defer u.Unlock()
	f(u)
}
