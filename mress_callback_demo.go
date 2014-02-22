package main

import (
	"github.com/thoj/go-ircevent" // imported as "irc"
	"strings"
	"time"
)

// The banana demo for event handling channel vs. direct message
func bananaTest(e *irc.Event, irc *irc.Connection, user, channel string) {
	time.Sleep(1 * time.Second)
	// ignore OTR
	if 0 == strings.Index(e.Message(), "?OTR") {
		return
	}
	if user == e.Arguments[0] {
		irc.Privmsg(e.Nick, "I'm not actually a banana, i am parrot!\n")
		time.Sleep(1 * time.Second)
		irc.Privmsg(e.Nick, "\""+e.Message()+"\"")
		time.Sleep(2 * time.Second)
		irc.Privmsg(e.Nick, "see ?\n")
	}
	if channel == e.Arguments[0] {
		if 0 != strings.Index(e.Message(), "mress:") {
			return
		}
		irc.Privmsg(channel, "I'm a banana!\n")
	}
}
