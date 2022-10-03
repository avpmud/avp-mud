package avp

import (
	"fmt"
	"strings"

	"github.com/avpmud/avp-mud/pkg/constant"
	"github.com/avpmud/avp-mud/pkg/race"
	"github.com/avpmud/avp-mud/pkg/state"
	"github.com/avpmud/avp-mud/pkg/user"
)

// ProcessLogin processes input for clients in the login state.
func (c *Client) ProcessCreate(kwarg []string, st int) {
	switch st {
	case state.STATE_CREATE_CONFIRM:
		switch true {
		case strings.EqualFold(kwarg[0], "y") || strings.EqualFold(kwarg[0], "yes"):
			c.Write <- "Enter a password: "
			c.State.Transaction(func(s *state.State) {
				s.State = state.STATE_CREATE_PASSWORD
			})
		default:
			c.Write <- "What is your name, soldier? "
			c.State.Transaction(func(s *state.State) {
				s.State = state.STATE_LOGIN_NAME
			})
		}
	case state.STATE_CREATE_PASSWORD:
		if len(kwarg[0]) < 6 {
			c.Write <- "That's not a good password!\nEnter a password: "
			return
		}
		// Hash the password
		hash := HashPassword(kwarg[0])
		// Update user data
		c.User.Transaction(func(u *user.User) {
			u.Password = hash
		})
		c.State.Transaction(func(s *state.State) {
			s.State = state.STATE_CREATE_RACE
		})
		// Prompt the user to select a race
		c.Write <- constant.RACE_SELECTION
	case state.STATE_CREATE_RACE:
		switch true {
		case strings.EqualFold(kwarg[0], "a") || strings.EqualFold(kwarg[0], "alien"):
			c.User.Transaction(func(u *user.User) {
				u.Race = race.RACE_ALIEN
			})
		case strings.EqualFold(kwarg[0], "h") || strings.EqualFold(kwarg[0], "human"):
			c.User.Transaction(func(u *user.User) {
				u.Race = race.RACE_HUMAN
			})
		case strings.EqualFold(kwarg[0], "p") || strings.EqualFold(kwarg[0], "predator"):
			c.User.Transaction(func(u *user.User) {
				u.Race = race.RACE_PREDATOR
			})
		default:
			c.Write <- constant.RACE_SELECTION
			return
		}
		c.User.RLock()
		name := c.User.Name
		c.User.RUnlock()
		if err := c.User.Save(); err != nil {
			c.mud.CallbackErr(err)
			c.Write <- "error, unable to save user, " + err.Error()
		}
		c.mud.BroadcastAll(fmt.Sprintf("[INFO] %s has joined the realm.\n", name), true)
		c.State.Transaction(func(s *state.State) {
			s.State = state.STATE_MAIN
		})
	}
}
