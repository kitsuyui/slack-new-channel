package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	_, err := LoadChannelFromFile("./test_assets/sample.json")
	if err != nil {
		t.Fail()
	}
}
