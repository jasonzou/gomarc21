//package gomarc21
package main

import (
    "reflect"
    "errors"
  "fmt"
)

type Tag struct {
	Tag1  byte  // for both Control and Data fields
	Tag2  byte  // for Data fields
	Tag3  byte  // for Data fields
}

func (tag *Tag) GetTag() string {
  myTag := fmt.Sprintf("%c%c%c", tag.Tag1, tag.Tag2, tag.Tag3)
  fmt.Println(myTag)
  fmt.Println(reflect.TypeOf(myTag))
  return myTag
}

func (tag *Tag) IsControlTag() (error, bool) {
  // <xsd:pattern value="00[1-9A-Za-z]{1}"/>
  // https://www.loc.gov/standards/marcxml/schema/MARC21slim.xsd
  if (tag.Tag1 == '0' && tag.Tag2 == '0'){
    if (tag.Tag3 >= '1' && tag.Tag3 <= 'z'){
      return nil, true
    }else{
      return errors.New("Invalid tag"), false
    }
  }
  return nil, false
}

func (tag *Tag) IsDataTag() (error, bool) {
  // <xsd:pattern value="(0([1-9A-Z][0-9A-Z])|0([1-9a-z][0-9a-z]))|(([1-9A-Z][0-9A-Z]{2})|([1-9a-z][0-9a-z]{2}))"/>
  // https://www.loc.gov/standards/marcxml/schema/MARC21slim.xsd
  
  // 0([1-9A-Z][0-9A-Z]
  // 0([1-9a-z][0-9a-z])
  if (tag.Tag1 >= '0' && tag.Tag1 <= 'z'){ 
    if (tag.Tag2 >= '1' && tag.Tag2 <= 'z'){ 
      if (tag.Tag3 >= '0' && tag.Tag3 <= 'z'){
        return nil, true
      }
    }
  }

  // [1-9A-Z][0-9A-Z]{2}
  // [1-9a-z][0-9a-z]{2}
  if (tag.Tag1 >= '1' && tag.Tag1 <= 'z'){ 
    if (tag.Tag2 >= '0' && tag.Tag2 <= 'z'){ 
      if (tag.Tag3 >= '0' && tag.Tag3 <= 'z'){
        return nil, true
      }
    }
  }

  err, controlTag := tag.IsControlTag()
  if err == nil {
    if controlTag == false {
      return errors.New("Invalid data tag"), false
    }else{
      return nil, false
    }
  }
      
  return errors.New("Invalid data tag"), false
}

func main(){
  var t Tag
  t.Tag1 = '1'
  t.Tag2 = '2'
  t.Tag3 = '3'

  fmt.Println(t.GetTag())
  err, test := t.IsControlTag()
  if err == nil {
    fmt.Println(test)
  }
  err, test = t.IsDataTag()
  if err == nil {
    fmt.Println(test, "data tag")
  }
  t.Tag1 = '0'
  t.Tag2 = '0'
  t.Tag3 = 'z'
  err, test = t.IsControlTag()
  if err == nil {
    fmt.Println(test)
  }
  t.Tag2 = '#'
  err, test = t.IsDataTag()
  fmt.Println(t.GetTag())
  if err == nil {
    fmt.Println(test, "data tag")
  }else{
    fmt.Println(err)
  }
}
