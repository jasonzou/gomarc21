package gomarc21

import (
	"fmt"
)

// SubField contains a Code and a Value.
// For example in:
//		=650  \0$aDiabetes$xComplications$zUnited States.
// an example of SubFieldValue will be:
// 		SubField{
//			Code: "a",
//			Value: "Diabetes"
//		}
type SubField struct {
	Code  string
	Value string
}

// SubField Code regex
//<xsd:pattern value="[\dA-Za-z!"#$%&'()*+,-./:;<=>?{}_^`~\[\]\\]{1}"/>

// String returns the subfield as a string.
func (sf SubField) String() string {
	return fmt.Sprintf("$%s%s", sf.Code, sf.Value)
}

// GetCode returns the code for the subfield
func (sf SubField) GetCode() string {
	return sf.Code
}

// GetText returns the text for the subfield
func (sf SubField) GetValue() string {
	return sf.Value
}

func (sf SubField) AddSubField(Code string, Value string) {
}

func (sf SubField) RemoveAllSubFields(Code string) {
}

func (sf SubField) ReplaceSubField(Code string, Value string) {
}
