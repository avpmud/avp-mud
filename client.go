package avp

import (
	"bufio"
	"fmt"
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
	c.conn = conn
	c.mud = mud
	c.State = new(state.State)
	c.State.Transaction(func(s *state.State) {
		s.State = state.STATE_LOGIN_NAME
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
					mud.CallbackErr(err)
					c.Quit()
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
				c.Quit()
				return
			}
			if err = readwriter.Flush(); err != nil {
				mud.CallbackErr(err)
				c.Quit()
				return
			}
		}
	}()

	// Send the Welcome Message
	c.Write <- constant.WELCOME

	return c
}

// Process processes commands.
func (c *Client) Process(msg string) {
	kwarg := strings.Split(msg, " ")

	c.State.RLock()
	st := c.State.State
	c.State.RUnlock()

	switch st {
	case state.STATE_LOGIN_NAME, state.STATE_LOGIN_PASSWORD:
		c.ProcessLogin(kwarg, st)
	case state.STATE_CREATE_CONFIRM, state.STATE_CREATE_PASSWORD, state.STATE_CREATE_RACE:
		c.ProcessCreate(kwarg, st)
	case state.STATE_MAIN, state.STATE_MAIN_COMBAT:
		handler, present := COMMANDS[st][strings.ToLower(kwarg[0])]
		if !present {
			c.Write <- "Huh?\n"
			return
		}
		handler(kwarg, c)
	default:
		c.Write <- fmt.Sprintf("error, invalid state: %v\n", st)
		c.mud.CallbackErr(fmt.Errorf("error, invalid state: %v", st))
	}
}

// Quit quits from the MUD.
func (c *Client) Quit() {
	c.mud.RLock()
	clients := c.mud.clients
	c.mud.RUnlock()

	// Delete this client from the MUD
	var nc []*Client
	for _, cl := range clients {
		if cl.conn != c.conn {
			nc = append(nc, cl)
		}
	}

	c.mud.Lock()
	c.mud.clients = nc
	c.mud.Unlock()

	c.conn.Close()
}
