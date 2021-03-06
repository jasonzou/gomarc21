package gomarc21

import (
	"testing"
)

func TestNewLeader(t *testing.T) {
	result, _ := NewLeader([]byte("00350cz  a2200157n  4500"))
	if result.RecordLength != 350 {
		t.Log("Record Length should be 350, but got", result.RecordLength)
		t.Error("Record Length is ", result.RecordLength)
	}
	if result.GetRaw() != "00350cz  a2200157n  4500" {
		t.Error("TestNewLeader is ", result.GetRaw())
	}

}
