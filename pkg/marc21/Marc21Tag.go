package gomarc21

import (
	"errors"
	"fmt"
)

type Tag struct {
	Tag1 byte // for both Control and Data fields
	Tag2 byte // for Data fields
	Tag3 byte // for Data fields
}

func NewTagByStr(tagStr string) (Tag, error) {
	tagBytes := []byte(tagStr)
	tag, err := NewTag(tagBytes)
	return tag, err
}

func NewTag(tagBytes []byte) (Tag, error) {
	if len(tagBytes) != 3 {
		msg := fmt.Sprintf("The length of a tag is 3, but found %d instead", len(tagBytes))
		err := errors.New(msg)
		return Tag{'0', '0', '0'}, err
	}

	var tag Tag
	// need to check each byte
	tag.Tag1 = tagBytes[0]
	tag.Tag2 = tagBytes[1]
	tag.Tag3 = tagBytes[2]
	if tag.Tag1 >= '0' && tag.Tag1 <= 'z' {
		if tag.Tag2 >= '0' && tag.Tag2 <= 'z' {
			if tag.Tag3 >= '0' && tag.Tag3 <= 'z' {
				return tag, nil
			}
		}
	}
	err := errors.New("the tag string is invalid")
	return Tag{'0', '0', '0'}, err
}

func (tag Tag) GetTag() string {
	myTag := fmt.Sprintf("%c%c%c", tag.Tag1, tag.Tag2, tag.Tag3)
	return myTag
}

func (tag Tag) IsControlTag() (bool, error) {
	// <xsd:pattern value="00[1-9A-Za-z]{1}"/>
	// https://www.loc.gov/standards/marcxml/schema/MARC21slim.xsd
	if tag.Tag1 == '0' && tag.Tag2 == '0' {
		if tag.Tag3 >= '1' && tag.Tag3 <= 'z' {
			return true, nil
		} else {
			return false, errors.New("invalid tag")
		}
	}
	return false, nil
}

func (tag Tag) IsDataTag() (bool, error) {
	// <xsd:pattern value="(0([1-9A-Z][0-9A-Z])|0([1-9a-z][0-9a-z]))|(([1-9A-Z][0-9A-Z]{2})|([1-9a-z][0-9a-z]{2}))"/>
	// https://www.loc.gov/standards/marcxml/schema/MARC21slim.xsd

	// 0([1-9A-Z][0-9A-Z]
	// 0([1-9a-z][0-9a-z])
	if tag.Tag1 >= '0' && tag.Tag1 <= 'z' {
		if tag.Tag2 >= '1' && tag.Tag2 <= 'z' {
			if tag.Tag3 >= '0' && tag.Tag3 <= 'z' {
				return true, nil
			}
		}
	}

	// [1-9A-Z][0-9A-Z]{2}
	// [1-9a-z][0-9a-z]{2}
	if tag.Tag1 >= '1' && tag.Tag1 <= 'z' {
		if tag.Tag2 >= '0' && tag.Tag2 <= 'z' {
			if tag.Tag3 >= '0' && tag.Tag3 <= 'z' {
				return true, nil
			}
		}
	}

	controlTag, err := tag.IsControlTag()
	if err == nil {
		if controlTag {
			return false, nil
		} else {
			return false, errors.New("invalid data tag")
		}
	}

	return false, errors.New("invalid data tag")
}
