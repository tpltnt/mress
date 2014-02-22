package main

import (
	"flag"
	"fmt"
	"github.com/jurka/goini"
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

// Print the message associated with the event to stdout.
// Useful for debugging
func msgStdout(e *irc.Event) {
	fmt.Println(e.Message())
}

// Build logger and choose commandline value over config file.
// Return created logger through channel (to facilitate concurrent setups).
func getLogger(destination, configfile string, logger chan *log.Logger) {
	dest := ""
	if len(destination) == 0 {
		// read config
		dest, _ = readConfigString(configfile, "maintainance", "log-destination", nil)
	} else {
		dest = destination
	}
	logger <- createLogger(&dest)
	return
}

// Get IRC channel and choose commandline value over config file.
// Return IRC channel through channel (to facilitate concurrent setups).
func getChannel(flag, configfile string, channel chan string, logger *log.Logger) {
	if logger == nil {
		channel <- ""
		return
	}
	if len(flag) == 0 {
		irc, err := readConfigString(configfile, "IRC", "channel", logger)
		if err != nil {
			logger.Println(err.Error())
			return
		}
		channel <- irc
	} else {
		channel <- flag
	}
	return
}

// Get IRC nickname and choose commandline value over config file.
// Return IRC nickname through channel (to facilitate concurrent setups).
func getNick(inick, configfile string, channel chan string, logger *log.Logger) {
	if logger == nil {
		return
	}
	cnick, _ := readConfigString(configfile, "IRC", "nickname", logger)
	//choose config over default value
	if "mress" == inick {
		//default and config value -> config
		channel <- cnick
	} else {
		//non-default flag -> flag (over config)
		channel <- inick
	}
}

// Get IRC server/hostname and choose commandline value over config file.
// Return IRC server through channel (to facilitate concurrent setups).
func getServer(iserver, configfile string, channel chan string, logger *log.Logger) {
	if logger == nil {
		return
	}
	cserver, _ := readConfigString(configfile, "IRC", "server", logger)
	//choose config over default value
	if len(iserver) == 0 {
		//default and config value -> config
		channel <- cserver
	} else {
		//non-default flag -> flag (over config)
		channel <- iserver
	}
}

// Get port to connect to and choose commandline value over config file.
// Return IRC server through channel (to facilitate concurrent setups).
func getPort(iport int, configfile string, channel chan int, logger *log.Logger) {
	if logger == nil {
		return
	}
	cport, err := readConfigInt(configfile, "IRC", "port", logger)
	if err != nil {
		logger.Println(err.Error())
		channel <- 0
		return
	}
	//choose config over default value
	if iport == 6697 {
		//default and config value -> config
		channel <- cport
	} else {
		//non-default flag -> flag (over config)
		channel <- iport
	}
}

// Read string from config file
func readConfigString(filename, section, key string, logger *log.Logger) (string, error) {
	if logger == nil {
		return "", fmt.Errorf("logger nil pointer\n")
	}
	conf, err := goini.LoadConfig(filename)
	if err != nil {
		return "", fmt.Errorf("failed to load configuration\n")
	}
	if len(section) == 0 {
		return "", fmt.Errorf("empty section string\n")
	}
	sec := conf.GetSection(section)
	if sec == nil {
		return "", fmt.Errorf("failed to load " + section + " section\n")
	}
	value, err := sec.GetString(key)
	if err == nil {
		return "", fmt.Errorf("failed to get the " + key + " value")
	}
	return value, nil
}

// Read bool from config file
func readConfigBool(filename, section, key string, logger *log.Logger) (bool, error) {
	if logger == nil {
		return false, fmt.Errorf("logger nil pointer\n")
	}
	conf, err := goini.LoadConfig(filename)
	if err != nil {
		return false, fmt.Errorf("failed to load configuration\n")
	}
	if len(section) == 0 {
		return false, fmt.Errorf("empty section string\n")
	}
	sec := conf.GetSection(section)
	if sec == nil {
		return false, fmt.Errorf("failed to load " + section + " section\n")
	}
	value, err := sec.GetBool(key)
	if err == nil {
		return false, fmt.Errorf("failed to get the " + key + " value")
	}
	return value, nil
}

// Read integer from config file
func readConfigInt(filename, section, key string, logger *log.Logger) (int, error) {
	if logger == nil {
		return 0, fmt.Errorf("logger nil pointer\n")
	}
	conf, err := goini.LoadConfig(filename)
	if err != nil {
		return 0, fmt.Errorf("failed to load configuration\n")
	}
	if len(section) == 0 {
		return 0, fmt.Errorf("empty section string\n")
	}
	sec := conf.GetSection(section)
	if sec == nil {
		return 0, fmt.Errorf("failed to load " + section + " section\n")
	}
	value, err := sec.GetInt(key)
	if err == nil {
		return 0, fmt.Errorf("failed to get the " + key + " value")
	}
	return value, nil
}

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
	flag.Parse()

	logchan := make(chan *log.Logger)
	go getLogger(*logdest, *configfile, logchan)
	logger := <-logchan
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

	// determine config values concurrently with go-routines
	// roughly according to need, choose non-default flags over
	// config file values
	nickchan := make(chan string)
	go getNick(*ircNick, *configfile, nickchan, logger)
	servchan := make(chan string)
	go getServer(*ircServer, *configfile, servchan, logger)
	portchan := make(chan int)
	go getPort(*ircPort, *configfile, portchan, logger)
	// tls
	// debug
	chanchan := make(chan string)
	go getChannel(*ircChannel, *configfile, chanchan, logger)
	// nick
	// password

	// collect all config values needed

	// create IRC connection
	nick := <-nickchan
	irccon := irc.IRC(nick, "mress")
	if nil == irccon {
		logger.Println("creating IRC connection failed")
	} else {
		logger.Println("creating IRC connection worked")
	}
	irccon.Password = *ircPasswd
	if 0 < len(*ircPasswd) {
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
		os.Exit(1)
	}
	logger.Println("connecting to server succeeded")

	// collect last config value needed
	channel := <-chanchan
	// add callbacks
	irccon.AddCallback("001", func(e *irc.Event) {
		logger.Println("joining " + channel)
		irccon.Join(channel)
	})

	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		offlineMessengerCommand(e, irccon, nick, logger)
	})
	irccon.AddCallback("JOIN", func(e *irc.Event) {
		offlineMessengerDrone(e, irccon, nick, channel, logger)
	})
	irccon.AddCallback("353", func(e *irc.Event) {
		offlineMessengerDrone(e, irccon, nick, channel, logger)
	})

	logger.Println("starting event loop")
	irccon.Loop()
}
