package main

import "testing"

func Test_serverLookupCoordinates_0(t *testing.T) {
	ipstring := "127.0.0.1"
	logger := createLogger("")
	server, err := readConfigString("test2.ini", "geolocation", "server", logger)
	if err != nil {
		t.Error(err.Error())
	}
	port, err := readConfigInt("test2.ini", "geolocation", "port", logger)
	if err != nil {
		t.Error(err.Error())
	}
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
