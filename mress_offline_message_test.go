package main

import (
	"github.com/thoj/go-ircevent"
	"os"
	"testing"
)

// test db initialization
func Test_initOfflineMessageDatabase_SL3_0(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", dbname: "mress-data", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}
	err = os.Remove(config.filename)
	if nil != err {
		t.Error(err.Error())
	}
}

func Test_initOfflineMessageDatabase_PG_0(t *testing.T) {
	config := MressDbConfig{backend: "postgres", filename: "testoffline.db", dbname: "mress-data", offlineMsgTable: "messages"}
	logger := createLogger("")
	dbuserchan := make(chan string)
	go getMressDbUser("", "test2.ini", dbuserchan, logger)
	dbpasswdchan := make(chan string)
	go getMressDbPassword("", "test2.ini", dbpasswdchan, logger)
	config.password = <-dbpasswdchan
	config.user = <-dbuserchan
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_initOfflineMessageDatabase_SL3_1(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err == nil {
		t.Error("empty filename did not yield error")
	}
}

func Test_initOfflineMessageDatabase_PG_1(t *testing.T) {
	config := MressDbConfig{backend: "postgres", dbname: "mress-data", user: "iuwefhf", password: "oidhfri", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Error("empty sqlite3 filename did yield error despite using Postgres backend: " + err.Error())
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

func Test_initOfflineMessageDatabase_SL3_4(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: ""}
	err := initOfflineMessageDatabase(config)
	if err == nil {
		t.Error("did not catch missing/empty offline message table name")
	}
}

func Test_initOfflineMessageDatabase_PG_4(t *testing.T) {
	config := MressDbConfig{backend: "postgres", dbname: "mress-data", offlineMsgTable: ""}
	err := initOfflineMessageDatabase(config)
	if err == nil {
		t.Error("did not catch missing/empty offline message table name")
	}
}

func Test_initOfflineMessageDatabase_5(t *testing.T) {
	config := MressDbConfig{backend: "postgres", user: "safrg", password: "supersecret", offlineMsgTable: ""}
	err := initOfflineMessageDatabase(config)
	if err == nil {
		t.Error("did not catch missing/empty postgres database name")
	} else {
		t.Log("testing for missing database name")
		t.Log(err.Error())
	}
}

func Test_initOfflineMessageDatabase_6(t *testing.T) {
	config := MressDbConfig{backend: "postgres", dbname: "mress-data", user: "", password: "supersecret", offlineMsgTable: ""}
	err := initOfflineMessageDatabase(config)
	if err == nil {
		t.Error("did not catch empty postgres database user name")
	}
}

func Test_initOfflineMessageDatabase_7(t *testing.T) {
	config := MressDbConfig{backend: "postgres", dbname: "mress-data", user: "", offlineMsgTable: ""}
	err := initOfflineMessageDatabase(config)
	if err == nil {
		t.Error("did not catch empty postgres database user password")
	}
}

func Test_initOfflineMessageDatabase_8(t *testing.T) {
	config := MressDbConfig{backend: "postgres", dbname: "fwfhweflhuew73", user: "", offlineMsgTable: ""}
	err := initOfflineMessageDatabase(config)
	if err == nil {
		t.Error("did not catch wrong postgres database name")
	}
}

// valid transaction
func Test_saveOfflineMessage_SL3_0(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}

	err = saveOfflineMessage(config, "testsource", "testtarget", "testmessage")
	if err != nil {
		t.Error(err.Error())
	}

	err = os.Remove(config.filename)
	if nil != err {
		t.Error(err.Error())
	}
}

func Test_saveOfflineMessage_PG_0(t *testing.T) {
	config := MressDbConfig{backend: "postgres", dbname: "mress-data", offlineMsgTable: "messages"}
	logger := createLogger("")
	dbuserchan := make(chan string)
	go getMressDbUser("", "test2.ini", dbuserchan, logger)
	dbpasswdchan := make(chan string)
	go getMressDbPassword("", "test2.ini", dbpasswdchan, logger)
	config.password = <-dbpasswdchan
	config.user = <-dbuserchan
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Log(config)
		t.Error(err.Error())
	}

	err = saveOfflineMessage(config, "testsource", "testtarget", "testmessage")
	if err != nil {
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
func Test_saveOfflineMessage_SL3_6(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", offlineMsgTable: "messages"}
	err := saveOfflineMessage(config, "test source", "testtarget", "testmessage")
	if err == nil {
		t.Error("empty database filename not detected")
	}
}

func Test_saveOfflineMessage_PG_6(t *testing.T) {
	config := MressDbConfig{backend: "postgres", offlineMsgTable: "messages"}
	logger := createLogger("")
	dbuserchan := make(chan string)
	go getMressDbUser("", "test2.ini", dbuserchan, logger)
	dbpasswdchan := make(chan string)
	go getMressDbPassword("", "test2.ini", dbpasswdchan, logger)
	dbnamechan := make(chan string)
	go getMressDbName("", "test2.ini", dbnamechan, logger)
	config.password = <-dbpasswdchan
	config.user = <-dbuserchan
	config.dbname = <-dbnamechan
	err := saveOfflineMessage(config, "testsource", "testtarget", "testmessage")
	if err != nil {
		t.Log(err.Error())
		t.Error("empty database filename shouldn't matter")
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
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := initOfflineMessageDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}
	con := &irc.Connection{}

	err = deliverOfflineMessage(config, "testuser", con)
	if err != nil {
		t.Log("valid call failed")
		t.Error(err.Error())
	}

	os.Remove(config.filename)
}

func Test_deliverOfflineMessage_1(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	con := &irc.Connection{}
	err := deliverOfflineMessage(config, "test user", con)
	if err == nil {
		t.Log("username with spaces shouldn't be accepted")
	}

	os.Remove(config.filename)
}

func Test_deliverOfflineMessage_2(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	con := &irc.Connection{}
	err := deliverOfflineMessage(config, "", con)
	if err == nil {
		t.Log("empty username shouldn't be accepted")
	}
	os.Remove(config.filename)
}

func Test_deliverOfflineMessage_3(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "", offlineMsgTable: "messages"}
	con := &irc.Connection{}
	err := deliverOfflineMessage(config, "testuser", con)
	if err == nil {
		t.Log("empty sqlite3 database filename shouldn't be accepted")
	}
	os.Remove(config.filename)
}

func Test_deliverOfflineMessage_4(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages"}
	err := deliverOfflineMessage(config, "testuser", nil)
	if err == nil {
		t.Log("nil connection pointer shouldn't be accepted")
	}
	os.Remove(config.filename)
}

func Test_deliverOfflineMessage_5(t *testing.T) {
	config := MressDbConfig{backend: "", filename: "testoffline.db", offlineMsgTable: "messages"}
	con := &irc.Connection{}
	err := deliverOfflineMessage(config, "testuser", con)
	if err == nil {
		t.Log("empty backend not detected")
	}
	os.Remove(config.filename)
}

func Test_deliverOfflineMessage_6(t *testing.T) {
	config := MressDbConfig{backend: "jksdgfrf", filename: "testoffline.db", offlineMsgTable: "messages"}
	con := &irc.Connection{}
	err := deliverOfflineMessage(config, "testuser", con)
	if err == nil {
		t.Log("invalid backend not detected")
	}
	os.Remove(config.filename)
}

func Test_deliverOfflineMessage_7(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "", offlineMsgTable: "messages"}
	con := &irc.Connection{}
	err := deliverOfflineMessage(config, "testuser", con)
	if err == nil {
		t.Log("empty sqlite3 filename not detected")
	}
	os.Remove(config.filename)
}

func Test_deliverOfflineMessage_8(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: ""}
	con := &irc.Connection{}
	err := deliverOfflineMessage(config, "testuser", con)
	if err == nil {
		t.Log("empty offline message table name not detected")
	}
	os.Remove(config.filename)
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
