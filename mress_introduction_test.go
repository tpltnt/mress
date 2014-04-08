package main

import (
	//	"log"
	"os"
	//	"strconv"
	"testing"
)

func Test_markAsSeen_0(t *testing.T) {
	dbconfig := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages", introductionTable: "intro"}
	user := "testuser"
	logger := createLogger("")
	err := initIntroductionTrackingDatabase(dbconfig)
	if err != nil {
		t.Error(err.Error())
	}
	err = markAsSeen(dbconfig, user, logger)
	if err != nil {
		t.Error(err.Error())
	}

	err = os.Remove(dbconfig.filename)
	if err != nil {
		t.Error(err.Error())
	}
}

// test db initialization
func Test_initIntroductionTrackingDatabase_SL3_0(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages", introductionTable: "intro"}
	err := initIntroductionTrackingDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}
	err = os.Remove(config.filename)
	if nil != err {
		t.Error(err.Error())
	}
}

func Test_initIntroductionTrackingDatabase_PG_0(t *testing.T) {
	config := MressDbConfig{backend: "postgres", dbname: "mress-data", offlineMsgTable: "messages", introductionTable: "intro"}
	logger := createLogger("")
	dbuserchan := make(chan string)
	go getMressDbUser("", "test2.ini", dbuserchan, logger)
	dbpasswdchan := make(chan string)
	go getMressDbPassword("", "test2.ini", dbpasswdchan, logger)
	config.password = <-dbpasswdchan
	config.user = <-dbuserchan
	err := initIntroductionTrackingDatabase(config)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_initIntroductionTrackingDatabase_SL3_1(t *testing.T) {
	config := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "message", introductionTable: ""}
	err := initIntroductionTrackingDatabase(config)
	if err == nil {
		t.Error("did not catch missing/empty introduction table name")
	}
}

func Test_initIntroductionTrackingDatabase_PG_1(t *testing.T) {
	config := MressDbConfig{backend: "postgres", dbname: "mress-data", offlineMsgTable: "message", introductionTable: ""}
	err := initIntroductionTrackingDatabase(config)
	if err == nil {
		t.Error("did not catch missing/empty introduction table name")
	}
}

func Test_hasBeenSeen_0(t *testing.T) {
	dbconfig := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages", introductionTable: "intro"}
	user := "testuser"
	logger := createLogger("stdout")
	err := initIntroductionTrackingDatabase(dbconfig)
	if err != nil {
		t.Error(err.Error())
	}

	err = markAsSeen(dbconfig, user, logger)
	if err != nil {
		t.Error(err.Error())
	}

	seen := hasBeenSeen(dbconfig, user, logger)
	if seen != true {
		t.Error("explicitly marked user not recognized as seen")
	}

	err = os.Remove(dbconfig.filename)
	if err != nil {
		t.Error(err.Error())
	}
}
