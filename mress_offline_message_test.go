package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thoj/go-ircevent"
	"os"
	"testing"
)

// test db initialization
func Test_initOfflineMessageDatabase_0(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}
	err = os.Remove(config.filename)
	if nil != err {
		t.Error(err.Error())
	}
}

func Test_initOfflineMessageDatabase_1(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err == nil {
		t.Error("empty filename did not yield error")
	}
}

func Test_initOfflineMessageDatabase_2(t *testing.T) {
	config := MressDbConfig{backend: "", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err == nil {
		t.Error("did not catch missing/empty backend")
	}
}

func Test_initOfflineMessageDatabase_3(t *testing.T) {
	config := MressDbConfig{backend: "kjsdgfjds", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err == nil {
		t.Error("did not catch invalid backend")
	}
}

func Test_initOfflineMessageDatabase_4(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: ""}
	err := initOfflineMessageDatabase(config)
	if err == nil {
		t.Error("did not catch missing/empty offline message table name")
	}
}

// valid transaction
func Test_saveOfflineMessage_0(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := saveOfflineMessage(config, "testsource", "testtarget", "testmessage")
	if err != nil {
		t.Error(err.Error())
	}
	err = os.Remove(config.filename)
	if nil != err {
		t.Error(err.Error())
	}
}

// empty target
func Test_saveOfflineMessage_1(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}
	err = saveOfflineMessage(config, "testsource", "", "testmessage")
	if err == nil {
		t.Error("empty target not detected")
	}
	err = os.Remove(config.filename)
	if nil != err {
		t.Error(err.Error())
	}
}

// target with space
func Test_saveOfflineMessage_2(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}
	err = saveOfflineMessage(config, "testsource", "test target", "testmessage")
	if err == nil {
		t.Error("target with space not detected")
	}
	err = os.Remove(config.filename)
	if nil != err {
		t.Error(err.Error())
	}
}

// emtpy message
func Test_saveOfflineMessage_3(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}
	err = saveOfflineMessage(config, "testsource", "testtarget", "")
	if err == nil {
		t.Error("empty message not detected")
	}
	err = os.Remove(config.filename)
	if nil != err {
		t.Error(err.Error())
	}
}

// empty source
func Test_saveOfflineMessage_4(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}
	err = saveOfflineMessage(config, "", "testtarget", "testmessage")
	if err == nil {
		t.Error("empty source not detected")
	}
	err = os.Remove(config.filename)
	if nil != err {
		t.Error(err.Error())
	}
}

// source with space
func Test_saveOfflineMessage_5(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}
	err = saveOfflineMessage(config, "test source", "testtarget", "testmessage")
	if err == nil {
		t.Error("source with space not detected")
	}
	err = os.Remove(config.filename)
	if nil != err {
		t.Error(err.Error())
	}
}

// empty db filename
func Test_saveOfflineMessage_6(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "", offlineMsgTable: "messages"}
	err := saveOfflineMessage(config, "test source", "testtarget", "testmessage")
	if err == nil {
		t.Error("empty database filename not detected")
	}
}

// empty offline msg table
func Test_saveOfflineMessage_7(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: ""}
	err := saveOfflineMessage(config, "test source", "testtarget", "testmessage")
	if err == nil {
		t.Error("empty offline message table name not detected")
	}
}

func Test_deliverOfflineMessage_0(t *testing.T) {
	// prepare db
	dbfile := "testmsg.db"
	db, err := sql.Open("sqlite3", dbfile)
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
	err = deliverOfflineMessage(dbfile, "testuser", con)
	if err != nil {
		t.Log("valid call failed")
		t.Error(err.Error())
	}

	os.Remove(dbfile)
}

func Test_deliverOfflineMessage_1(t *testing.T) {
	dbfile := "testmsg.db"
	con := &irc.Connection{}
	err := deliverOfflineMessage(dbfile, "test user", con)
	if err == nil {
		t.Log("username with spaces shouldn't be accepted")
	}

	os.Remove(dbfile)
}

func Test_deliverOfflineMessage_2(t *testing.T) {
	dbfile := "testmsg.db"
	con := &irc.Connection{}
	err := deliverOfflineMessage(dbfile, "", con)
	if err == nil {
		t.Log("empty username shouldn't be accepted")
	}
	os.Remove(dbfile)
}

func Test_deliverOfflineMessage_3(t *testing.T) {
	dbfile := ""
	con := &irc.Connection{}
	err := deliverOfflineMessage(dbfile, "testuser", con)
	if err == nil {
		t.Log("nil connection pointer shouldn't be accepted")
	}
	os.Remove(dbfile)
}

func Test_deliverOfflineMessage_4(t *testing.T) {
	dbfile := "testmsg.db"
	err := deliverOfflineMessage(dbfile, "testuser", nil)
	if err == nil {
		t.Log("nil connection pointer shouldn't be accepted")
	}
	os.Remove(dbfile)
}

// callbacks shouldn't explode
func Test_offlineMessengerCommand_0(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	con := &irc.Connection{}
	logger := createLogger("")
	offlineMessengerCommand(event, con, "testuser", config, logger)
	os.Remove(config.filename)
}

func Test_offlineMessengerCommand_1(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	con := &irc.Connection{}
	logger := createLogger("")
	offlineMessengerCommand(nil, con, "testuser", config, logger)
	os.Remove(config.filename)
}

func Test_offlineMessengerCommand_2(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	logger := createLogger("")
	offlineMessengerCommand(event, nil, "testuser", config, logger)
	os.Remove(config.filename)
}

func Test_offlineMessengerCommand_3(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	con := &irc.Connection{}
	logger := createLogger("")
	offlineMessengerCommand(event, con, "test user", config, logger)
	os.Remove(config.filename)
}

func Test_offlineMessengerCommand_4(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	con := &irc.Connection{}
	logger := createLogger("")
	offlineMessengerCommand(event, con, "", config, logger)
	os.Remove(config.filename)
}

func Test_offlineMessengerCommand_5(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	con := &irc.Connection{}
	offlineMessengerCommand(event, con, "testuser", config, nil)
	os.Remove(config.filename)
}

func Test_offlineMessengerCommand_6(t *testing.T) {
	config := MressDbConfig{backend: "", filename: "testoffline.db", offlineMsgTable: "messages"}
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	con := &irc.Connection{}
	offlineMessengerCommand(event, con, "testuser", config, nil)
}

func Test_offlineMessengerCommand_7(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "", offlineMsgTable: "messages"}
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	con := &irc.Connection{}
	offlineMessengerCommand(event, con, "testuser", config, nil)
}

func Test_offlineMessengerCommand_8(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.sb", offlineMsgTable: ""}
	args := []string{"bla bla foo bar baz"}
	event := &irc.Event{Arguments: args}
	con := &irc.Connection{}
	offlineMessengerCommand(event, con, "testuser", config, nil)
}
