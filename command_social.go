package avp

import (
	"fmt"
	"strings"
)

func CMD_CHAT(kwarg []string, client *Client) {
	if len(kwarg) < 2 {
		client.Write <- "Chat yes, chat we must, but chat what?\n"
		return
	}

	client.User.RLock()
	msg := fmt.Sprintf("[CHAT] %s: %s\n", client.User.Name, strings.Join(kwarg[1:], " "))
	client.User.RUnlock()

	client.mud.BroadcastAll(msg, true)
}
