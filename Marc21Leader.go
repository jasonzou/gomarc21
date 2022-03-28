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

// NewLeader creates a Leader from the data in the MARC record.
func NewLeader(bytes []byte) (Leader, error) {
	if len(bytes) != 24 {
		return Leader{}, errors.New("incomplete leader")
	}

	recordLen, err := strconv.Atoi(string(bytes[0:5]))

	if err != nil {
		msg := fmt.Sprintf("could not determine the length of the record from leader (%s)", string(bytes))
		err = errors.New(msg)
		return Leader{}, err
	}

	offset, err := strconv.Atoi(string(bytes[12:17]))
	if err != nil {
		msg := fmt.Sprintf("could not determine data offset from leader (%s)", string(bytes))
		err = errors.New(msg)
		return Leader{}, err
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
	return leader, err
}

func (l Leader) String() string {
	return fmt.Sprintf("=LDR  %s", string(l.raw))
}

func (l Leader) GetRaw() string {
	return string(l.raw)
}

func (l Leader) GetRecordLength() int {
	// record length
	return l.RecordLength
}

func (l Leader) GetBaseAddressOfData() int {
	// record base address of data of a record
	return l.BaseAddressOfData
}
