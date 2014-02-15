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
	var logfile io.Writer
	var err error
	if len(*destination) > 0 {
		switch *destination {
		case "stdout":
			logfile = os.Stdout
		case "stderr":
			logfile = os.Stderr
		default:
			logfile, err = os.OpenFile(*destination, os.O_WRONLY, 0244)
		}
	} else {
		logfile, err = os.OpenFile("/dev/null", os.O_RDWR, 666)
	}
	if nil != err {
		fmt.Fprint(os.Stderr, "opening logging destination failed\n")
		fmt.Fprint(os.Stderr, err.Error()+"\n")
	}
	logger := log.New(logfile, "[mress] ", 0)
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

	if *debug {
		logger.Println("[info] debug mode enabled")
	}

	ircobj := irc.IRC(*nick, "mress")
	if nil == ircobj {
		logger.Println("creating IRC connection failed")
	} else {
		logger.Println("creating IRC connection worked")
	}
	// configure IRC connection
	if *useTLS {
		ircobj.UseTLS = true
		logger.Println("using TLS encrypted connection")
	} else {
		ircobj.UseTLS = false
		logger.Println("using cleartext connection")
	}
	ircobj.Password = *passwd
	if 0 < len(*passwd) {
		logger.Println("password is used")
	}
	// connect to server
	socketstring := *ircServer + ":" + strconv.Itoa(*ircPort)
	logger.Println("connecting to " + socketstring)
	logger.Println("joining " + *ircChannel)
}
