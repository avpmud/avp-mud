package avp

import (
	"github.com/avpmud/avp-mud/pkg/state"
)

var COMMANDS map[int]map[string]func(kwarg []string, client *Client)

func init() {
	// Initialize COMMANDS map: map[STATE]map[COMMAND_STRING]COMMAND_HANDLER
	COMMANDS = make(map[int]map[string]func(kwarg []string, client *Client))

	// Initialize sub maps: STATE -> map[COMMAND STRING]COMMAND_HANDLER
	COMMANDS[state.STATE_MAIN] = make(map[string]func(kwarg []string, client *Client))

	// Initialize commands
	COMMANDS[state.STATE_MAIN]["chat"] = CMD_CHAT
}
