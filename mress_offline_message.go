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
	// sanity checks
	err := validateMressDbConfig(config)
	if err != nil {
		return err
	}
	err = nil
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
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("database connection failed: " + err.Error())
	}
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
	err := validateMressDbConfig(dbconfig)
	if err != nil {
		return err
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
	err = nil
	// TODO fix ugly hack
	db, _ := sql.Open("", "")
	if dbconfig.backend == "sqlite3" {
		db, err = sql.Open("sqlite3", dbconfig.filename)
	}
	if dbconfig.backend == "postgres" {
		db, err = sql.Open("postgres", "host=localhost user="+dbconfig.user+" password="+dbconfig.password+" dbname="+dbconfig.dbname+" sslmode=disable")
	}
	if err != nil {
		return fmt.Errorf("failed to open database file: " + err.Error())
	}
	defer db.Close()
	// prepare transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction failed: " + err.Error())
	}
	var stmt *sql.Stmt = nil
	if dbconfig.backend == "sqlite3" {
		stmt, err = tx.Prepare("INSERT INTO messages (target, source, content) VALUES (?, ?, ?)")
	}
	if dbconfig.backend == "postgres" {
		stmt, err = tx.Prepare("INSERT INTO messages (target, source, content) VALUES ($1, $2, $3)")
	}
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
	err := validateMressDbConfig(dbconfig)
	if err != nil {
		return err
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
	err = nil
	// TODO fix ugly hack
	db, _ := sql.Open("", "")
	if dbconfig.backend == "sqlite3" {
		db, err = sql.Open("sqlite3", dbconfig.filename)
	}
	if dbconfig.backend == "postgres" {
		db, err = sql.Open("postgres", "host=localhost user="+dbconfig.user+" password="+dbconfig.password+" dbname="+dbconfig.dbname+" sslmode=disable")
	}
	if err != nil {
		return fmt.Errorf("failed to open database file: " + err.Error())
	}
	defer db.Close()

	// query db
	var rows *sql.Rows = nil
	if dbconfig.backend == "sqlite3" {
		rows, err = db.Query("SELECT source, content FROM messages WHERE target = ?", user)
	}
	if dbconfig.backend == "postgres" {
		rows, err = db.Query("SELECT source, content FROM messages WHERE target = $1", user)
	}
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
		if dbconfig.backend == "sqlite3" {
			_, err = db.Exec("DELETE FROM messages WHERE target = ? AND source = ?", user, source)
		}
		if dbconfig.backend == "postgres" {
			_, err = db.Exec("DELETE FROM messages WHERE target = $1 AND source = $2", user, source)
		}
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
	if !((dbconfig.backend == "sqlite3") || (dbconfig.backend == "postgres")) {
		return
	}
	if dbconfig.backend == "sqlite3" {
		if len(dbconfig.filename) == 0 {
			return
		}
	}
	if dbconfig.backend == "postgres" {
		if len(dbconfig.dbname) == 0 {
			return
		}
		if len(dbconfig.user) == 0 {
			return
		}
		if len(dbconfig.password) == 0 {
			return
		}
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
		irc.Privmsg(e.Nick, "Sorry "+e.Nick+", something went wrong and I couldn't store your message. :(\n")
		return
	}
	logger.Println("offline message saved")
	irc.Privmsg(e.Nick, "Yes "+e.Nick+", I will deliver your message as soon as possible.\n")
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
	if !((dbconfig.backend == "sqlite3") || (dbconfig.backend == "postgres")) {
		return
	}
	if dbconfig.backend == "sqlite3" {
		if len(dbconfig.filename) == 0 {
			return
		}
	}
	if dbconfig.backend == "postgres" {
		if len(dbconfig.dbname) == 0 {
			return
		}
		if len(dbconfig.user) == 0 {
			return
		}
		if len(dbconfig.password) == 0 {
			return
		}
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
