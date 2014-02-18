package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thoj/go-ircevent"
	"os"
	"testing"
)

// test stdout destination
// TODO acually check stdout
func Test_create_Logger_1(t *testing.T) {
	dest := "stdout"
	logger := createLogger(&dest)
	if logger == nil {
		t.Error("creating logger to '" + dest + "' returned 'nil'")
	}
}

// test stderr destination
// TODO actually check stderr
func Test_create_Logger_2(t *testing.T) {
	dest := "stderr"
	logger := createLogger(&dest)
	if logger == nil {
		t.Error("creating logger to '" + dest + "' returned 'nil'")
	}
}

// test logfile destination
func Test_create_Logger_3(t *testing.T) {
	dest := "test-logger.log"
	logger := createLogger(&dest)
	if logger == nil {
		t.Error("creating logger to '" + dest + "' returned 'nil'")
	}
	logger.Println("basic logger test")

	logfile, err := os.Open(dest)
	if nil != err {
		t.Error(err.Error())
	}
	filecontent := make([]byte, 100)
	count, err := logfile.Read(filecontent)
	if nil != err {
		t.Error("reading logfile failed: " + err.Error())
	}
	if 46 != count {
		t.Error("read wrong number of bytes")
	}

	logfile.Close()
	err = os.Remove(dest)
	if nil != err {
		t.Error(err.Error())
	}
}

// valid transaction
func Test_saveOfflineMessage_0(t *testing.T) {
	err := saveOfflineMessage("testsource", "testtarget", "testmessage")
	if err != nil {
		t.Error(err.Error())
	}
	err = os.Remove("./messages.db")
	if nil != err {
		t.Error(err.Error())
	}
}

// empty target
func Test_saveOfflineMessage_1(t *testing.T) {
	err := saveOfflineMessage("testsource", "", "testmessage")
	if err == nil {
		t.Error("empty target not detected")
	}
}

// target with space
func Test_saveOfflineMessage_2(t *testing.T) {
	err := saveOfflineMessage("testsource", "test target", "testmessage")
	if err == nil {
		t.Error("target with space not detected")
	}
}

// emtpy message
func Test_saveOfflineMessage_3(t *testing.T) {
	err := saveOfflineMessage("testsource", "testtarget", "")
	if err == nil {
		t.Error("empty message not detected")
	}
}

// empty source
func Test_saveOfflineMessage_4(t *testing.T) {
	err := saveOfflineMessage("", "testtarget", "testmessage")
	if err == nil {
		t.Error("empty source not detected")
	}
}

// source with space
func Test_saveOfflineMessage_5(t *testing.T) {
	err := saveOfflineMessage("test source", "testtarget", "testmessage")
	if err == nil {
		t.Error("source with space not detected")
	}
}

func Test_deliverOfflineMessage_0(t *testing.T) {
	// prepare db
	db, err := sql.Open("sqlite3", "./messages.db")
	if err != nil {
		t.Error("failed to open database file: " + err.Error())
	}
	defer db.Close()
	sql := `CREATE TABLE IF NOT EXISTS messages (target TEXT, source TEXT, content TEXT);`
	_, err = db.Exec(sql)
	if err != nil {
		t.Error("failed to create database table: " + err.Error())
	}

	con := &irc.Connection{}
	err = deliverOfflineMessage("testuser", con)
	if err != nil {
		t.Log("valid call failed")
		t.Error(err.Error())
	}

	os.Remove("./messages.db")
}

func Test_deliverOfflineMessage_1(t *testing.T) {
	con := &irc.Connection{}
	err := deliverOfflineMessage("test user", con)
	if err == nil {
		t.Log("username with spaces shouldn't be accepted")
	}
}

func Test_deliverOfflineMessage_2(t *testing.T) {
	con := &irc.Connection{}
	err := deliverOfflineMessage("", con)
	if err == nil {
		t.Log("empty username shouldn't be accepted")
	}
}

func Test_deliverOfflineMessage_3(t *testing.T) {
	err := deliverOfflineMessage("testuser", nil)
	if err == nil {
		t.Log("nil connection pointer shouldn't be accepted")
	}
}

// callbacks shouldn't explode
func Test_offlineMessengerCommand_0(t *testing.T) {
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	con := &irc.Connection{}
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	offlineMessengerCommand(event, con, "testuser", logger)
}

func Test_offlineMessengerCommand_1(t *testing.T) {
	con := &irc.Connection{}
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	offlineMessengerCommand(nil, con, "testuser", logger)
}

func Test_offlineMessengerCommand_2(t *testing.T) {
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	offlineMessengerCommand(event, nil, "testuser", logger)
}

func Test_offlineMessengerCommand_3(t *testing.T) {
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	con := &irc.Connection{}
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	offlineMessengerCommand(event, con, "test user", logger)
}

func Test_offlineMessengerCommand_4(t *testing.T) {
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	con := &irc.Connection{}
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	offlineMessengerCommand(event, con, "", logger)
}

func Test_offlineMessengerCommand_5(t *testing.T) {
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	con := &irc.Connection{}
	offlineMessengerCommand(event, con, "testuser", nil)
}
