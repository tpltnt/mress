package main

/* Introduce bot to first-time users.
 * Track nicks as already seen to avoid muliple introductions etc. */
import (
	"github.com/thoj/go-ircevent" // imported as "irc"
	"strings"
	"time"
)

// Callback on JOIN to introduce bot functionality.
func runIntroduction(e *irc.Event, irc *irc.Connection, user, channel string) {
	// ignore OTR
	if 0 == strings.Index(e.Message(), "?OTR") {
		return
	}

	irc.Privmsg(e.Nick, "Hi "+e.Nick+"\n")
	time.Sleep(1 * time.Second)
	irc.Privmsg(e.Nick, "I am "+user+", the bot for this channel.\n")
	time.Sleep(1 * time.Second)
	irc.Privmsg(e.Nick, "Currently I can enable offline messaging for you.\n")
	time.Sleep(2 * time.Second)
	irc.Privmsg(e.Nick, "To leave a message for another user just type the following:\n")
	time.Sleep(2 * time.Second)
	irc.Privmsg(e.Nick, "/msg "+user+" tell <username>: <your message>\n")
	time.Sleep(3 * time.Second)
	irc.Privmsg(e.Nick, "I will then deliver it to the user you mentioned as soon as (s)he joins the channel.\n")
}

// Mark a nick as seen (and log to a database).
func markAsSeen(dbconfig MressDbConfig, user string) error {
	return fmt.Errorf("not implemented yet")
	return nil
}
