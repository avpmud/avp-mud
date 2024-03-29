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

func CMD_TELL(kwarg []string, client *Client) {
	if len(kwarg) < 3 {
		client.Write <- "Tell yes, tell we must, but tell who what?\n"
		return
	}

	client.User.RLock()
	name := client.User.Name
	client.User.RUnlock()

	msg := fmt.Sprintf("%s tells you, '%s'\n", name, strings.Join(kwarg[2:], " "))
	client.mud.RLock()
	var found bool
	for _, cl := range client.mud.clients {
		if strings.EqualFold(kwarg[1], name) {
			cl.Write <- msg
			found = true
		}
		if found {
			break
		}
	}
	client.mud.RUnlock()
	if !found {
		client.Write <- "Nobody by that name seems to be around.\n"
	}
}

func CMD_WHO(kwarg []string, client *Client) {
	var msg string
	client.mud.RLock()
	for _, cl := range client.mud.clients {
		cl.User.RLock()
		msg += cl.User.Name + "\n"
		cl.User.RUnlock()
	}
	client.mud.RUnlock()
	client.Write <- "Online users ------------\n" + msg + "-------------------------\n"
}
