package gomarc21

import (
	"os"
	"testing"
)

func openTestMARC(t *testing.T) (data *os.File) {
	data, err := os.Open("data/test.mrc")
	if err != nil {
		t.Fatal(err)
	}
	return
}

func TestGetRaw1(test *testing.T) {
}

// ControlNum returns the control number (tag 001) for the record.
func TestControlNum(test *testing.T) {
}

// ControlNum returns the control number (tag 001) for the record.
func TestTitle(test *testing.T) {
}
