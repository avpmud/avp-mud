package main

import (
	"log"

	"github.com/avpmud/avp-mud"
)

var mud = new(avp.MUD)

func main() {
	for {
		if err := mud.ListenAndServe("0.0.0.0:8080"); err != nil {
			log.Println(err)
		} else {
			return
		}
	}
}
