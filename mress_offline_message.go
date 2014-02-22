package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thoj/go-ircevent" // imported as "irc"
	"log"
	"strings"
)

// Inital setup of database. Handle things as needed to reduce
// false alarms.
func initOfflineMessageDatabase() error {
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
	return nil
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

// Retrieve and deliver previously stored message for user.
func deliverOfflineMessage(user string, con *irc.Connection) error {
	// sanity checks
	if len(user) == 0 {
		return fmt.Errorf("user of zero-length")
	}
	if 0 != strings.Count(user, " ") {
		return fmt.Errorf("user not allowed to contain whitespace")
	}
	if con == nil {
		return fmt.Errorf("connection pointer is nil")
	}

	// prepare db
	db, err := sql.Open("sqlite3", "./messages.db")
	if err != nil {
		return fmt.Errorf("failed to open database file: " + err.Error())
	}
	defer db.Close()

	// query db
	rows, err := db.Query("SELECT source, content FROM messages WHERE target = ?", user)
	if err != nil {
		return fmt.Errorf("query failed: " + err.Error())
	}
	defer rows.Close()

	// process retrieved information
	source := ""
	message := ""
	for rows.Next() {
		rows.Scan(&source, &message)
		con.Privmsg(user, "message from "+source+": "+message+"\n")
	}

	// delete this message from db if needed
	if (len(source) != 0) && (len(message) != 0) {
		_, err = db.Exec("DELETE FROM messages WHERE target = ? AND source = ?", user, source)
		if err != nil {
			return fmt.Errorf("executing DELETE failed: " + err.Error())
		}
	}

	return nil
}

// Implements the offline messenger command to deliver messages to other upon JOIN.
// To be in used as a callback for PRIVMSG.
// mress command: tell <nick>: <message>
// See also offlineMessengerDrone()
func offlineMessengerCommand(e *irc.Event, irc *irc.Connection, user string, logger *log.Logger) {
	// sanity checks
	if e == nil {
		return
	}
	if irc == nil {
		return
	}
	if len(user) == 0 {
		return
	}
	if logger == nil {
		return
	}
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

	// store the message
	target := strings.Fields(e.Message())[1]
	target = strings.Trim(target, ":")
	msgstart := strings.Index(e.Message(), ":") + 1
	err := saveOfflineMessage(e.Nick, target, e.Message()[msgstart:])
	if err != nil {
		logger.Println("offline message command failed")
		logger.Println(err.Error())
	}
	logger.Println("offline message saved")
}

// Deliver a message from a database. To be used as a callback for JOIN.
// This implements the delivery part of the offline messenger command.
// See also offlineMessengerCommand()
func offlineMessengerDrone(e *irc.Event, irc *irc.Connection, user, channel string, logger *log.Logger) {
	// sanity checks
	if e == nil {
		return
	}
	if irc == nil {
		return
	}
	if len(user) == 0 {
		return
	}
	if len(channel) == 0 {
		return
	}
	if logger == nil {
		return
	}

	// check for being a callback for an event intended
	// JOIN and 353 (names list)
	if !((e.Code == "JOIN") || (e.Code == "353")) {
		return
	}
	// ignore OTR -> potentially dead code?
	if 0 == strings.Index(e.Message(), "?OTR") {
		return
	}

	// TODO: handle self-join: if mress enters channel, deliver messages
	// 353 hf_testbot2 @ #ircscribble :hf_testbot2 tzugh @herr_flupke\r\n
	if e.Code == "353" {
		// e.Nick is empty for 353
		// strip "@" from op name
		nickline := strings.Replace(e.Message(), "@", "", -1)
		nicklist := strings.Fields(nickline)
		for i := 0; i < len(nicklist); i++ {
			err := deliverOfflineMessage(nicklist[i], irc)
			if err != nil {
				logger.Println("delivering stale messages had problems")
				logger.Println(err.Error())
			}
		}
		return
	}
	// handle others joining
	err := deliverOfflineMessage(e.Nick, irc)
	if err != nil {
		logger.Println("message delivery had problems")
		logger.Println(err.Error())
	}
}
