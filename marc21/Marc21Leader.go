package gomarc21

import (
	"errors"
	"fmt"
	"strconv"
)

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
	raw           []byte
	RecordLength  int  // 00 - 04 [\d ]{5}
	RecordStatus  byte // 05 byte position [\dA-Za-z ]
  TypeOfRecord	byte // 06 [\dA-Za-z]
	BibLevel      byte // 07 
	TypeOfControl byte // 08
  CharCodingScheme byte // 09
  IndicatorCount byte // 10 - 2
  SubfieldCodeCount byte // 11 - 2
  BaseAddressOfData int  // 12-16
	EncodingLevel byte // 17
	DescrCatForm  byte // 18
	MultipartLevel     byte // 19
  LengthOfFieldPortion byte // 20 - 4
  LengthOfTheStartingCharPositionPortion byte // 21 - 5
  LengthOfTheImplementationDefinedPortion byte // 22 - 0
  Undefined byte // 23 - 0
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
    raw:           bytes[:24],
    RecordLength:  recordLen,
    RecordStatus : bytes[5],
    TypeOfRecord	:bytes[6],
    BibLevel:      bytes[7], 
    TypeOfControl: bytes[8],
    CharCodingScheme: bytes[9],
    IndicatorCount: bytes[10],
    SubfieldCodeCount: bytes[11],
    BaseAddressOfData: offset,
    EncodingLevel: bytes[17],
    DescrCatForm:  bytes[18],
    MultipartLevel:     bytes[19],
    LengthOfFieldPortion: bytes[20],
    LengthOfTheStartingCharPositionPortion: bytes[21],
    LengthOfTheImplementationDefinedPortion: bytes[22],
    Undefined:  bytes[23],
	}
	return err, leader
}

func (l Leader) String() string {
	return fmt.Sprintf("=LDR  %s", string(l.raw))
}

func (l Leader) Raw() string {
	return fmt.Sprintf("%s",string(l.raw))
}

func (l Leader) GetRecordLength() int {
  // record length
	return l.RecordLength
}

func (l Leader) GetBaseAddressOfData() int {
  // record base address of data of a record
	return l.BaseAddressOfData
}


