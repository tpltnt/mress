package main

import (
	"fmt"
	"github.com/jurka/goini"
	"github.com/thoj/go-ircevent" // imported as "irc"
	"io"
	"log"
	"os"
)

// Create a Logger which logs to the given destination
// Valid destinations are files (+path), stdout and stderr
func createLogger(destination string) *log.Logger {
	var logdest io.Writer = nil
	var logfile *os.File = nil
	var logger *log.Logger = nil
	var err error
	if len(destination) > 0 {
		switch destination {
		case "stdout":
			logdest = os.Stdout
		case "stderr":
			logdest = os.Stderr
		default:
			// assuming the logfile already exists
			logfile, err = os.OpenFile(destination, os.O_WRONLY|os.O_APPEND, 0644)
			if nil != err {
				// it didn't, so create a new one
				logfile, err = os.Create(destination)
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
	logger <- createLogger(dest)
	return
}

// Get IRC channel and choose commandline value over config file.
// Return IRC channel through channel (to facilitate concurrent setups).
// A returning empty channel indicates errors.
func getChannel(flag, configfile string, channel chan string, logger *log.Logger) {
	if logger == nil {
		channel <- ""
		return
	}
	if len(flag) == 0 {
		irc, err := readConfigString(configfile, "IRC", "channel", logger)
		if err != nil {
			logger.Println(err.Error())
			channel <- ""
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
// A returning empty nick indicates errors.
func getNick(inick, configfile string, channel chan string, logger *log.Logger) {
	if logger == nil {
		return
	}
	cnick, _ := readConfigString(configfile, "IRC", "nickname", logger)
	// if emtpy flag -> choose config
	if len(inick) == 0 {
		channel <- cnick
		return
	}
	//choose config over default value
	if "mress" == inick {
		//default and config value -> config
		channel <- cnick
	} else {
		//non-default flag -> flag (over config)
		channel <- inick
	}
	return
}

// Get IRC server/hostname and choose commandline value over config file.
// Return IRC server through channel (to facilitate concurrent setups).
func getPassword(ipasswd, configfile string, channel chan string, logger *log.Logger) {
	if logger == nil {
		return
	}
	cpasswd, _ := readConfigString(configfile, "IRC", "password", logger)
	//choose config over default value
	if len(ipasswd) == 0 {
		//default and config value -> config
		channel <- cpasswd
	} else {
		//non-default flag -> flag (over config)
		channel <- ipasswd
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
// A port number of 0 indicates errors.
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
	//choose config over "empty" value
	if iport == 0 {
		//default and config value -> config
		channel <- cport
	} else {
		//non-default flag -> flag (over config)
		channel <- iport
	}
}

// read name of database file for offline messages
func getOfflineDBfilename(dbfile, configfile string, channel chan string, logger *log.Logger) {
	cdb, err := readConfigString(configfile, "offline messaging", "dbfile", logger)
	if err != nil {
		logger.Println(err.Error())
		channel <- ""
		return
	}
	//choose config over "empty" value
	if len(dbfile) == 0 {
		channel <- cdb
	} else {
		channel <- dbfile
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
	if err != nil {
		return "", fmt.Errorf("failed to get the " + key + " value")
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
	if err != nil {
		return 0, fmt.Errorf("failed to get the " + key + " value")
	}
	return value, nil
}
