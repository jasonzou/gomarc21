package gomarc21

import (
	"fmt"
	"testing"
)

func TestGetTag(test *testing.T) {
	var t Tag
	t.Tag1 = '1'
	t.Tag2 = '2'
	t.Tag3 = '3'

	if t.GetTag() != "123" {
		test.Error("Tag no should be 123, but ", t.GetTag())
	}
}

func TestIsControlTag(test *testing.T) {
}

func TestIsDataTag(test *testing.T) {
}

func main() {
	var t Tag
	t.Tag1 = '1'
	t.Tag2 = '2'
	t.Tag3 = '3'

	fmt.Println(t.GetTag())
	err, test := t.IsControlTag()
	if err == nil {
		fmt.Println(test)
	}
	err, test = t.IsDataTag()
	if err == nil {
		fmt.Println(test, "data tag")
	}
	t.Tag1 = '0'
	t.Tag2 = '0'
	t.Tag3 = 'z'
	err, test = t.IsControlTag()
	if err == nil {
		fmt.Println(test)
	}
	t.Tag2 = '#'
	err, test = t.IsDataTag()
	fmt.Println(t.GetTag())
	if err == nil {
		fmt.Println(test, "data tag")
	} else {
		fmt.Println(err)
	}
}
