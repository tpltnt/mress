package main

import (
	"flag"
	"fmt"
	"github.com/thoj/go-ircevent" // imported as "irc"
	"log"
	"os"
	"strconv"
)

func main() {
	configfile := flag.String("config", "config.ini", "configuration file (lower priority if other flags are defined)")
	logdest := flag.String("log", "", "destination (filename, stdout, stderr) of the log")
	ircNick := flag.String("nick", "mress", "nickname")
	ircPasswd := flag.String("passwd", "", "server/ident password")
	ircServer := flag.String("server", "", "IRC server hostname")
	ircPort := flag.Int("port", 6697, "IRC server port")
	ircChannel := flag.String("channel", "", "IRC channel to join")
	useTLS := flag.Bool("use-tls", true, "use TLS encrypted connection")
	debug := flag.Bool("debug", false, "enable debugging (+flags)")
	offlineMsgDb := flag.String("offline-msg-db", "messages.db", "filename of sqlite3 database for offline messages")
	flag.Parse()

	logchan := make(chan *log.Logger)
	go getLogger(*logdest, *configfile, logchan)
	logger := <-logchan
	if nil == logger {
		fmt.Fprint(os.Stderr, "creating logger failed")
		os.Exit(1)
	}

	// determine config values concurrently with go-routines
	// "fork" roughly according to need, choose non-default
	// flags over config file values, collect config values
	// later as needed
	nickchan := make(chan string)
	go getNick(*ircNick, *configfile, nickchan, logger)
	passwdchan := make(chan string)
	go getPassword(*ircPasswd, *configfile, passwdchan, logger)
	servchan := make(chan string)
	go getServer(*ircServer, *configfile, servchan, logger)
	portchan := make(chan int)
	go getPort(*ircPort, *configfile, portchan, logger)
	chanchan := make(chan string)
	go getChannel(*ircChannel, *configfile, chanchan, logger)
	// to disable TLS and/or use debugging should always
	// be conscious decisions and are therefore not part
	// of the config.
	offlinedbchan := make(chan string)
	go getOfflineDBfilename(*offlineMsgDb, *configfile, offlinedbchan, logger)
	// create IRC connection
	nick := <-nickchan
	irccon := irc.IRC(nick, "mress")
	if nil == irccon {
		logger.Println("creating IRC connection failed")
	} else {
		logger.Println("creating IRC connection worked")
	}
	irccon.Password = <-passwdchan
	if 0 < len(irccon.Password) {
		logger.Println("password is used")
	}
	// configure IRC connection
	if *useTLS {
		irccon.UseTLS = true
		logger.Println("using TLS encrypted connection")
	} else {
		irccon.UseTLS = false
		logger.Println("using cleartext connection")
	}
	if *debug {
		irccon.Debug = true
	}

	// connect to server
	socketstring := <-servchan + ":" + strconv.Itoa(<-portchan)
	logger.Println("connecting to " + socketstring)
	err := irccon.Connect(socketstring)
	if err != nil {
		logger.Println("connecting to server failed")
		logger.Println(err.Error())
		os.Exit(2)
	}
	logger.Println("connecting to server succeeded")

	// collect last config value needed
	channel := <-chanchan
	// add callbacks
	irccon.AddCallback("001", func(e *irc.Event) {
		logger.Println("joining " + channel)
		irccon.Join(channel)
	})

	offlmsgdb := <-offlinedbchan
	irccon.AddCallback("001", func(e *irc.Event) {
		err := initOfflineMessageDatabase(offlmsgdb)
		if err != nil {
			logger.Println(err.Error())
		}
	})
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		offlineMessengerCommand(e, irccon, nick, offlmsgdb, logger)
	})
	irccon.AddCallback("JOIN", func(e *irc.Event) {
		offlineMessengerDrone(e, irccon, offlmsgdb, nick, channel, logger)
	})
	irccon.AddCallback("353", func(e *irc.Event) {
		offlineMessengerDrone(e, irccon, offlmsgdb, nick, channel, logger)
	})

	logger.Println("starting event loop")
	irccon.Loop()
}
