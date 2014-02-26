package main

import "testing"

func Test_serverLookupCoordinates_0(t *testing.T) {
	ipstring := "207.241.224.2"
	logger := createLogger("")
	server, err := readConfigString("test2.ini", "geolocation", "server", logger)
	if err != nil {
		t.Log("please host your own lookup service https://github.com/fiorix/freegeoip")
		t.Error(err.Error())
	}
	port, err := readConfigInt("test2.ini", "geolocation", "port", logger)
	if err != nil {
		t.Log("please host your own lookup service https://github.com/fiorix/freegeoip")
		t.Error(err.Error())
	}
	lat, lon, err := serverLookupCoordinates(ipstring, server, port)
	if err != nil {
		t.Log("please host your own lookup service https://github.com/fiorix/freegeoip")
		t.Error(err.Error())
	}
	if lat != 37.7811 {
		t.Log(lat)
		t.Error("wrong latitude returned")
	}
	if lon != -122.4625 {
		t.Log(lon)
		t.Error("wrong longitude returned")
	}
}

func Test_serverLookupCoordinates_1(t *testing.T) {
	ipstring := "127.0.0.1"
	logger := createLogger("")
	server, err := readConfigString("test2.ini", "geolocation", "server", logger)
	if err != nil {
		t.Log("please host your own lookup service https://github.com/fiorix/freegeoip")
		t.Error(err.Error())
	}
	port, err := readConfigInt("test2.ini", "geolocation", "port", logger)
	if err != nil {
		t.Log("please host your own lookup service https://github.com/fiorix/freegeoip")
		t.Error(err.Error())
	}
	lat, lon, err := serverLookupCoordinates(ipstring, server, port)
	if err != nil {
		t.Log("please host your own lookup service https://github.com/fiorix/freegeoip")
		t.Error(err.Error())
	}
	if lat != 0 {
		t.Error("wrong latitude returned")
	}
	if lon != 0 {
		t.Error("wrong longitude returned")
	}
}

func Test_serverLookupCoordinates_2(t *testing.T) {
	ipstring := ""
	logger := createLogger("")
	server, err := readConfigString("test2.ini", "geolocation", "server", logger)
	if err != nil {
		t.Log("please host your own lookup service https://github.com/fiorix/freegeoip")
		t.Error(err.Error())
	}
	port, err := readConfigInt("test2.ini", "geolocation", "port", logger)
	if err != nil {
		t.Log("please host your own lookup service https://github.com/fiorix/freegeoip")
		t.Error(err.Error())
	}
	lat, lon, err := serverLookupCoordinates(ipstring, server, port)
	if err == nil {
		t.Error("empty ip string not dectected")
	}
	if lat != 0 {
		t.Error("wrong latitude returned")
	}
	if lon != 0 {
		t.Error("wrong longitude returned")
	}
}
