//package gomarc21
package main

import (
	"errors"
	"fmt"
	"reflect"
)

type Tag struct {
	Tag1 byte // for both Control and Data fields
	Tag2 byte // for Data fields
	Tag3 byte // for Data fields
}

func NewTag(tagStr string) (error, Tag){
	tagBytes := []byte(tagStr)
	err, tag = NewTag(tagBytes)
	return err, tag
}

func NewTag(tagBytes []bytes) (error, Tag){
	if len(tagBytes) != 3 {
		msg := fmt.Sprintf("The length of a tag is 3, but found %d instead", len(tagBytes))
		err := error.New(msg)
		return err, nil
	}
	var tag Tag
	// need to check each byte 
	tag.Tag1 = tagBytes[0]
	tag.Tag2 = tagBytes[1]
	tag.Tag3 = tagBytes[2]
	if tag.Tag1 >= '0' && tag.Tag1 <= 'z' {
		if tag.Tag2 >= '0' && tag.Tag2 <= 'z' {
			if tag.Tag3 >= '0' && tag.Tag3 <= 'z' {
				return nil, tag
			}
		}
	}
	err := error.New("The tag string is invalid.")
	return err, nil
}

func (tag *Tag) GetTag() string {
	myTag := fmt.Sprintf("%c%c%c", tag.Tag1, tag.Tag2, tag.Tag3)
	return myTag
}

func (tag *Tag) IsControlTag() (error, bool) {
	// <xsd:pattern value="00[1-9A-Za-z]{1}"/>
	// https://www.loc.gov/standards/marcxml/schema/MARC21slim.xsd
	if tag.Tag1 == '0' && tag.Tag2 == '0' {
		if tag.Tag3 >= '1' && tag.Tag3 <= 'z' {
			return nil, true
		} else {
			return errors.New("Invalid tag"), false
		}
	}
	return nil, false
}

func (tag *Tag) IsDataTag() (error, bool) {
	// <xsd:pattern value="(0([1-9A-Z][0-9A-Z])|0([1-9a-z][0-9a-z]))|(([1-9A-Z][0-9A-Z]{2})|([1-9a-z][0-9a-z]{2}))"/>
	// https://www.loc.gov/standards/marcxml/schema/MARC21slim.xsd

	// 0([1-9A-Z][0-9A-Z]
	// 0([1-9a-z][0-9a-z])
	if tag.Tag1 >= '0' && tag.Tag1 <= 'z' {
		if tag.Tag2 >= '1' && tag.Tag2 <= 'z' {
			if tag.Tag3 >= '0' && tag.Tag3 <= 'z' {
				return nil, true
			}
		}
	}

	// [1-9A-Z][0-9A-Z]{2}
	// [1-9a-z][0-9a-z]{2}
	if tag.Tag1 >= '1' && tag.Tag1 <= 'z' {
		if tag.Tag2 >= '0' && tag.Tag2 <= 'z' {
			if tag.Tag3 >= '0' && tag.Tag3 <= 'z' {
				return nil, true
			}
		}
	}

	err, controlTag := tag.IsControlTag()
	if err == nil {
		if controlTag == false {
			return errors.New("Invalid data tag"), false
		} else {
			return nil, false
		}
	}

	return errors.New("Invalid data tag"), false
}