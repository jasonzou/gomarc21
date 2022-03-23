package gomarc21

import (
	"testing"
)

// GetTag returns the tag for the datafield
func TestGetTag(test *testing.T) {
	var df DataField
	df.Tag, _ = NewTagByStr("020")
	if df.Tag != df.GetTag() {
		test.Error("df Tag error", df.Tag)
	}

}

// GetTag returns the tag for the datafield
func TestString(test *testing.T) {

}

// Contains returns true if the field contains the passed string.
func TestContains(test *testing.T) {
}

// GetInd1 returns the indicator 1 value for the datafield
func TestGetIndicator1(test *testing.T) {

}

// GetInd2 returns the indicator 2 value for the DataField
func TestGetIndicator2(test *testing.T) {

}

func TestFormatIndicator(test *testing.T) {
}
