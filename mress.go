package main

import (
	"flag"
	"fmt"
	"github.com/thoj/go-ircevent" // imported as "irc"
)

func main() {
	fmt.Println("starting up ...")
	configfile := flag.String("config", "", "configuration file")
	flag.Parse()
	if nil == configfile {
		fmt.Println("no config file given")
	}
	ircobj := irc.IRC("<nick>", "<user>") //Create new ircobj
	if nil == ircobj {
		fmt.Println("no IRC object")
	}
}
