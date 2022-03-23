package gomarc21

import (
	"errors"
	"fmt"
	"strconv"
)

/* ----------------------------------------------------------------------

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
	Raw              []byte
	Tag              Tag
	FieldLength      int
	StartingPosition int
}

func (dirEntry DirectoryEntry) GetRaw() string {
	// need to make sure 12 bytes
	return string(dirEntry.Raw)
}
func (dirEntry DirectoryEntry) GetFieldLength() int {
	return dirEntry.FieldLength
}
func (dirEntry DirectoryEntry) GetStartingPosition() int {
	return dirEntry.StartingPosition
}
func (dirEntry DirectoryEntry) GetTag() string {
	return dirEntry.Tag.GetTag()
}
func NewDirectoryEntry(entryBytes []byte) (DirectoryEntry, error) {
	if len(entryBytes) != 12 {
		msg := fmt.Sprintf("The length of directry entry is not 12, but %d instead", len(entryBytes))
		err := errors.New(msg)
		return DirectoryEntry{}, err
	}

	tag, err := NewTag(entryBytes[0:3])
	if err != nil {
		return DirectoryEntry{}, errors.New("tag is invalid")
	}

	fieldLength, err := strconv.Atoi(string(entryBytes[3:7]))
	if err != nil {
		return DirectoryEntry{}, errors.New("fieldLength is invalid")
	}

	startingPosition, err := strconv.Atoi(string(entryBytes[7:12]))
	if err != nil {
		return DirectoryEntry{}, errors.New("starting position is invalid")
	}

	directoryEntry := DirectoryEntry{
		Raw:              entryBytes[:12],
		Tag:              tag,
		FieldLength:      fieldLength,
		StartingPosition: startingPosition,
	}

	return directoryEntry, nil
}
