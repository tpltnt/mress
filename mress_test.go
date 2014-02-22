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
	config := "test.ini"
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
	config := "test.ini"
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
	config := "test.ini"
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
	config := "empty_test.ini"
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
	config := "test.ini"
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
	config := "test.ini"
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
	config := "test.ini"
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
	config := "empty_test.ini"
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
	conf := "test.ini"
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
	conf := "test.ini"
	logchan := make(chan *log.Logger)
	go getLogger(dest, conf, logchan)
	logger := <-logchan
	if logger == nil {
		t.Error("handling empty destination string failed")
	}
}

func Test_getLogger_2(t *testing.T) {
	dest := "/dev/null"
	conf := ""
	logchan := make(chan *log.Logger)
	go getLogger(dest, conf, logchan)
	logger := <-logchan
	if logger == nil {
		t.Error("handling empty file path failed")
	}
}

func Test_getChannel_0(t *testing.T) {
	testflag := "#bar"
	config := "test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getChannel(testflag, config, testchan, logger)
	cchannel := <-testchan
	if cchannel != "#bar" {
		t.Error("read wrong channel")
	}
}

func Test_getChannel_1(t *testing.T) {
	testflag := ""
	config := "test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getChannel(testflag, config, testchan, logger)
	cchannel := <-testchan
	if cchannel != "#foo" {
		t.Error("read wrong channel")
	}
}

func Test_getChannel_2(t *testing.T) {
	testflag := "#bar"
	config := "test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getChannel(testflag, config, testchan, logger)
	cchannel := <-testchan
	if cchannel != "#bar" {
		t.Error("did not select flag over config value")
	}
}

func Test_getChannel_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getChannel(testflag, config, testchan, logger)
	cchannel := <-testchan
	if cchannel != "" {
		t.Error("did not handle empty/missing channel strings")
	}
}

func Test_getNick_0(t *testing.T) {
	testflag := "testbot"
	config := "test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getNick(testflag, config, testchan, logger)
	cnick := <-testchan
	if cnick != testflag {
		t.Error("read wrong nick")
	}
}

func Test_getNick_1(t *testing.T) {
	testflag := ""
	config := "test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getNick(testflag, config, testchan, logger)
	cnick := <-testchan
	if cnick != "mress" {
		t.Error("read wrong nick (" + cnick + ") from config")
	}
}

func Test_getNick_2(t *testing.T) {
	testflag := "testbot"
	config := "test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getNick(testflag, config, testchan, logger)
	cnick := <-testchan
	if cnick != "testbot" {
		t.Error("did not select flag over config value")
	}
}

func Test_getNick_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getNick(testflag, config, testchan, logger)
	cnick := <-testchan
	if cnick != "" {
		t.Error("did not handle empty/missing nick strings")
	}
}

//getPassword(ipasswd, configfile string, channel chan string, logger *log.Logger)
func Test_getNick_0(t *testing.T) {
	testflag := "1234foobar"
	config := "test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getNick(testflag, config, testchan, logger)
	cnick := <-testchan
	if cnick != testflag {
		t.Error("read wrong nick")
	}
}

func Test_getNick_1(t *testing.T) {
	testflag := ""
	config := "test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getNick(testflag, config, testchan, logger)
	cnick := <-testchan
	if cnick != "mress" {
		t.Error("read wrong nick (" + cnick + ") from config")
	}
}

func Test_getNick_2(t *testing.T) {
	testflag := "testbot"
	config := "test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getNick(testflag, config, testchan, logger)
	cnick := <-testchan
	if cnick != "testbot" {
		t.Error("did not select flag over config value")
	}
}

func Test_getNick_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logdest := "/dev/null"
	logger := createLogger(&logdest)
	go getNick(testflag, config, testchan, logger)
	cnick := <-testchan
	if cnick != "" {
		t.Error("did not handle empty/missing nick strings")
	}
}

// getServer(iserver, configfile string, channel chan string, logger *log.Logger)
// getPort(iport int, configfile string, channel chan int, logger *log.Logger)
