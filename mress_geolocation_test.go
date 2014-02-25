package main

import "testing"

func Test_serverLookupCoordinates_0(t *testing.T) {
	ipstring := "127.0.0.1"
	server := "foo.bar"
	port := 2384
	lat, lon, err := serverLookupCoordinates(ipstring, server, port)
	if err != nil {
		t.Error(err.Error())
	}
	if lat != 0 {
		t.Error("wrong latitude returned")
	}
	if lon != 0 {
		t.Error("wrong longitude returned")
	}
}
