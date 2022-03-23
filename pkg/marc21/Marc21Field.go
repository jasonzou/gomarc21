package gomarc21

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// MakeField creates a field objet with the data received.
func MakeDataField(tag string, data []byte) (DataField, error) {
	f := DataField{}
	tempTag, err := NewTagByStr(tag)
	if err == nil {
		f.Tag = tempTag
	} else {
		return f, fmt.Errorf("invalid tag %s", tag)
	}

	// It's a control field
	isControlTag, _ := f.Tag.IsControlTag()
	if isControlTag {
		return f, errors.New("use MakeControlField instead; control field - wrong type")
	}

	if len(data) > 2 {
		f.Indicator1 = string(data[0])
		f.Indicator2 = string(data[1])
	} else {
		return f, errors.New("invalid indicators detected")
	}

	if len(data) < 4 { // Each data field contains at least one subfield code.
		return f, errors.New("bad SubFields length")
	}

	for _, sf := range bytes.Split(data[3:], []byte{SUBFIELD_INDICATOR}) {
		if len(sf) > 1 {
			f.SubFields = append(f.SubFields, SubField{string(sf[0]), string(sf[1:])})
		}
	}
	return f, nil
}

// indicate what subfields are to be returned.
func (df DataField) GetSubFields(filter string) []SubField {
	values := []SubField{}
	for _, sub := range df.SubFields {
		if strings.Contains(filter, sub.Code) {
			value := SubField{
				Code:  sub.Code,
				Value: sub.Value,
			}
			values = append(values, value)
		}
	}
	return values
}
