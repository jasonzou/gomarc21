package gomarc21

import (
	"fmt"
	"testing"
)

func TestGettag1(test *testing.T) {
	t, err := NewTagByStr("1234")
	if err != nil {
		test.Error("derorr")
	}

	if t.GetTag() != "" {
		test.Error("Tag no should be 123, but ", t.GetTag())
	}
}

func TestIsControlTag(test *testing.T) {
}

func TestIsDataTag(test1 *testing.T) {

	var t Tag
	t, _ = NewTagByStr("123")

	fmt.Println(t.GetTag())
	test, err := t.IsControlTag()
	if err == nil {
		fmt.Println(test)
	}
	test, err = t.IsDataTag()
	if err == nil {
		fmt.Println(test, "data tag")
	}
	test, err = t.IsControlTag()
	if err == nil {
		fmt.Println(test)
	}
	test, err = t.IsDataTag()
	fmt.Println(t.GetTag())
	if err == nil {
		fmt.Println(test, "data tag")
	} else {
		fmt.Println(err)
	}
}
