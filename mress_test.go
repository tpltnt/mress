package main

import (
	//"github.com/thoj/go-ircevent"
	"log"
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

func Test_create_Logger_4(t *testing.T) {
	logger := createLogger(nil)
	if logger != nil {
		t.Error("creating with 'nil' destination didn't fail")
	}
}

func Test_readConfigInt_0(t *testing.T) {
	config := "config.ini"
	section := "IRC"
	key := "port"
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	if logger == nil {
		t.Log("creating test logger failed")
	}
	port, err := readConfigInt(config, section, key, logger)
	if err != nil {
		t.Fatal(err.Error())
	}
	if port != 6697 {
		t.Error("wrong integer read")
	}
}

func Test_readConfigInt_1(t *testing.T) {
	config := ""
	section := "IRC"
	key := "port"
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	if logger == nil {
		t.Log("creating test logger failed")
	}
	_, err := readConfigInt(config, section, key, logger)
	if err == nil {
		t.Error("failed to detect empty configuration file path")
	}
}

func Test_readConfigInt_2(t *testing.T) {
	config := "config.ini"
	section := ""
	key := "port"
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	if logger == nil {
		t.Log("creating test logger failed")
	}
	_, err := readConfigInt(config, section, key, logger)
	if err == nil {
		t.Fatal("failed to detect empty section string")
	}
}

func Test_readConfigInt_3(t *testing.T) {
	config := "config.ini"
	section := "IRC"
	key := ""
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	if logger == nil {
		t.Log("creating test logger failed")
	}
	_, err := readConfigInt(config, section, key, logger)
	if err == nil {
		t.Error("failed to detect empty key string")
	}
}

func Test_readConfigInt_4(t *testing.T) {
	config := "empty_config.ini"
	section := "IRC"
	key := "port"
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	if logger == nil {
		t.Log("creating test logger failed")
	}
	_, err := readConfigInt(config, section, key, logger)
	if err == nil {
		t.Error("failed to detect missing entries in config")
	}
}

func Test_readConfigString_0(t *testing.T) {
	config := "config.ini"
	section := "IRC"
	key := "server"
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	if logger == nil {
		t.Log("creating test logger failed")
	}
	server, err := readConfigString(config, section, key, logger)
	if err != nil {
		t.Fatal(err.Error())
	}
	if server != "chat.freenode.net" {
		t.Error("wrong server read")
	}
}

func Test_readConfigString_1(t *testing.T) {
	config := ""
	section := "IRC"
	key := "server"
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	if logger == nil {
		t.Log("creating test logger failed")
	}
	_, err := readConfigString(config, section, key, logger)
	if err == nil {
		t.Error("failed to detect empty configuration file path")
	}
}

func Test_readConfigString_2(t *testing.T) {
	config := "config.ini"
	section := ""
	key := "server"
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	if logger == nil {
		t.Log("creating test logger failed")
	}
	_, err := readConfigString(config, section, key, logger)
	if err == nil {
		t.Fatal("failed to detect empty section string")
	}
}

func Test_readConfigString_3(t *testing.T) {
	config := "config.ini"
	section := "IRC"
	key := ""
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	if logger == nil {
		t.Log("creating test logger failed")
	}
	_, err := readConfigString(config, section, key, logger)
	if err == nil {
		t.Error("failed to detect empty key string")
	}
}

func Test_readConfigString_4(t *testing.T) {
	config := "empty_config.ini"
	section := "IRC"
	key := "server"
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	if logger == nil {
		t.Log("creating test logger failed")
	}
	_, err := readConfigString(config, section, key, logger)
	if err == nil {
		t.Error("failed to detect missing entries in config")
	}
}

func Test_getLogger_0(t *testing.T) {
	dest := "/dev/null"
	conf := "config.ini"
	logchan := make(chan *log.Logger)
	go getLogger(dest, conf, logchan)
	//getLogger(destination, configfile string, logger chan *log.Logger)
	logger := <-logchan
	if logger == nil {
		t.Error("creating logger failed")
	}
}

func Test_getLogger_1(t *testing.T) {
	dest := ""
	conf := "config.ini"
	logchan := make(chan *log.Logger)
	go getLogger(dest, conf, logchan)
	//getLogger(destination, configfile string, logger chan *log.Logger)
	logger := <-logchan
	if logger != nil {
		t.Error("detecting empty destination string failed")
	}
}

func Test_getLogger_2(t *testing.T) {
	dest := "/dev/null"
	conf := ""
	logchan := make(chan *log.Logger)
	go getLogger(dest, conf, logchan)
	//getLogger(destination, configfile string, logger chan *log.Logger)
	logger := <-logchan
	if logger != nil {
		t.Error("detecting empty file path failed")
	}
}
