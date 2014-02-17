package main

import (
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

// deliverOfflineMessage(user string, con *irc.Connection) error
// offlineMessengerCommand(e *irc.Event, irc *irc.Connection, user string, logger *log.Logger)
