package main

import (
	//"github.com/thoj/go-ircevent"
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
