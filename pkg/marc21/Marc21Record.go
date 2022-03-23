package gomarc21

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type Record struct {
	Raw           []byte
	Leader        Leader
	Entries       []DirectoryEntry
	ControlFields []ControlField
	DataFields    []DataField
}

func (rec Record) GetRaw() string {
	return string(rec.Raw)
}

func (rec Record) GetMrk() string {
	str := rec.Leader.String()
	for _, cf := range rec.ControlFields {
		str = str + cf.String()
	}
	for _, df := range rec.DataFields {
		str = str + df.String()
	}
	return str
}

// ControlNum returns the control number (tag 001) for the record.
func (rec Record) ControlNum() string {
	for _, cf := range rec.ControlFields {
		if cf.Tag.GetTag() == "001" {
			return cf.Data
		}
	}
	return ""
}

// ControlNum returns the control number (tag 001) for the record.
func (rec Record) Title() string {
	for _, cf := range rec.DataFields {
		if cf.Tag.GetTag() == "245" {
			return cf.String()
		}
	}
	return ""
}

// parse directory of a marc record
func ParseDirectory(rawRec []byte) (dir []DirectoryEntry, err error) {
	log.Print("Parse Directory")
	for i := LEADER_LEN; rawRec[i] != END_OF_FIELD; i += 12 {
		var entry DirectoryEntry

		entry, err = NewDirectoryEntry(rawRec[i : i+12])
		if err != nil {
			return nil, err
		}
		log.Print(entry)

		dir = append(dir, entry)
	}
	log.Print("Parse Directory end")
	return dir, nil
}

// parse leader of a marc record
func ParseLeader(rawRec []byte) (Leader Leader, err error) {
	leader, err := NewLeader(rawRec)

	if err != nil {
		err = fmt.Errorf("invalid record length: %s", err)
		return leader, err
	}
	return leader, nil
}

// parse control fields from the raw MARC record bytes
func ParseControlFields(rawRec []byte, baseAddress int, dir []DirectoryEntry) (cfs []ControlField, err error) {

	var parseErrorTags []string
	var controlNumber string

	for _, d := range dir {
		isControlTag, err := d.Tag.IsControlTag()
		if err != nil {
			return nil, errors.New("issues in the control tags")
		}
		if isControlTag {
			start := baseAddress + d.StartingPosition
			b := rawRec[start : start+d.FieldLength]

			if b[len(b)-1] == END_OF_FIELD {
				if d.Tag.GetTag() == "001" {
					if controlNumber == "" {
						controlNumber = string(b[:len(b)-1])
					} else {
						parseErrorTags = append(parseErrorTags, d.Tag.GetTag())
					}
				}
				cfs = append(cfs, ControlField{Tag: d.Tag, Data: string(b[:len(b)-1])})
			} else {
				parseErrorTags = append(parseErrorTags, d.Tag.GetTag())
			}
		}
	}

	if len(parseErrorTags) > 0 {
		badTags := strings.Join(parseErrorTags, ", ")
		log.Printf("Control fields extraction error for ControlNumber %q (fields: %s)\n", controlNumber, badTags)
	}

	return cfs, nil
}

// extractDatafields extracts the data fields/sub-fields from the raw MARC record bytes
func ParseDataFields(rawRec []byte, baseAddress int, dir []DirectoryEntry) (dfs []DataField, err error) {
	for _, d := range dir {
		isControlTag, err := d.Tag.IsControlTag()
		if err != nil {
			return nil, errors.New("wrong field type")
		}

		if !isControlTag {
			start := baseAddress + d.StartingPosition
			b := rawRec[start : start+d.FieldLength]

			if b[len(b)-1] != END_OF_FIELD {
				return nil, errors.New("extractDatafields: Field terminator not found at end of field")
			}

			df := DataField{
				Tag:        d.Tag,
				Indicator1: string(b[0]),
				Indicator2: string(b[1]),
			}

			for _, t := range bytes.Split(b[2:d.FieldLength-1], []byte{SUBFIELD_INDICATOR}) {
				if len(t) > 0 {
					df.SubFields = append(df.SubFields, SubField{Code: string(t[0]), Value: string(t[1:])})
				}
			}
			dfs = append(dfs, df)
		}
	}

	return dfs, nil
}

// GetDatafields returns datafields for the record that match the
// specified comma separated list of tags. If no tags are specified
// (empty string) then all datafields are returned
func (rec Record) GetDatafields(tags string) (dfs []DataField) {
	if tags == "" {
		return rec.DataFields
	}

	for _, t := range strings.Split(tags, ",") {
		for _, df := range rec.DataFields {
			if df.Tag.GetTag() == t {
				dfs = append(dfs, df)
			}
		}
	}
	return dfs
}

// GetSubfields returns subfields for the datafield that match the
// specified codes. If no codes are specified (empty string) then all

// ParseRecord takes the bytes for a MARC record and returns the parsed
// record structure
func ParseRecord(rawRec []byte) (rec Record, err error) {

	log.Print(rawRec)
	rec = Record{}

	rec.Leader, err = ParseLeader(rawRec[:24])
	if err != nil {
		log.Print(err)
		return rec, err
	}
	log.Print(rec.Leader.String())
	log.Print("Leader parsed ok")

	rec.Entries, err = ParseDirectory(rawRec)
	if err != nil {
		log.Print(err)
		return rec, err
	}
	log.Print("Directory Entries parsed ok")
	for index, entry := range rec.Entries {
		log.Print(index, entry)
	}

	baseDataAddress := rec.Leader.BaseAddressOfData
	log.Print("base Data address is okay")

	log.Print(baseDataAddress)
	log.Print(rawRec)
	rec.ControlFields, err = ParseControlFields(rawRec, baseDataAddress, rec.Entries)
	if err != nil {
		return rec, err
	}
	log.Print("control fields is okay")
	for _, cf := range rec.ControlFields {
		log.Print(cf.String())
	}
	log.Print("control fields is iiiokay")

	rec.DataFields, err = ParseDataFields(rawRec, baseDataAddress, rec.Entries)
	if err != nil {
		return rec, err
	}
	log.Print("data fields is okay")
	for _, df := range rec.DataFields {
		log.Print(df.String())
	}
	log.Print("control fields is iiiokay")

	return rec, nil
}

// ReadRecord returns a single MARC record from a reader.
func ReadRecord(reader io.Reader) (record Record, err error) {
	record = Record{}
	// read 5 bytes into a leader
	data := make([]byte, 5)
	n, err := io.ReadFull(reader, data)
	if err != nil {
		return record, errors.New("reading the record length of a leader errors")
	}
	if n != 5 {
		return record, errors.New("the record is too short")
	}

	recordLength, err := strconv.Atoi(string(data[0:5]))
	log.Print(recordLength)
	if err != nil {
		return record, errors.New("the record length is wrong")
	}
	tmpRecord := make([]byte, recordLength-5)

	n, err = io.ReadFull(reader, tmpRecord)
	if err != nil {
		log.Print(err)
		return record, errors.New("reading a whole record errors")
	}

	if n != recordLength-5 {
		log.Print(err)
		return record, errors.New("the record is too short not a complete record")
	}
	log.Print(tmpRecord)

	rawRecord := append(data, tmpRecord...)
	log.Print(rawRecord)

	record, err = ParseRecord(rawRecord)
	if err != nil {
		return record, err
	}

	return record, nil
}
