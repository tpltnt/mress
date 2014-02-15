package main

import (
	"flag"
	"fmt"
	"github.com/thoj/go-ircevent" // imported as "irc"
	"log"
	"os"
	"io"
)

func main() {
	fmt.Println("starting up ...")
	configfile := flag.String("config", "", "configuration file")
	logdest := flag.String("log", "", "destination (filename, stdout) of the log")
	debug := flag.Bool("debug", false, "enable debugging (+flags)")
	flag.Parse()

	var logfile io.Writer
	var err error
	if len(*logdest) > 0 {
		if "stdout" == *logdest {
			logfile = os.Stdout
		} else {
			logfile, err = os.OpenFile(*logdest, os.O_WRONLY, 0244)
		}
	} else {
		logfile, err = os.OpenFile("/dev/null", os.O_RDWR, 666)
	}
	if nil != err {
		fmt.Fprint(os.Stderr, "opening logging destination failed")
	}
	logger := log.New(logfile, "[mress]", 0)
	if nil == logger {
		fmt.Fprint(os.Stderr, "creating logger failed")
		os.Exit(1)
	}

	if nil == configfile {
		fmt.Println("[info] no config file given, using defaults")
	}

	if *debug {
		logger.Println("[info] debug mode enabled")
	}

	logger.Print("creating IRC connection ")
	ircobj := irc.IRC("<nick>", "<user>")
	if nil == ircobj {
		logger.Print("failed\n")
	} else {
		logger.Print("worked\n")
	}
}
