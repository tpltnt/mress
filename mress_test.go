package main

import (
	"testing"
)

// test stdout destination
func Test_create_Logger_1(t *testing.T) {
	dest := "stdout"
	logger := createLogger(&dest)
	if logger == nil {
		t.Error("creating logger to '" + dest + "' returned 'nil'")
	}
}

// test stderr destination
func Test_create_Logger_2(t *testing.T) {
	dest := "stderr"
	logger := createLogger(&dest)
	if logger == nil {
		t.Error("creating logger to '" + dest + "' returned 'nil'")
	}
}
