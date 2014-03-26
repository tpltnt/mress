package main

import (
	//	"log"
	//	"os"
	//	"strconv"
	"testing"
)

func Test_markAsSeen_0(t *testing.T) {
	dbconfig := MressDbConfig{backend: "sqlite3", filename: "testoffline.db", offlineMsgTable: "messages", introductionTable: "intro"}
	user := "testuser"
	err := markAsSeen(dbconfig, user)
	if err != nil {
		t.Error(err.Error())
	}
}
