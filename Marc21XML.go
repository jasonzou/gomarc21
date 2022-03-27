package gomarc21

import (
	"encoding/xml"
	"fmt"
)

func (record Record) RecordAsXML() (string, error) {
	x := XmlRecord{
		Leader:        string(record.Leader.raw),
		ControlFields: record.ControlFields,
		DataFields:    record.DataFields,
	}

	b, err := xml.MarshalIndent(x, "", "")
	fmt.Print(string(b))

	return string(b), err
}
