package gomarc21

import (
	"encoding/xml"
)

const (
	SPACE               = 0x20
	END_OF_RECORD       = 0x1D
	END_OF_FIELD        = 0x1E
	SUBFIELD_INDICATOR  = 0x1F
	DIRECTORY_ENTRY_LEN = 12
	LEADER_LEN          = 24
	MAX_RECORD_LEN      = 99999
)

type Tag string

// A MARC Record Set is for containing zero or more MARC records
type Collection struct {
	Name    xml.Name  `xml:"collection"`
	Records []*Record `xml:"record"`
}

// Controlfield contains a controlfield entry
type ControlField struct {
	Tag  Tag    `xml:"tag,attr"`
	Data string `xml:",chardata"`
}

// SubField contains a Code and a Value.
// For example in:
//		=650  \0$aDiabetes$xComplications$zUnited States.
// an example of SubFieldValue will be:
// 		SubField{
//			Code: "a",
//			Value: "Diabetes"
//		}
type SubField struct {
	Code string `xml:"code,attr"`
	Data string `xml:",chardata"`
}

// Datafield contains a datafield entry
type DataField struct {
	Tag        Tag        `xml:"tag,attr"`
	Indicator1 string     `xml:"ind1,attr"`
	Indicator2 string     `xml:"ind2,attr"`
	SubFields  []SubField `xml:"subfield"`
}

type XmlRecord struct {
	XMLName       xml.Name       `xml:"record"`
	Leader        string         `xml:"leader"`
	ControlFields []ControlField `xml:"controlfield"`
	DataFields    []DataField    `xml:"datafield"`
}

type JsonRecord struct {
	Leader        string         `json:"leader"`
	ControlFields []ControlField `json:"controlfield"`
	DataFields    []DataField    `json:"datafield"`
}

// Field represents a field inside a MARC record. Notice that the
// field could be a "control" field (tag 001-009) or a "data" field
// (any other tag)
//
// For example in:
//		=650  \0$aDiabetes$xComplications$zUnited States.
// Field would be:
// 		Field{
//			Tag: "650",
//			Value: ""
//			Indicator1: " ",
//			Indicator2: "0",
//			SubFields (see SubField definition above)
//	}

// Indicator regex
// <xsd:pattern value="[\da-z ]{1}"/>

// Field defines an interface that is satisfied by the Control and Data field
// types.
type Field interface {
	String() string
	GetTag() string
	Contains(string) bool
	AsJson() string
}

const (
	Bibliography   = iota // a bibliography (or bib-holding) record
	Holdings              // a holdings record
	Authority             // an authority record
	Classification        // a classification record
	Community             // a community Information record
	FormatUnknown         // unknown format
)

/*
var marcFormatName = map[int]string{
	Bibliography:   "Bibliography",
	Holdings:       "Holdings",
	Authority:      "Authority",
	Classification: "Classification",
	Community:      "Community Information",
	FormatUnknown:  "Unknown",
}
*/
// Leader pattern regex
//<xsd:pattern value="[\d ]{5}[\dA-Za-z ]{1}[\dA-Za-z]{1}[\dA-Za-z ]{3}(2| )(2| )[\d ]{5}[\dA-Za-z ]{3}(4500| )"/>

// Leader represents the leader of the MARC record.
// a sample authority record
//=LDR  00350cz  a2200157n  4500
//=001  cash10000\
//=003  CaOONL
//=005  20200519142027.8
//=008  850515\neanknnbabn\\\\\\\\\\\n\ana\\\\\\
//=040  \\$aCaOOP$beng$cCaOOP$dCaOONL
//=016  \\$a0200H4979
//=035  \\$a(CaOONL)182982
//=042  \\$anlc
//=151  \\$aBaffin Bay
//=667  \\$aCSH3.
//=751  \6$aBaffin, Baie de
type Leader struct {
	raw                                     []byte `json:"leader"`
	RecordLength                            int    `json:"-"` // 00 - 04 [\d ]{5} `json:"name"`
	RecordStatus                            byte   // 05 byte position [\dA-Za-z ]
	TypeOfRecord                            byte   // 06 [\dA-Za-z]
	BibLevel                                byte   // 07
	TypeOfControl                           byte   // 08
	CharCodingScheme                        byte   // 09
	IndicatorCount                          byte   // 10 - 2
	SubfieldCodeCount                       byte   // 11 - 2
	BaseAddressOfData                       int    // 12-16
	EncodingLevel                           byte   // 17
	DescrCatForm                            byte   // 18
	MultipartLevel                          byte   // 19
	LengthOfFieldPortion                    byte   // 20 - 4
	LengthOfTheStartingCharPositionPortion  byte   // 21 - 5
	LengthOfTheImplementationDefinedPortion byte   // 22 - 0
	Undefined                               byte   // 23 - 0
}

/*
source: http://www.loc.gov/marc/bibliographic/bdintro.html

    Directory - A series of entries that contain the tag, length, and
    starting location of each variable field within a record. Each
    entry is 12 character positions in length. Directory entries for
    variable control fields appear first, sequenced by the field tag in
    increasing numerical order. Entries for variable data fields
    follow, arranged in ascending order according to the first
    character of the tag. The stored sequence of the variable data
    fields in a record does not necessarily correspond to the order of
    the corresponding Directory entries. Duplicate tags are
    distinguished only by the location of the respective fields within
    the record. The Directory ends with a field terminator character
    (ASCII 1E hex).

Also:
    http://www.loc.gov/marc/holdings/hdintro.html
    http://www.loc.gov/marc/authority/adintro.html
    http://www.loc.gov/marc/classification/cdintro.html
    http://www.loc.gov/marc/community/ciintro.html

=========================================================================

source: http://www.loc.gov/marc/bibliographic/bddirectory.html

    CHARACTER POSITIONS

    00-02 - Tag
        Three ASCII numeric or ASCII alphabetic characters (upper case
        or lower case, but not both) that identify an associated
        variable field.

    03-06 - Field length
        Four ASCII numeric characters that specify the length of the
        variable field, including indicators, subfield codes, data, and
        the field terminator. A Field length number of less than four
        digits is right justified and unused positions contain zeros.

    07-11 - Starting character position
        Five ASCII numeric characters that specify the starting
        character position of the variable field relative to the Base
        address of data (Leader/12-16) of the record. A Starting
        character position number of less than five digits is right
        justified and unused positions contain zeros.
-------------------------------------------------------------------------*/
type DirectoryEntry struct {
	raw              []byte
	Tag              Tag
	FieldLength      int
	StartingPosition int
}

type Record struct {
	raw           []byte
	Leader        Leader
	Entries       []DirectoryEntry
	ControlFields []ControlField
	DataFields    []DataField
}
