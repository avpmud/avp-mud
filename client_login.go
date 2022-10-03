package avp

import (
	"fmt"
	"strings"

	"github.com/avpmud/avp-mud/pkg/state"
	"github.com/avpmud/avp-mud/pkg/user"
)

// ProcessLogin processes input for clients in the login state.
func (c *Client) ProcessLogin(kwarg []string, st int) {
	switch st {
	case state.STATE_LOGIN_NAME:
		if len(kwarg[0]) < 4 {
			c.Write <- "That's no name!\n\nWhat's your name, soldier? "
			return
		}

		// Properly format name
		name := strings.ToUpper(string(kwarg[0][0])) + kwarg[0][1:]
		c.User.Transaction(func(u *user.User) {
			u.Name = name
		})

		// If the user exists, prompt for password
		if c.User.Exists() {
			if err := c.User.Load(); err != nil {
				c.Write <- "Sorry, looks like a corrupted user file. Aborting.\n"
				c.mud.CallbackErr(err)
				c.Quit()
				return
			}
			c.Write <- "Password: "
			c.State.Transaction(func(s *state.State) {
				s.State = state.STATE_LOGIN_PASSWORD
			})
			return
		}

		// Else, enter character creation
		c.Write <- fmt.Sprintf("Create new character %s? [y/N] ", kwarg[0])
		c.State.Transaction(func(s *state.State) {
			s.State = state.STATE_CREATE_CONFIRM
		})
	case state.STATE_LOGIN_PASSWORD:
		c.User.RLock()
		name, pass := c.User.Name, c.User.Password
		c.User.RUnlock()

		// If password is valid, login
		if CheckPasswordHash(kwarg[0], pass) {
			c.mud.BroadcastAll(fmt.Sprintf("[INFO] %s has entered the realm.\n", name), true)
			c.State.Transaction(func(s *state.State) {
				s.State = state.STATE_MAIN
			})
			return
		}

		// Else, disconnect user
		c.Write <- "Wrong password.\n"
		c.Quit()
	}
}
