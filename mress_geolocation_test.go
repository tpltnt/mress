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

func Test_serverLookupCoordinates_3(t *testing.T) {
	ipstring := "127.0.0.1"
	logger := createLogger("")
	port, err := readConfigInt("test2.ini", "geolocation", "port", logger)
	if err != nil {
		t.Log("please host your own lookup service https://github.com/fiorix/freegeoip")
		t.Error(err.Error())
	}

	lat, lon, err := serverLookupCoordinates(ipstring, "", port)
	if err == nil {
		t.Error("empty server string not dectected")
	}
	if lat != 0 {
		t.Error("wrong latitude returned")
	}
	if lon != 0 {
		t.Error("wrong longitude returned")
	}
}

func Test_serverLookupCoordinates_4(t *testing.T) {
	ipstring := "127.0.0.1"
	logger := createLogger("")
	server, err := readConfigString("test2.ini", "geolocation", "server", logger)
	if err != nil {
		t.Log("please host your own lookup service https://github.com/fiorix/freegeoip")
		t.Error(err.Error())
	}

	lat, lon, err := serverLookupCoordinates(ipstring, server, 0)
	if err == nil {
		t.Error("invalid port not dectected")
	}
	if lat != 0 {
		t.Error("wrong latitude returned")
	}
	if lon != 0 {
		t.Error("wrong longitude returned")
	}
}

func Test_haversin_0(t *testing.T) {
	// 0.2298488470659301412995316962785116981338447896910388...
	if 0.2298488470659301412995316962785116981338447896910388 > haversin(1.0) {
		t.Error("haversine returned too small value")
	}
	if 0.2298488470659301412995316962785116981338447896910389 < haversin(1.0) {
		t.Error("haversine returned too large value")
	}
}

func Test_calcDistance_0(t *testing.T) {
	lat1 := 31.2
	lon1 := 121.5
	lat2 := 31.2
	lon2 := 121.5
	result, err := calcDistance(lat1, lon1, lat2, lon2)
	if err != nil {
		t.Error(err.Error())
	}
	if result != 0.0 {
		t.Error("calculated wrong distance")
	}
}

func Test_calcDistance_1(t *testing.T) {
	lat1 := 90.0
	lon1 := 121.5
	lat2 := 31.2
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_calcDistance_2(t *testing.T) {
	lat1 := 90.1
	lon1 := 121.5
	lat2 := 31.2
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err == nil {
		t.Error("did not catch latitude 1 outside valid range")
	}
}

func Test_calcDistance_3(t *testing.T) {
	lat1 := -90.0
	lon1 := 121.5
	lat2 := 31.2
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_calcDistance_4(t *testing.T) {
	lat1 := -90.1
	lon1 := 121.5
	lat2 := 31.2
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err == nil {
		t.Error("did not catch latitude 1 outside valid range")
	}
}

func Test_calcDistance_5(t *testing.T) {
	lat1 := 31.2
	lon1 := 121.5
	lat2 := 90.0
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_calcDistance_6(t *testing.T) {
	lat1 := 31.2
	lon1 := 121.5
	lat2 := 90.1
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err == nil {
		t.Error("did not catch latitude 2 outside valid range")
	}
}

func Test_calcDistance_7(t *testing.T) {
	lat1 := 31.2
	lon1 := 121.5
	lat2 := -90.0
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_calcDistance_8(t *testing.T) {
	lat1 := 31.2
	lon1 := 121.5
	lat2 := -90.1
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err == nil {
		t.Error("did not catch latitude 2 outside valid range")
	}
}

func Test_calcDistance_9(t *testing.T) {
	lat1 := 31.2
	lon1 := 180.0
	lat2 := 31.2
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_calcDistance_10(t *testing.T) {
	lat1 := 31.2
	lon1 := 180.1
	lat2 := 31.2
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err == nil {
		t.Error("did not catch longitude 1 outside valid range")
	}
}

func Test_calcDistance_11(t *testing.T) {
	lat1 := 31.2
	lon1 := -180.0
	lat2 := 31.2
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_calcDistance_12(t *testing.T) {
	lat1 := 31.2
	lon1 := -180.1
	lat2 := 31.2
	lon2 := 121.5
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err == nil {
		t.Error("did not catch longitude 2 outside valid range")
	}
}

func Test_calcDistance_13(t *testing.T) {
	lat1 := 31.2
	lon1 := 121.5
	lat2 := 31.2
	lon2 := 180.0
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_calcDistance_14(t *testing.T) {
	lat1 := 31.2
	lon1 := 121.5
	lat2 := 31.2
	lon2 := 180.1
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err == nil {
		t.Error("did not catch longitude 2 outside valid range")
	}
}

func Test_calcDistance_15(t *testing.T) {
	lat1 := 31.2
	lon1 := 121.5
	lat2 := 31.2
	lon2 := -180.0
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_calcDistance_16(t *testing.T) {
	lat1 := 31.2
	lon1 := 121.5
	lat2 := 31.2
	lon2 := -180.1
	_, err := calcDistance(lat1, lon1, lat2, lon2)
	if err == nil {
		t.Error("did not catch longitude 2 outside valid range")
	}
}
