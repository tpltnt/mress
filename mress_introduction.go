package main

/* Introduce bot to first-time users.
 * Track nicks as already seen to avoid muliple introductions etc. */
import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thoj/go-ircevent" // imported as "irc"
	"log"
	"strings"
	"time"
)

// Callback on JOIN to introduce bot functionality.
func runIntroduction(e *irc.Event, irc *irc.Connection, user, channel string) {
	// ignore OTR
	if 0 == strings.Index(e.Message(), "?OTR") {
		return
	}

	irc.Privmsg(e.Nick, "Hi "+e.Nick+"\n")
	time.Sleep(1 * time.Second)
	irc.Privmsg(e.Nick, "I am "+user+", the bot for this channel.\n")
	time.Sleep(1 * time.Second)
	irc.Privmsg(e.Nick, "Currently I can enable offline messaging for you.\n")
	time.Sleep(2 * time.Second)
	irc.Privmsg(e.Nick, "To leave a message for another user just type the following:\n")
	time.Sleep(2 * time.Second)
	irc.Privmsg(e.Nick, "/msg "+user+" tell <username>: <your message>\n")
	time.Sleep(3 * time.Second)
	irc.Privmsg(e.Nick, "I will then deliver it to the user you mentioned as soon as (s)he joins the channel.\n")
}

// Mark a nick as seen (and log to a database).
func markAsSeen(dbconfig MressDbConfig, user string, logger *log.Logger) error {
	err := validateMressDbConfig(dbconfig)
	if err != nil {
		return err
	}
	if len(user) == 0 {
		return fmt.Errorf("emtpy user name given")
	}
	if logger == nil {
		return fmt.Errorf("no logger given, only nil-pointer")
	}

	err = nil
	//TODO: clean up ugly hack
	db, _ := sql.Open("", "")
	if dbconfig.backend == "sqlite3" {
		db, err = sql.Open("sqlite3", dbconfig.filename)
	}
	if dbconfig.backend == "postgres" {
		db, err = sql.Open("postgres", "host=localhost user=mress-bot password="+dbconfig.password+" dbname="+dbconfig.dbname+" sslmode=disable")
	}
	if err != nil {
		return fmt.Errorf("failed to open database: " + err.Error())
	}
	defer db.Close()

	// prepare transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction failed: " + err.Error())
	}
	var stmt *sql.Stmt = nil
	if dbconfig.backend == "sqlite3" {
		stmt, err = tx.Prepare("INSERT INTO " + dbconfig.introductionTable + " (nickname) VALUES (?)")
	}
	if dbconfig.backend == "postgres" {
		stmt, err = tx.Prepare("INSERT INTO " + dbconfig.introductionTable + " (nickname) VALUES ($1)")
	}
	if err != nil {
		logger.Println("marking a nickname as having received the introduction failed")
		logger.Println(err)
	}
	defer stmt.Close()

	// execute transaction
	_, err = stmt.Exec(user)
	if err != nil {
		return fmt.Errorf("executing INSERT failed: " + err.Error())
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commiting to database failed: " + err.Error())
	}

	return nil
}

// Check if a nickname has been seen before. This also returns false in case
// of any errors.
func hasBeenSeen(dbconfig MressDbConfig, user string, logger *log.Logger) bool {
	if logger == nil {
		return false
	}
	err := validateMressDbConfig(dbconfig)
	if err != nil {
		logger.Println("hasBeenSeen()")
		logger.Println("invalid database configuration structure found")
		logger.Println(dbconfig)
		return false
	}
	if len(user) == 0 {
		logger.Println("hasBeenSeen()")
		logger.Println("given username is empty")
		return false
	}

	err = nil
	//TODO: clean up ugly hack
	db, _ := sql.Open("", "")
	if dbconfig.backend == "sqlite3" {
		db, err = sql.Open("sqlite3", dbconfig.filename)
	}
	if dbconfig.backend == "postgres" {
		db, err = sql.Open("postgres", "host=localhost user=mress-bot password="+dbconfig.password+" dbname="+dbconfig.dbname+" sslmode=disable")
	}
	if err != nil {
		logger.Println("hasBeenSeen()")
		logger.Println("failed to open database: " + err.Error())
		return false
	}
	defer db.Close()

	var founduser string
	if dbconfig.backend == "sqlite3" {
		err = db.QueryRow("SELECT username FROM "+dbconfig.introductionTable+" WHERE nickname=?", user).Scan(&founduser)
	}
	if dbconfig.backend == "postgres" {
		err = db.QueryRow("SELECT username FROM "+dbconfig.introductionTable+" WHERE nickname=$?", user).Scan(&founduser)
	}
	if err != nil {
		logger.Println("hasBeenSeen()")
		logger.Println("query the nickname failed")
		logger.Println(err)
	}
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		logger.Println("hasBeenSeen()")
		logger.Println("query the nickname failed")
		logger.Println(err)
		return false
	default:
		if len(founduser) == 0 {
			// should never happen
			logger.Println("hasBeenSeen()")
			logger.Println("found nickname empty")
			return false
		} else {
			return true
		}
	}
	return false
}

// initialize database
func initIntroductionTrackingDatabase(dbconfig MressDbConfig) error {
	err := validateMressDbConfig(dbconfig)
	if err != nil {
		return err
	}
	err = nil
	//TODO: clean up ugly hack
	db, _ := sql.Open("", "")
	if dbconfig.backend == "sqlite3" {
		db, err = sql.Open("sqlite3", dbconfig.filename)
	}
	if dbconfig.backend == "postgres" {
		db, err = sql.Open("postgres", "host=localhost user=mress-bot password="+dbconfig.password+" dbname="+dbconfig.dbname+" sslmode=disable")
	}
	if err != nil {
		return fmt.Errorf("failed to open database: " + err.Error())
	}
	defer db.Close()

	sql := "CREATE TABLE IF NOT EXISTS " + dbconfig.introductionTable + " (nickname TEXT);"
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
