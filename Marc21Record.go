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

func (rec Record) GetRaw() string {
	return string(rec.raw)
}

func (rec Record) RecordAsMrk() (string, error) {
	return rec.GetMrk(), nil
}

func (rec Record) GetMrk() string {
	str := rec.Leader.String() + "\n"
	for _, cf := range rec.ControlFields {
		str = str + cf.String() + "\n"
	}
	for _, df := range rec.DataFields {
		str = str + df.String() + "\n"
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
	for i := LEADER_LEN; rawRec[i] != END_OF_FIELD; i += 12 {
		var entry DirectoryEntry

		entry, err = NewDirectoryEntry(rawRec[i : i+12])
		if err != nil {
			return nil, err
		}
		dir = append(dir, entry)
	}
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
					df.SubFields = append(df.SubFields, SubField{Code: string(t[0]), Data: string(t[1:])})
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

	rec = Record{}

	rec.Leader, err = ParseLeader(rawRec[:24])
	if err != nil {
		log.Print(err)
		return rec, err
	}

	rec.Entries, err = ParseDirectory(rawRec)
	if err != nil {
		return rec, err
	}
	for index, entry := range rec.Entries {
		log.Print(index, entry)
	}

	baseDataAddress := rec.Leader.BaseAddressOfData
	rec.ControlFields, err = ParseControlFields(rawRec, baseDataAddress, rec.Entries)
	if err != nil {
		return rec, err
	}

	rec.DataFields, err = ParseDataFields(rawRec, baseDataAddress, rec.Entries)
	if err != nil {
		return rec, err
	}

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
	if err != nil {
		return record, errors.New("the record length is wrong")
	}
	tmpRecord := make([]byte, recordLength-5)

	n, err = io.ReadFull(reader, tmpRecord)
	if err != nil {
		return record, errors.New("reading a whole record errors")
	}

	if n != recordLength-5 {
		return record, errors.New("the record is too short not a complete record")
	}

	rawRecord := append(data, tmpRecord...)

	record, err = ParseRecord(rawRecord)
	if err != nil {
		return record, err
	}

	return record, nil
}

// NextRecord reads the next MARC record and returns the unparsed bytes
func NextRecord(reader io.Reader) (rawRec []byte, err error) {

	// Read the first 5 bytes, determine the record length and
	//    read the remainder of the record
	rawLen := make([]byte, 5)
	_, err = reader.Read(rawLen)
	if err != nil {
		return nil, err
	}

	recLen, err := strconv.Atoi(string(rawLen[0:5]))
	if err != nil {
		return nil, err
	}

	// Ensure that we have a "sane" record length?
	if recLen <= LEADER_LEN {
		err = errors.New("MARC record is too short")
		return nil, err
	} else if recLen > MAX_RECORD_LEN {
		err = errors.New("MARC record is too long")
		return nil, err
	}

	rawRec = make([]byte, recLen)
	// ensure that the raw len is available for the leader
	copy(rawRec, rawLen)

	// Read the remainder of the record
	_, err = reader.Read(rawRec[5:recLen])
	if err != nil {
		return nil, err
	}

	// The last byte should be a record terminator
	if rawRec[len(rawRec)-1] != END_OF_RECORD {
		return nil, errors.New("Record terminator not found at end of record")
	}
	// other issues?

	return rawRec, nil
}

// ParseNextRecord reads the next MARC record and returns the parsed
// record structure
func ParseNextRecord(reader io.Reader) (record Record, err error) {

	rawRec, err := NextRecord(reader)
	if err != nil {
		return Record{}, err
	}

	record, err = ParseRecord(rawRec)
	if err != nil {
		return Record{}, err
	}

	return record, nil
}
