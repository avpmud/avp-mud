package avp

import (
	"bufio"
	"net"
	"strings"

	"github.com/avpmud/avp-mud/pkg/constant"
	"github.com/avpmud/avp-mud/pkg/input"
	"github.com/avpmud/avp-mud/pkg/state"
	"github.com/avpmud/avp-mud/pkg/user"
)

// Client represents
type Client struct {
	conn  net.Conn
	mud   *MUD
	State *state.State
	User  *user.User
	Write chan string
}

// ListenAndServe concurrently handles read/writes for a Client.
func (c *Client) ListenAndServe(conn net.Conn, mud *MUD) *Client {
	// Initial configuration
	c.mud = mud
	c.State = new(state.State)
	c.State.Transaction(func(s *state.State) {
		s.State = state.STATE_LOGIN
	})
	c.User = new(user.User)
	c.Write = make(chan string)

	// Read/Write Loops
	go func() {
		readwriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		// Read Loops
		go func() {
			var msg string
			var err error
			for {
				msg, err = readwriter.ReadString('\n')
				if err != nil {
					c.conn.Close()
					mud.CallbackErr(err)
					return
				}
				c.Process(input.Sanitize(msg))
			}
		}()
		// Write Loop
		var msg string
		var err error
		for {
			msg = <-c.Write
			if _, err = readwriter.WriteString(msg); err != nil {
				mud.CallbackErr(err)
				return
			}
			if err = readwriter.Flush(); err != nil {
				c.conn.Close()
				mud.CallbackErr(err)
				return
			}
		}
	}()

	// Send the Welcome Message
	c.Write <- constant.WELCOME

	return c
}

func (c *Client) Process(msg string) {
	kwarg := strings.Split(msg, " ")

	if !c.State.IsLoggedIn() {
		if len(kwarg[0]) < 4 {
			c.Write <- "That's no name!\n\nWhat's your name, soldier? "
			return
		}
		c.User.Transaction(func(u *user.User) {
			u.Name = strings.ToUpper(string(kwarg[0][0])) + kwarg[0][1:]
		})
		c.State.Transaction(func(s *state.State) {
			s.State = state.STATE_MAIN
		})
		return
	}

	c.State.RLock()
	state := c.State.State
	c.State.RUnlock()

	handler, present := COMMANDS[state][strings.ToLower(kwarg[0])]
	if !present {
		c.Write <- "Huh?\n"
		return
	}

	handler(kwarg, c)
}
