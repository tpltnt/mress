package main

import (
	"flag"
	"fmt"
	"github.com/thoj/go-ircevent" // imported as "irc"
	"log"
	"os"
	"io"
)


// Create a Logger which logs to the given destination
func createLogger(destination *string) *log.Logger {
	var logfile io.Writer
	var err error
	if len(*destination) > 0 {
		if "stdout" == *destination {
			logfile = os.Stdout
		} else {
			logfile, err = os.OpenFile(*destination, os.O_WRONLY, 0244)
		}
	} else {
		logfile, err = os.OpenFile("/dev/null", os.O_RDWR, 666)
	}
	if nil != err {
		fmt.Fprint(os.Stderr, "opening logging destination failed")
	}
	logger := log.New(logfile, "[mress] ", 0)
	return logger
}


func main() {
	fmt.Println("starting up ...")
	configfile := flag.String("config", "", "configuration file (lower priority if other flags are defined)")
	logdest := flag.String("log", "", "destination (filename, stdout) of the log")
	nick := flag.String("nick", "mress", "nickname")
	passwd := flag.String("passwd", "", "server/ident password")
	debug := flag.Bool("debug", false, "enable debugging (+flags)")
	flag.Parse()

	logger := createLogger(logdest)
	if nil == logger {
		fmt.Fprint(os.Stderr, "creating logger failed")
		os.Exit(1)
	}

	if 0 == len(*configfile) {
		fmt.Println("[info] no config file given, using defaults")
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
	if 0 < len(*passwd) {
		ircobj.Password = passwd
	}
}
