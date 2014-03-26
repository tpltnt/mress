package main

/* Introduce bot to first-time users.
 * Track nicks as already seen to avoid muliple introductions etc. */
import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thoj/go-ircevent" // imported as "irc"
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
func markAsSeen(dbconfig MressDbConfig, user string) error {
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

	return fmt.Errorf("not implemented yet")
	return nil
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
