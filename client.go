package avp

import (
	"bufio"
	"log"
	"net"

	"github.com/avpmud/avp-mud/pkg/constants"
	"github.com/avpmud/avp-mud/pkg/input"
)

// Client represents
type Client struct {
	conn  net.Conn
	mud   *MUD
	Write chan string
}

// ListenAndServe concurrently handles read/writes for a Client.
func (c *Client) ListenAndServe(conn net.Conn, mud *MUD) *Client {
	c.mud = mud
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
			log.Println("sending message,", msg)
			if _, err = readwriter.WriteString(msg); err != nil {
				mud.CallbackErr(err)
				return
			}
			if err = readwriter.Flush(); err != nil {
				mud.CallbackErr(err)
				return
			}
		}
	}()

	// Send the Welcome Message
	c.Write <- constants.WELCOME

	return c
}

func (c *Client) Process(msg string) {
	log.Println("Processing", msg)
	c.mud.RLock()
	for _, client := range c.mud.clients {
		client.Write <- "INPUT RECEIVED"
	}
}
