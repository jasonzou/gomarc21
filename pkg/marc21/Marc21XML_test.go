package gomarc21

import (
	"os"
	"testing"
)

func TestRecordAsXML(test *testing.T) {
	data, err := os.Open("data/test_1a.mrc")
	if err != nil {
		test.Fatal(err)
	}
	defer data.Close()

	var rec Record
	rec, err = ParseNextRecord(data)
	if rec.Leader.RecordLength != 0 {
		test.Error(rec.RecordAsXML())
	}
}
