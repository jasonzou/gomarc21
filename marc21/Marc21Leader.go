package gomarc21

import (
	"errors"
	"fmt"
	"strconv"
)

/* --------------------------------------------------------------------

source: https://www.loc.gov/marc/specifications/specrecstruc.html

    The leader is the first field in the record and has a fixed length
    of 24 octets (character positions 0-23). Only ASCII graphic
    characters are allowed in the Leader.
=======================================================================

source: http://www.loc.gov/marc/bibliographic/bdintro.html

    Leader - Data elements that primarily provide information for the
    processing of the record. The data elements contain numbers or
    coded values and are identified by relative character position. The
    Leader is fixed in length at 24 character positions and is the
    first field of a MARC record.

Also:
    http://www.loc.gov/marc/holdings/hdleader.html
    http://www.loc.gov/marc/authority/adleader.html
    http://www.loc.gov/marc/classification/cdleader.html
    http://www.loc.gov/marc/community/cileader.html

=========================================================================*/

const (
	Bibliography   = iota // a bibliography (or bib-holding) record
	Holdings              // a holdings record
	Authority             // an authority record
	Classification        // a classification record
	Community             // a community Information record
	FormatUnknown         // unknown format
)

var marcFormatName = map[int]string{
	Bibliography:   "Bibliography",
	Holdings:       "Holdings",
	Authority:      "Authority",
	Classification: "Classification",
	Community:      "Community Information",
	FormatUnknown:  "Unknown",
}

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
	raw                                     []byte
	RecordLength                            int  // 00 - 04 [\d ]{5}
	RecordStatus                            byte // 05 byte position [\dA-Za-z ]
	TypeOfRecord                            byte // 06 [\dA-Za-z]
	BibLevel                                byte // 07
	TypeOfControl                           byte // 08
	CharCodingScheme                        byte // 09
	IndicatorCount                          byte // 10 - 2
	SubfieldCodeCount                       byte // 11 - 2
	BaseAddressOfData                       int  // 12-16
	EncodingLevel                           byte // 17
	DescrCatForm                            byte // 18
	MultipartLevel                          byte // 19
	LengthOfFieldPortion                    byte // 20 - 4
	LengthOfTheStartingCharPositionPortion  byte // 21 - 5
	LengthOfTheImplementationDefinedPortion byte // 22 - 0
	Undefined                               byte // 23 - 0
}

// NewLeader creates a Leader from the data in the MARC record.
func NewLeader(bytes []byte) (error, Leader) {
	if len(bytes) != 24 {
		return errors.New("Incomplete leader"), Leader{}
	}

	recordLen, err := strconv.Atoi(string(bytes[0:5]))

	if err != nil {
		msg := fmt.Sprintf("Could not determine the length of the record from leader (%s)", string(bytes))
		err = errors.New(msg)
		recordLen = -1
	}

	offset, err := strconv.Atoi(string(bytes[12:17]))
	if err != nil {
		msg := fmt.Sprintf("Could not determine data offset from leader (%s)", string(bytes))
		err = errors.New(msg)
		offset = -1
	}

	leader := Leader{
		raw:                                     bytes[:24],
		RecordLength:                            recordLen,
		RecordStatus:                            bytes[5],
		TypeOfRecord:                            bytes[6],
		BibLevel:                                bytes[7],
		TypeOfControl:                           bytes[8],
		CharCodingScheme:                        bytes[9],
		IndicatorCount:                          bytes[10],
		SubfieldCodeCount:                       bytes[11],
		BaseAddressOfData:                       offset,
		EncodingLevel:                           bytes[17],
		DescrCatForm:                            bytes[18],
		MultipartLevel:                          bytes[19],
		LengthOfFieldPortion:                    bytes[20],
		LengthOfTheStartingCharPositionPortion:  bytes[21],
		LengthOfTheImplementationDefinedPortion: bytes[22],
		Undefined:                               bytes[23],
	}
	return err, leader
}

func (l Leader) String() string {
	return fmt.Sprintf("=LDR  %s", string(l.raw))
}

func (l Leader) Raw() string {
	return fmt.Sprintf("%s", string(l.raw))
}

func (l Leader) GetRecordLength() int {
	// record length
	return l.RecordLength
}

func (l Leader) GetBaseAddressOfData() int {
	// record base address of data of a record
	return l.BaseAddressOfData
}

// RecordFormat indicates the high level nature of the record and is
// used to differentiate between Bibliography, Holdings, Authority,
// Classification, and Community record formats.
func (rec Record) RecordFormat() int {
	code := pluckByte(rec.Leader.Text, 6) //?????
	switch code {
	case "q":
		return Community
	case "z":
		return Authority
	case "w":
		return Classification
	case "u", "v", "x", "y":
		return Holdings
	case "a", "c", "d", "e", "f", "g", "i", "j", "k", "m", "o", "p", "r", "t":
		return Bibliography
	}
	return FmtUnknown
}

// RecordFormatName indicates the name of the format of the record and is
// used to differentiate between Bibliography, Holdings, Authority,
// Classification, and Community record formats.
func (rec Record) RecordFormatName() string {
	f := rec.RecordFormat()
	n := marcFormatName[f]
	return n
}

//  06 - Type of record
var recordType = map[string]string{
	// Bibliography
	"a": "Language material",
	"c": "Notated music",
	"d": "Manuscript notated music",
	"e": "Cartographic material",
	"f": "Manuscript cartographic material",
	"g": "Projected medium",
	"i": "Nonmusical sound recording",
	"j": "Musical sound recording",
	"k": "Two-dimensional nonprojectable graphic",
	"m": "Computer file",
	"o": "Kit",
	"p": "Mixed materials",
	"r": "Three-dimensional artifact or naturally occurring object",
	"t": "Manuscript language material",
	// Holding
	"u": "Unknown",
	"v": "Multipart item holdings",
	"x": "Single-part item holdings",
	"y": "Serial item holdings",
	// Classification
	"w": "Classification data",
	// Authority
	"z": "Authority data",
	// Community
	"q": "Community information",
}

// RecordType returns the one character code and label indicating
// the "06 - Type of record" for the record. Use RecordFormat
// to determine the record format (bibliographic, holdings, etc.)
func (rec Record) RecordType() (code, label string) {
	return shortCodeLookup(recordType, rec.Leader.Text, 6)
}

//  09 - Character coding scheme
var characterCodingScheme = map[string]string{
	" ": "MARC-8",
	"a": "UCS/Unicode",
}

// CharacterCodingScheme returns the code and label indicating the
// "09 - Character coding scheme" of the record (MARC-8 or UCS/Unicode).
func (rec Record) CharacterCodingScheme() (code, label string) {
	return shortCodeLookup(characterCodingScheme, rec.Leader.Text, 9)
}

////////////////////////////////////////////////////////////////////////
// Lookup data for Bibliography records

var bibliographyMaterialType = map[string]string{
	"BK": "Books",
	"CF": "Computer Files",
	"MP": "Maps",
	"MU": "Music",
	"CR": "Continuing Resources",
	"VM": "Visual Materials",
	"MX": "Mixed Materials",
}

// BibliographyMaterialType returns the code and description of the type
// of material documented by Bibliography the record.
// {"Books",  "Computer Files", "Maps", "Music", "Continuing Resources",
// "Visual Materials" or "Mixed Materials"}
func (rec Record) BibliographyMaterialType() (code, label string) {
	rectype := pluckByte(rec.Leader.Text, 6)

	// the simple record type to material type mappings
	var rtmap = map[string]string{
		"c": "MU",
		"d": "MU",
		"i": "MU",
		"j": "MU",
		"e": "MP",
		"f": "MP",
		"g": "VM",
		"k": "VM",
		"o": "VM",
		"r": "VM",
		"m": "CF",
		"p": "MX",
	}

	code, ok := rtmap[rectype]
	if !ok {
		// no simple match
		biblevel := pluckByte(rec.Leader.Text, 7)
		switch biblevel {
		case "a", "c", "d", "m":
			if rectype == "a" || rectype == "t" {
				code = "BK"
			}
		case "b", "i", "s":
			if rectype == "a" {
				code = "CR"
			}
		}
	}

	label = bibliographyMaterialType[code]
	return code, label
}

////////////////////////////////////////////////////////////////////////

// GetText returns the text for the leader
func (ldr Leader) GetText() string {
	return ldr.Text
}
