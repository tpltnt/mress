package main

import (
	//	"log"
	//	"os"
	//	"strconv"
	"testing"
)

func Test_markAsSeen_0(t *testing.T) {
	irce := irc.Event{}
	irccon := irc.Connection{}
	user := "testuser"
	channel := "#testchannel"
	err := runIntroduction(irce, irccon, user, channel)
	if err != nil {
		t.Error(err.Error())
	}
}
