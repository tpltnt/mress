package main

import (
	"flag"
	"fmt"
	"github.com/thoj/go-ircevent" // imported as "irc"
	"io"
	"log"
	"os"
	"strconv"
)

// Create a Logger which logs to the given destination
// Valid destinations are files (+path), stdout and stderr
func createLogger(destination *string) *log.Logger {
	var logdest io.Writer = nil
	var logfile *os.File = nil
	var logger *log.Logger = nil
	var err error
	if len(*destination) > 0 {
		switch *destination {
		case "stdout":
			logdest = os.Stdout
		case "stderr":
			logdest = os.Stderr
		default:
			// assuming the logfile already exists
			logfile, err = os.OpenFile(*destination, os.O_WRONLY|os.O_APPEND, 0644)
			if nil != err {
				// it didn't, so create a new one
				logfile, err = os.Create(*destination)
				if nil != err {
					fmt.Fprintln(os.Stderr, err.Error())
					return nil
				}
				err = logfile.Chmod(0644)
			}
		}
	} else {
		logfile, err = os.OpenFile("/dev/null", os.O_RDWR, 666)
	}
	if nil != err {
		fmt.Fprint(os.Stderr, "opening logging destination failed\n")
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		return nil
	}
	if nil != logfile {
		logger = log.New(logfile, "[mress] ", log.LstdFlags)
	} else {
		logger = log.New(logdest, "[mress] ", log.LstdFlags)
	}
	return logger
}

func main() {
	configfile := flag.String("config", "", "configuration file (lower priority if other flags are defined)")
	logdest := flag.String("log", "", "destination (filename, stdout, stderr) of the log")
	nick := flag.String("nick", "mress", "nickname")
	passwd := flag.String("passwd", "", "server/ident password")
	ircServer := flag.String("server", "irc.freenode.net", "IRC server hostname")
	ircPort := flag.Int("port", 6697, "IRC server port")
	ircChannel := flag.String("channel", "", "IRC channel to join")
	useTLS := flag.Bool("use-tls", true, "use TLS encrypted connection")
	debug := flag.Bool("debug", false, "enable debugging (+flags)")
	flag.Parse()

	logger := createLogger(logdest)
	if nil == logger {
		fmt.Fprint(os.Stderr, "creating logger failed")
		os.Exit(1)
	}

	if 0 == len(*configfile) {
		fmt.Println("no config file given, using defaults")
	} else {
		fmt.Fprintln(os.Stderr, "configuration file parsing not implemented yet")
		os.Exit(1)
	}

	if len(*ircChannel) == 0 {
		logger.Println("no channel given to join")
		os.Exit(1)
	}

	irccon := irc.IRC(*nick, "mress")
	if nil == irccon {
		logger.Println("creating IRC connection failed")
	} else {
		logger.Println("creating IRC connection worked")
	}
	// configure IRC connection
	if *useTLS {
		irccon.UseTLS = true
		logger.Println("using TLS encrypted connection")
	} else {
		irccon.UseTLS = false
		logger.Println("using cleartext connection")
	}
	irccon.Password = *passwd
	if 0 < len(*passwd) {
		logger.Println("password is used")
	}
	if *debug {
		irccon.Debug = true
	}

	// connect to server
	socketstring := *ircServer + ":" + strconv.Itoa(*ircPort)
	logger.Println("connecting to " + socketstring)
	err := irccon.Connect(socketstring)
	if err != nil {
		logger.Println("connecting to server failed")
		logger.Println(err.Error())
		os.Exit(1)
	}
	logger.Println("connecting to server succeeded")

	// add callbacks
	irccon.AddCallback("001", func(e *irc.Event) {
		logger.Println("joining " + *ircChannel)
		irccon.Join(*ircChannel)
	})

	logger.Println("starting event loop")
	irccon.Loop()
}
