package gomarc21

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"

	marc "github.com/jasonzou/gomarc21/marc21.old"
)

// A MARC Record Set is for containing zero or more MARC records
type Collection struct {
	Name    xml.Name  `xml:"collection"`
	Records []*Record `xml:"record"`
}

type controlField struct {
	Tag   string `xml:"tag,attr"`
	Value string `xml:",chardata"`
}

type subfield struct {
	Code  string `xml:"code,attr"`
	Value string `xml:",chardata"`
}
type dataField struct {
	Tag       string     `xml:"tag,attr"`
	Ind1      string     `xml:"ind1,attr"`
	Ind2      string     `xml:"ind2,attr"`
	Subfields []subfield `xml:"subfield"`
}

type xmlRecord struct {
	XMLName       xml.Name       `xml:"record"`
	Leader        string         `xml:"leader"`
	ControlFields []controlField `xml:"controlfield"`
	DataFields    []dataField    `xml:"datafield"`
}

const xmlProlog = `<?xml version="1.0" encoding="UTF-8"?>`
const xmlRootBegin = `<collection xmlns="http://www.loc.gov/MARC21/slim" xmlns:marc="http://www.loc.gov/MARC21/slim">`
const xmlRootEnd = `</collection>`

func toXML(params ProcessFileParams) error {
	if count == 0 {
		return nil
	}

	file, err := os.Open(params.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Printf("%s\n%s\n", xmlProlog, xmlRootBegin)

	var i, out int
	marc := marc.NewMarcFile(file)
	for marc.Scan() {

		r, err := marc.Record()
		if err == io.EOF {
			break
		}

		if err != nil {
			printError(r, "PARSE ERROR", err)
			if params.debug {
				continue
			}
			return err
		}

		if i++; i < start {
			continue
		}

		if r.Contains(params.searchValue) && r.HasFields(params.hasFields) {
			str, err := recordToXML(r, params)
			if err != nil {
				if params.debug {
					printError(r, "XML PARSE ERROR", err)
					continue
				}
				panic(err)
			}
			fmt.Printf("%s\r\n", str)
			if out++; out == count {
				break
			}
		}
	}
	fmt.Printf("%s\n", xmlRootEnd)

	return marc.Err()
}

func recordToXML(r marc.Record, params ProcessFileParams) (string, error) {
	x := xmlRecord{
		Leader: r.Leader.Raw(),
	}

	for _, f := range r.Filter(params.filters, params.exclude) {
		if f.IsControlField() {
			x.ControlFields = append(x.ControlFields, controlField{Tag: f.Tag, Value: f.Value})
		} else {
			df := dataField{Tag: f.Tag, Ind1: f.Indicator1, Ind2: f.Indicator2}
			for _, s := range f.SubFields {
				df.Subfields = append(df.Subfields, subfield{Code: s.Code, Value: s.Value})
			}
			x.DataFields = append(x.DataFields, df)
		}
	}

	indent := ""
	if params.debug {
		indent = " "
	}
	b, err := xml.MarshalIndent(x, indent, indent)
	return string(b), err
}

func printError(r marc.Record, errType string, err error) {
	str := "== RECORD WITH ERROR STARTS HERE\n"
	str += fmt.Sprintf("%s:\n%s\n", errType, err.Error())
	str += r.DebugString() + "\n"
	str += "== RECORD WITH ERROR ENDS HERE\n\n"
	fmt.Print(str)
}

// RecordAsMARC converts a Record into a MARC record byte array
func (rec Record) RecordAsXML() (marc []byte, err error) {

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
