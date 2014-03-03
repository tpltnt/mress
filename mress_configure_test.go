package main

import (
	//"github.com/thoj/go-ircevent"
	"log"
	"os"
	"strconv"
	"testing"
)

// test config datastructure for validity
func Test_validateMressDbConfig_SL3_0(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "configtest.db", offlineMsgTable: "messages"}
	err := validateMressDbConfig(config)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_validateMressDbConfig_SL3_1(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "configtest.db"}
	err := validateMressDbConfig(config)
	if err == nil {
		t.Error("missing offline message table name should yield an error")
	}
}

func Test_validateMressDbConfig_PG_1(t *testing.T) {
	config := MressDbConfig{backend: "postgres", dbname: "mress-data", user: "db-user", password: "supersecret"}
	err := validateMressDbConfig(config)
	if err == nil {
		t.Error("missing offline message table name should yield an error")
	} else {
		t.Log(err.Error())
	}
}

func Test_validateMressDbConfig_SL3_2(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", offlineMsgTable: "messages"}
	err := validateMressDbConfig(config)
	if err == nil {
		t.Error("missing filename should yield an error")
	}
}

func Test_validateMressDbConfig_PG_3(t *testing.T) {
	config := MressDbConfig{backend: "postgres", user: "db-user", password: "supersecret", offlineMsgTable: "messages"}
	err := validateMressDbConfig(config)
	if err == nil {
		t.Error("missing database name should yield an error")
	} else {
		t.Log(err.Error())
	}
}

func Test_validateMressDbConfig_PG_4(t *testing.T) {
	config := MressDbConfig{backend: "postgres", dbname: "mress-data", password: "supersecret", offlineMsgTable: "messages"}
	err := validateMressDbConfig(config)
	if err == nil {
		t.Error("missing database user should yield an error")
	} else {
		t.Log(err.Error())
	}
}

func Test_validateMressDbConfig_PG_5(t *testing.T) {
	config := MressDbConfig{backend: "postgres", dbname: "mress-data", user: "db-user", offlineMsgTable: "messages"}
	err := validateMressDbConfig(config)
	if err == nil {
		t.Error("missing database password yield an error")
	} else {
		t.Log(err.Error())
	}
}

// test stdout destination
// TODO acually check stdout
func Test_create_Logger_1(t *testing.T) {
	dest := "stdout"
	logger := createLogger(dest)
	if logger == nil {
		t.Error("creating logger to '" + dest + "' returned 'nil'")
	}
}

// test stderr destination
// TODO actually check stderr
func Test_create_Logger_2(t *testing.T) {
	dest := "stderr"
	logger := createLogger("stderr")
	if logger == nil {
		t.Error("creating logger to '" + dest + "' returned 'nil'")
	}
}

// test logfile destination
func Test_create_Logger_3(t *testing.T) {
	dest := "test-logger.log"
	logger := createLogger(dest)
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
	logger := createLogger("")
	if logger == nil {
		t.Error("creating with empty destination did fail")
	}
}

func Test_readConfigInt_0(t *testing.T) {
	config := "test.ini"
	section := "IRC"
	key := "port"
	logger := createLogger("")
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
	logger := createLogger("")
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
	logger := createLogger("")
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
	logger := createLogger("")
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
	logger := createLogger("")
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
	logger := createLogger("")
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
	logger := createLogger("")
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
	logger := createLogger("")
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
	logger := createLogger("")
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
	logger := createLogger("")
	if logger == nil {
		t.Log("creating test logger failed")
	}
	_, err := readConfigString(config, section, key, logger)
	if err == nil {
		t.Error("failed to detect missing entries in config")
	}
}

func Test_getLogger_0(t *testing.T) {
	dest := ""
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
	dest := ""
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
	logger := createLogger("")
	go getChannel(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "#bar" {
		t.Error("read wrong channel")
	}
}

func Test_getChannel_1(t *testing.T) {
	testflag := ""
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getChannel(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "#foo" {
		t.Error("read wrong channel")
	}
}

func Test_getChannel_2(t *testing.T) {
	testflag := "#bar"
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getChannel(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "#bar" {
		t.Error("did not select flag over config value")
	}
}

func Test_getChannel_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getChannel(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "" {
		t.Error("did not handle empty/missing channel strings")
	}
}

func Test_getNick_0(t *testing.T) {
	testflag := "testbot"
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getNick(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != testflag {
		t.Error("read wrong nick")
	}
}

func Test_getNick_1(t *testing.T) {
	testflag := ""
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getNick(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "mress" {
		t.Error("read wrong nick (" + cstring + ") from config")
	}
}

func Test_getNick_2(t *testing.T) {
	testflag := "testbot"
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getNick(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "testbot" {
		t.Error("did not select flag over config value")
	}
}

func Test_getNick_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getNick(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "" {
		t.Error("did not handle empty/missing nick strings")
	}
}

func Test_getPassword_0(t *testing.T) {
	testflag := "424242"
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getPassword(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != testflag {
		t.Error("read wrong password")
	}
}

func Test_getPassword_1(t *testing.T) {
	testflag := ""
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getPassword(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "1234foobar" {
		t.Error("read wrong password (" + cstring + ") from config")
	}
}

func Test_getPassword_2(t *testing.T) {
	testflag := "424242"
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getNick(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "424242" {
		t.Error("did not select flag over config value")
	}
}

func Test_getPassword_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getNick(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "" {
		t.Error("did not handle empty/missing nick strings")
	}
}

func Test_getServer_0(t *testing.T) {
	testflag := "example.org"
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getServer(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != testflag {
		t.Error("read wrong server")
	}
}

func Test_getServer_1(t *testing.T) {
	testflag := ""
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getServer(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "chat.freenode.net" {
		t.Error("read wrong server (" + cstring + ") from config")
	}
}

func Test_getServer_2(t *testing.T) {
	testflag := "example.org"
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getServer(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "example.org" {
		t.Error("did not select flag over config value")
	}
}

func Test_getServer_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getServer(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "" {
		t.Error("did not handle empty/missing server strings")
	}
}

func Test_getPort_0(t *testing.T) {
	testflag := 23
	config := "test.ini"
	testchan := make(chan int)
	logger := createLogger("")
	go getPort(testflag, config, testchan, logger)
	cint := <-testchan
	if cint != testflag {
		t.Error("read wrong port")
	}
}

func Test_getPort_1(t *testing.T) {
	testflag := 0
	config := "test.ini"
	testchan := make(chan int)
	logger := createLogger("")
	go getPort(testflag, config, testchan, logger)
	cint := <-testchan
	if cint != 6697 {
		t.Error("read wrong port (" + strconv.Itoa(cint) + ") from config")
	}
}

func Test_getPort_2(t *testing.T) {
	testflag := 23
	config := "test.ini"
	testchan := make(chan int)
	logger := createLogger("")
	go getPort(testflag, config, testchan, logger)
	cint := <-testchan
	if cint != 23 {
		t.Error("did not select flag over config value")
	}
}

func Test_getPort_3(t *testing.T) {
	testflag := 0
	config := "empty_test.ini"
	testchan := make(chan int)
	logger := createLogger("")
	go getPort(testflag, config, testchan, logger)
	cint := <-testchan
	if cint != 0 {
		t.Error("did not handle missing port numbers")
	}
}

// test determining database filename
func Test_getOfflineDBfilename_0(t *testing.T) {
	testflag := "foobar.db"
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getOfflineDBfilename(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != testflag {
		t.Error("read wrong database filename")
	}
}

func Test_getOfflineDBfilename_1(t *testing.T) {
	testflag := ""
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getOfflineDBfilename(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "messages.db" {
		t.Error("read wrong filename (" + cstring + ") from config")
	}
}

func Test_getOfflineDBfilename_2(t *testing.T) {
	testflag := "foobar.db"
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getOfflineDBfilename(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "foobar.db" {
		t.Error("did not select flag over config value")
	}
}

func Test_getOfflineDBfilename_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getOfflineDBfilename(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "" {
		t.Error("did not handle empty/missing database filename")
	}
}

// test determining table name for offline messages
func Test_getOfflineTableName_0(t *testing.T) {
	testflag := "blah"
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getOfflineTableName(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != testflag {
		t.Error("read wrong offline message table name")
	}
}

func Test_getOfflineTableName_1(t *testing.T) {
	testflag := ""
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getOfflineTableName(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "offline-messages" {
		t.Error("read wrong filename (" + cstring + ") from config")
	}
}

func Test_getOfflineTableName_2(t *testing.T) {
	testflag := "blah"
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getOfflineDBfilename(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "blah" {
		t.Error("did not select flag over config value")
	}
}

func Test_getOfflineTableName_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getOfflineTableName(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "" {
		t.Error("did not handle empty/missing database table in config")
	}
}

// test determining database name
func Test_getMressDbName_0(t *testing.T) {
	testflag := "foobar"
	config := "test2.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbName(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != testflag {
		t.Error("read wrong database name")
	}
}

func Test_getMressDbName_1(t *testing.T) {
	testflag := ""
	config := "test2.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbName(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "mress-data" {
		t.Error("read wrong database name (" + cstring + ") from config")
	}
}

func Test_getMressDbName_2(t *testing.T) {
	testflag := "foobar"
	config := "test2.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbName(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "foobar" {
		t.Error("did not select flag over config value")
	}
}

func Test_getMressDbName_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbName(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "" {
		t.Error("did not handle empty/missing database name")
	}
}

// test determining database user name
func Test_getMressDbUser_0(t *testing.T) {
	testflag := "testuser"
	config := "test2.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbUser(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != testflag {
		t.Error("read wrong database user")
	}
}

func Test_getMressDbUser_1(t *testing.T) {
	testflag := ""
	config := "test2.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbUser(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "mress-bot" {
		t.Error("read wrong database user (" + cstring + ") from config")
	}
}

func Test_getMressDbUser_2(t *testing.T) {
	testflag := "testuser"
	config := "test2.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbUser(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "testuser" {
		t.Error("did not select flag over config value")
	}
}

func Test_getMressDbUser_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbUser(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "" {
		t.Error("did not handle empty/missing database user")
	}
}

// test determining database password
func Test_getMressDbPassword_0(t *testing.T) {
	testflag := "secretpassword"
	config := "test2.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbPassword(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != testflag {
		t.Error("read wrong database secretpassword")
	}
}

func Test_getMressDbPassword_1(t *testing.T) {
	testflag := ""
	config := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbPassword(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "1234" {
		t.Error("read wrong database password (" + cstring + ") from config")
	}
}

func Test_getMressDbPassword_2(t *testing.T) {
	testflag := "secretpassword"
	config := "test2.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbPassword(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "secretpassword" {
		t.Error("did not select flag over config value")
	}
}

func Test_getMressDbPassword_3(t *testing.T) {
	testflag := ""
	config := "empty_test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getMressDbPassword(testflag, config, testchan, logger)
	cstring := <-testchan
	if cstring != "" {
		t.Error("did not handle empty/missing password entry")
	}
}

// test the geoIP server
func Test_getGeoipServer_0(t *testing.T) {
	configfile := "test.ini"
	testchan := make(chan string)
	logger := createLogger("")
	go getGeoipServer(configfile, testchan, logger)
	cstring := <-testchan
	if cstring != "geoip.foo.bar" {
		t.Error("wrong geoip server read")
	}
}

func Test_getGeoipPort_0(t *testing.T) {
	configfile := "test.ini"
	testchan := make(chan int)
	logger := createLogger("")
	go getGeoipPort(configfile, testchan, logger)
	cint := <-testchan
	if cint != 2345 {
		t.Error("wrong geoip server port read")
	}
}
