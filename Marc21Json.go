package gomarc21

import (
	"encoding/json"
	"fmt"
)

func (record Record) RecordAsJson() (string, error) {
	j := JsonRecord{
		Leader:        string(record.Leader.raw),
		ControlFields: record.ControlFields,
		DataFields:    record.DataFields,
	}

	leaderJson, err := json.MarshalIndent(j.Leader, "", "")
	if err != nil {
		fmt.Printf("%s\r\n", err)
	}

	var fieldStr string
	fieldStr += "\"fields\": ["
	for _, cf := range j.ControlFields {
		//.MarshalIndent(cf)
		fieldStr += cf.AsJson() + ","
	}

	for _, df := range j.DataFields {
		//.MarshalIndent(cf)
		fieldStr += df.AsJson() + ","
	}
	fieldStr = fieldStr[:len(fieldStr)-1]
	fieldStr += "]"

	b := "{  \"leader\":" + string(leaderJson) + "," + fieldStr + "}"
	fmt.Print(string(b))

	return string(b), err
}
