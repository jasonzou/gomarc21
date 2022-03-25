package gomarc21

import (
	"errors"
	"fmt"
)

func NewTagByStr(tagStr string) (Tag, error) {
	tagBytes := []byte(tagStr)
	tag, err := NewTag(tagBytes)
	return tag, err
}

func NewTag(tagBytes []byte) (Tag, error) {
	if len(tagBytes) != 3 {
		msg := fmt.Sprintf("The length of a tag is 3, but found %d instead", len(tagBytes))
		err := errors.New(msg)
		return Tag(""), err
	}

	var tag Tag
	// need to check each byte
	tag1 := tagBytes[0]
	tag2 := tagBytes[1]
	tag3 := tagBytes[2]
	if tag1 >= '0' && tag1 <= 'z' {
		if tag2 >= '0' && tag2 <= 'z' {
			if tag3 >= '0' && tag3 <= 'z' {
				tag = Tag(string(tagBytes))
				//fmt.Sprint("%c%c%c", tag1, tag2, tag3)
				return tag, nil
			}
		}
	}
	err := errors.New("the tag string is invalid")
	return Tag(""), err
}

func (tag Tag) GetTag() string {
	return string(tag)
}

func (tag Tag) IsControlTag() (bool, error) {
	if len(string(tag)) != 3 {
		msg := fmt.Sprintf("The length of a tag is 3, but found %d instead", len(tag))
		err := errors.New(msg)
		return false, err
	}

	tag1 := tag[0]
	tag2 := tag[1]
	tag3 := tag[2]
	// <xsd:pattern value="00[1-9A-Za-z]{1}"/>
	// https://www.loc.gov/standards/marcxml/schema/MARC21slim.xsd
	if tag1 == '0' && tag2 == '0' {
		if tag3 >= '1' && tag3 <= 'z' {
			return true, nil
		} else {
			return false, errors.New("invalid tag")
		}
	}
	return false, nil
}

func (tag Tag) IsDataTag() (bool, error) {
	if len(string(tag)) != 3 {
		msg := fmt.Sprintf("The length of a tag is 3, but found %d instead", len(tag))
		err := errors.New(msg)
		return false, err
	}

	tag1 := tag[0]
	tag2 := tag[1]
	tag3 := tag[2]
	// <xsd:pattern value="(0([1-9A-Z][0-9A-Z])|0([1-9a-z][0-9a-z]))|(([1-9A-Z][0-9A-Z]{2})|([1-9a-z][0-9a-z]{2}))"/>
	// https://www.loc.gov/standards/marcxml/schema/MARC21slim.xsd

	// 0([1-9A-Z][0-9A-Z]
	// 0([1-9a-z][0-9a-z])
	if tag1 >= '0' && tag1 <= 'z' {
		if tag2 >= '1' && tag2 <= 'z' {
			if tag3 >= '0' && tag3 <= 'z' {
				return true, nil
			}
		}
	}

	// [1-9A-Z][0-9A-Z]{2}
	// [1-9a-z][0-9a-z]{2}
	if tag1 >= '1' && tag1 <= 'z' {
		if tag2 >= '0' && tag2 <= 'z' {
			if tag3 >= '0' && tag3 <= 'z' {
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
