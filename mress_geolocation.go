package main

import (
	"fmt"
	//	"github.com/thoj/go-ircevent" // imported as "irc"
	"bufio"
	"encoding/json"
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Convert given IPv4 string (dotted quad) to geo coordinates.
// The channels transmits first latitude, then longitude.
// The IP lookup is done by a freegeoip service: https://github.com/fiorix/freegeoip
func serverLookupCoordinates(ip string, server string, port int) (lat, lon float32, err error) {
	if len(ip) < 7 {
		return 0.0, 0.0, fmt.Errorf("given IPv4 too short")
	}
	if strings.Count(ip, ".") != 3 {
		return 0.0, 0.0, fmt.Errorf("given IPv4 doesn't contain 4 dots, %d")
	}
	r, err := regexp.Compile(`\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf(err.Error())
	}
	if r.MatchString(ip) != true {
		return 0.0, 0.0, fmt.Errorf("IPv4 regex did not match")
	}
	if len(server) == 0 {
		return 0.0, 0.0, fmt.Errorf("given server string is empty")
	}
	if port == 0 {
		return 0.0, 0.0, fmt.Errorf("given port is invalid")
	}
	serverstring := server + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", serverstring, 30*time.Second)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("server connection failed")
	}
	defer conn.Close()
	lookupstring := "GET /json/" + ip + " HTTP/1.1\r\n\r\n"
	fmt.Fprintf(conn, lookupstring)
	reader := bufio.NewReader(conn)
	status, err := reader.ReadString('\n')
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("reading status from buffer failed")
	}
	// dispose uninteresting response lines
	for i := 0; i < 5; i++ {
		_, _ = reader.ReadString('\n')
	}
	// line with json data doesn't end with \n
	jsonstring, err := reader.ReadString('}')
	if err != nil {
		return 0.0, 0.0, fmt.Errorf(err.Error())
	}
	if strings.Contains(status, "404 Not Found") {
		return 0.0, 0.0, fmt.Errorf("Ressource not found (404)")
	}
	// decode JSON data
	type Geodata struct {
		Ip, Country_code, Country_name          string
		Region_code, Region_name, City, Zipcode string
		Latitude, Longitude                     float64
		Metro_code, Areacode                    string
	}
	dec := json.NewDecoder(strings.NewReader(jsonstring))
	var gip Geodata
	err = dec.Decode(&gip)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf(err.Error())
	}

	return float32(gip.Latitude), float32(gip.Longitude), nil
}

// Haversin math helper function.
func haversin(value float64) float64 {
	return math.Sin(value/2) * math.Sin(value/2)
}

// Calculate the distance between two (GPS) points using the
// Haversine formula. This yields an error of ca. 0.5%.
// Using Vincenty's formulae yields better precision of needed.
func calcDistance(lat1, lon1, lat2, lon2 float32) (float32, error) {
	// first haversin term
	//	fmt.Println(term1)
	return 0.0, fmt.Errorf("not implemented yet")
}
