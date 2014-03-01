package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thoj/go-ircevent" // imported as "irc"
	"log"
	"strings"
)

// Inital setup of the database. Handle things as needed
// to reduce false alarms.
func initOfflineMessageDatabase(config MressDbConfig) error {
	if len(config.backend) == 0 {
		return fmt.Errorf("empty backend string given")
	}
	if !((config.backend == "sqlite3") || (config.backend == "postgres")) {
		return fmt.Errorf("backend/database not supported")
	}
	if config.backend == "sqlite3" {
		if len(config.filename) == 0 {
			return fmt.Errorf("empty filename given")
		}
	}
	if config.backend == "postgres" {
		if len(config.dbname) == 0 {
			return fmt.Errorf("empty database name given")
		}
		if len(config.password) == 0 {
			return fmt.Errorf("empty database password given")
		}
	}
	if len(config.offlineMsgTable) == 0 {
		return fmt.Errorf("no offline message table name given")
	}
	var err error = nil
	//TODO: clean up ugly hack
	db, _ := sql.Open("", "")
	if config.backend == "sqlite3" {
		db, err = sql.Open("sqlite3", config.filename)
	}
	if config.backend == "postgres" {
		db, err = sql.Open("postgres", "host=localhost user=mress-bot password="+config.password+" dbname="+config.dbname+" sslmode=disable")
	}
	if err != nil {
		return fmt.Errorf("failed to open database: " + err.Error())
	}
	defer db.Close()
	sql := "CREATE TABLE IF NOT EXISTS " + config.offlineMsgTable + " (target TEXT, source TEXT, content TEXT);"
	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("failed to create database table: " + err.Error())
	}
	return nil
}

// Store a message for a target (user). If saving fails, this fact
// is going to be logged (but not the message content)
func saveOfflineMessage(dbconfig MressDbConfig, source, target, message string) error {
	// sanity checks
	if len(dbconfig.backend) == 0 {
		return fmt.Errorf("no backend given")
	}
	if dbconfig.backend != "sqlite3" {
		return fmt.Errorf("backend not supportend")
	}
	if len(dbconfig.filename) == 0 {
		return fmt.Errorf("empty database filename")
	}
	if len(dbconfig.offlineMsgTable) == 0 {
		return fmt.Errorf("no name for offline message table given")
	}
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
	var err error = nil
	// TODO fix ugly hack
	db, _ := sql.Open("", "")
	if dbconfig.backend == "sqlite3" {
		db, err = sql.Open("sqlite3", dbconfig.filename)
	}
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
func deliverOfflineMessage(dbconfig MressDbConfig, user string, con *irc.Connection) error {
	// sanity checks
	if dbconfig.backend != "sqlite3" {
		return fmt.Errorf("backend not supported")
	}
	if len(dbconfig.filename) == 0 {
		return fmt.Errorf("database filename is empty")
	}
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
	var err error = nil
	// TODO fix ugly hack
	db, _ := sql.Open("", "")
	if dbconfig.backend == "sqlite3" {
		db, err = sql.Open("sqlite3", dbconfig.filename)
	}
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
func offlineMessengerCommand(e *irc.Event, irc *irc.Connection, user string, dbconfig MressDbConfig, logger *log.Logger) {
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
	if dbconfig.backend != "sqlite3" {
		return
	}
	if len(dbconfig.filename) == 0 {
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
	err := saveOfflineMessage(dbconfig, e.Nick, target, e.Message()[msgstart:])
	if err != nil {
		logger.Println("offline message command failed")
		logger.Println(err.Error())
	}
	logger.Println("offline message saved")
}

// Deliver a message from a database. To be used as a callback for JOIN.
// This implements the delivery part of the offline messenger command.
// See also offlineMessengerCommand()
func offlineMessengerDrone(e *irc.Event, irc *irc.Connection, dbconfig MressDbConfig, user, channel string, logger *log.Logger) {
	// sanity checks
	if e == nil {
		return
	}
	if irc == nil {
		return
	}
	if dbconfig.backend != "sqlite3" {
		return
	}
	if len(dbconfig.filename) == 0 {
		return
	}
	if len(dbconfig.offlineMsgTable) == 0 {
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
			err := deliverOfflineMessage(dbconfig, nicklist[i], irc)
			if err != nil {
				logger.Println("delivering stale messages had problems")
				logger.Println(err.Error())
			}
		}
		return
	}
	// handle others joining
	err := deliverOfflineMessage(dbconfig, e.Nick, irc)
	if err != nil {
		logger.Println("message delivery had problems")
		logger.Println(err.Error())
	}
}
