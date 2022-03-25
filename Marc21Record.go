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

// RecordAsMARC converts a Record into a MARC record byte array
func (rec Record) RecordAsMARC() (marc []byte, err error) {

	if rec.Leader.String() == "" {
		err = errors.New("record Leader is undefined")
		return marc, err
	}

	var dir []DirectoryEntry
	var rawDir []byte
	var cfs []byte
	var dfs []byte
	var startPos int

	// Pack the control fields
	for _, cf := range rec.ControlFields {
		if cf.String() == "" {
			continue
		}

		b := []byte(cf.String())
		b = append(b, END_OF_FIELD)
		cfs = append(cfs, b...)

		dir = append(dir, DirectoryEntry{Tag: cf.Tag, StartingPosition: startPos, FieldLength: len(b)})

		startPos += len(b)
	}

	// Pack the data fields/sub-fields
	for _, df := range rec.DataFields {

		ind1 := df.GetIndicator1()
		ind2 := df.GetIndicator2()

		b := []byte(ind1)
		b = append(b, []byte(ind2)...)

		for _, sf := range df.SubFields {
			b = append(b)
			b = append(b, []byte(sf.GetCode())...)
			b = append(b, []byte(sf.GetValue())...)
		}
		b = append(b, SUBFIELD_INDICATOR)

		dfs = append(dfs, b...)

		dir = append(dir, DirectoryEntry{Tag: df.Tag, StartingPosition: startPos, FieldLength: len(b)})

		startPos += len(b)
	}

	// Generate the directory
	for _, de := range dir {
		rawDir = append(rawDir, []byte(de.Tag.GetTag())...)
		rawDir = append(rawDir, []byte(fmt.Sprintf("%04d", de.FieldLength))...)
		rawDir = append(rawDir, []byte(fmt.Sprintf("%05d", de.StartingPosition))...)
	}
	rawDir = append(rawDir, SUBFIELD_INDICATOR)

	// Build the leader
	recLen := []byte(fmt.Sprintf("%05d", 24+len(rawDir)+len(cfs)+len(dfs)+1))
	recBaseDataAddress := []byte(fmt.Sprintf("%05d", 24+len(rawDir)))

	ldr := []byte(rec.Leader.String())
	for i := 0; i <= 4; i++ {
		ldr[i] = recLen[i]
		ldr[i+12] = recBaseDataAddress[i]
	}

	// Final assembly
	marc = append(marc, ldr...)
	marc = append(marc, rawDir...)
	marc = append(marc, cfs...)
	marc = append(marc, dfs...)
	marc = append(marc, END_OF_RECORD)

	return marc, nil
}

// TODO: Add support for JSONL (JSON line delimited) format that makes JSON
// easier to parse with Unix tools like grep, tail, and so on.
func (record Record) RecordAsJson() (jsonstr string, err error) {
	jsonstr = record.Leader.String()
	for _, cf := range record.ControlFields {
		jsonstr += cf.String()
	}
	for _, df := range record.DataFields {
		jsonstr = jsonstr + df.String()
	}
	return jsonstr, nil
}

// GetSubFields returns a slice of subfields that match the given tag
// and code.
//func (record *Record) GetSubFields(tag string, code byte) (subfields []*SubField) {
/*
	subfields = make([]*SubField, 0, 4)
	fields := record.GetFields(tag)
	for _, field := range fields {
		switch data := field.(type) {
		case *DataField:
			for _, subfield := range data.SubFields {
				if subfield.Code == code {
					subfields = append(subfields, subfield)
				}
			}
		}
	}
	return */
//}
