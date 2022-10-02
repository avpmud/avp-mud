package avp

import (
	"fmt"
)

func CMD_QUIT(kwarg []string, client *Client) {
	client.User.RLock()
	msg := fmt.Sprintf("[INFO] %s has left the realm.\n", client.User.Name)
	client.User.RUnlock()

	client.conn.Close()
	client.mud.BroadcastAll(msg, true)
}
