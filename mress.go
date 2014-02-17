package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thoj/go-ircevent" // imported as "irc"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

// Store a message for a target (user). If saving fails, this fact
// is going to be logged (but not the message content)
func saveOfflineMessage(source, target, message string) error {
	// sanity checks
	if len(source) == 0 {
		return fmt.Errorf("source of zero-length")
	}
	if 0 != strings.Count(source, " ") {
		return fmt.Errorf("source not allowed to contain whitespace")
	}
	if len(target) == 0 {
		return fmt.Errorf("target of zero-length")
	}
	if 0 != strings.Count(target, " ") {
		return fmt.Errorf("target not allowed to contain whitespace")
	}
	if len(message) == 0 {
		return fmt.Errorf("message of zero lenght")
	}

	// prepare db
	db, err := sql.Open("sqlite3", "./messages.db")
	if err != nil {
		return fmt.Errorf("failed to open database file: " + err.Error())
	}
	defer db.Close()
	sql := `CREATE TABLE IF NOT EXISTS messages (target TEXT, source TEXT, content TEXT);`
	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("failed to create database table: " + err.Error())
	}

	// prepare transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction failed: " + err.Error())
	}
	stmt, err := tx.Prepare("INSERT INTO messages (target, source, content) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// execute transaction
	_, err = stmt.Exec(target, source, message)
	if err != nil {
		return fmt.Errorf("executing INSERT failed: " + err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commiting to database failed: " + err.Error())
	}
	return nil
}

// Retrieve previously stored message for user.
func retrieveOfflineMessage(user string) error {
	// sanity checks
	if len(user) == 0 {
		return fmt.Errorf("user of zero-length")
	}
	if 0 != strings.Count(user, " ") {
		return fmt.Errorf("user not allowed to contain whitespace")
	}

	// prepare db
	db, err := sql.Open("sqlite3", "./messages.db")
	if err != nil {
		return fmt.Errorf("failed to open database file: " + err.Error())
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT source, content FROM messages WHERE target = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare query for stored message: " + err.Error())
	}
	defer stmt.Close()

	// query db
	rows, err := db.Query(user)
	if err != nil {
		return fmt.Errorf("query failed: " + err.Error())
	}
	defer rows.Close()

	// process retrieved information
	for rows.Next() {
		var source string
		var message string
		rows.Scan(&source, &message)
		//TODO: handle information
	}
	rows.Close()

	return nil
}

// Implements the offline messenger command to deliver messages to other upon JOIN.
// To be in used as a callback for PRIVMSG.
// mress command: tell <nick>: <message>
// See also offlineMessengerDrone()
func offlineMessengerCommand(e *irc.Event, irc *irc.Connection, user string, logger *log.Logger) {
	// ignore OTR
	if 0 == strings.Index(e.Message(), "?OTR") {
		return
	}
	// reject non-direct messages
	if user != e.Arguments[0] {
		return
	}
	// detect command -> reject non-command
	if 0 != strings.Index(e.Message(), "tell ") {
		return
	}
	if 5 > strings.Index(e.Message(), ":") {
		return
	}

	// extract message recipient
	target := strings.Fields(e.Message())[1]
	target = strings.Trim(target, ":")
	msgstart := strings.Index(e.Message(), ":")
	err := saveOfflineMessage(e.Nick, target, e.Message()[msgstart:])
	if err != nil {
		logger.Println("offline message command failed")
		logger.Println(err.Error())
	}
}

// Deliver a message from a database. To be used as a callback for JOIN.
// This implements the delivery part of the offline messenger command.
// See also offlineMessengerCommand()
func offlineMessengerDrone(e *irc.Event, irc *irc.Connection, user, channel string, logger *log.Logger) {
}

// The banana test
func bananaTest(e *irc.Event, irc *irc.Connection, user, channel string) {
	time.Sleep(1 * time.Second)
	// ignore OTR
	if 0 == strings.Index(e.Message(), "?OTR") {
		return
	}
	if user == e.Arguments[0] {
		irc.Privmsg(e.Nick, "I'm not actually a banana, i am parrot!\n")
		time.Sleep(1 * time.Second)
		irc.Privmsg(e.Nick, "\""+e.Message()+"\"")
		time.Sleep(2 * time.Second)
		irc.Privmsg(e.Nick, "see ?\n")
	}
	if channel == e.Arguments[0] {
		if 0 != strings.Index(e.Message(), "mress:") {
			return
		}
		irc.Privmsg(channel, "I'm a banana!\n")
	}
}

// Print the message associated with the event to stdout.
// Useful for debugging
func msgStdout(e *irc.Event) {
	fmt.Println(e.Message())
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

	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		offlineMessengerCommand(e, irccon, *nick, logger)
	})

	logger.Println("starting event loop")
	irccon.Loop()
}
