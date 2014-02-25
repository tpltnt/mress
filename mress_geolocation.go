package main

import (
	"fmt"
	//	"github.com/thoj/go-ircevent" // imported as "irc"
	"regexp"
	"strings"
)

// Convert given IPv4 string (dotted quad) to geo coordinates.
// The channels transmits first latitude, then longitude. An error is in
func serverLookupCoordinates(ip string, server string, port int) (lat, lon float32, err error) {
	if len(ip) < 8 {
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
	return 0.0, 0.0, fmt.Errorf("not fully implemented yet")
}
