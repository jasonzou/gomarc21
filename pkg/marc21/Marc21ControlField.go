package gomarc21

import (
	"fmt"
	"strings"
)

/*
https://www.loc.gov/marc/specifications/specrecstruc.html

    Control fields in MARC 21 formats are assigned tags beginning with
    two zeroes. They are comprised of data and a field terminator; they
    do not contain indicators or subfield codes. The control number
    field is assigned tag 001 and contains the control number of the
    record. Each record contains only one control number field (with
    tag 001), which is to be located at the base address of data.

http://www.loc.gov/marc/bibliographic/bdintro.html

    Variable control fields - The 00X fields. These fields are
    identified by a field tag in the Directory but they contain neither
    indicator positions nor subfield codes. The variable control fields
    are structurally different from the variable data fields. They may
    contain either a single data element or a series of fixed-length
    data elements identified by relative character position.

 http://www.loc.gov/marc/bibliographic/bd00x.html
 http://www.loc.gov/marc/holdings/hd00x.html
 http://www.loc.gov/marc/authority/ad00x.html
 http://www.loc.gov/marc/classification/cd00x.html
 http://www.loc.gov/marc/community/ci00x.html
*/

// Implement the Stringer interface for "Pretty-printing"
func (cf ControlField) String() string {
	return fmt.Sprintf("=%s  %s", cf.Tag.GetTag(), cf.Data)
}

// GetTag returns the tag for the controlfield
func (cf ControlField) GetTag() Tag {
	return cf.Tag
}

// GetTag returns the tag for the controlfield
func (cf ControlField) Contains(str string) bool {
	str = strings.ToLower(str)
	isControlField, err := cf.Tag.IsControlTag()
	if err == nil {
		if isControlField {
			return strings.Contains(strings.ToLower(cf.Data), str)
		}
	}
	return false
}

// GetDatareturns the text for the ControlField
func (cf ControlField) GetData() string {
	return cf.Data
}
