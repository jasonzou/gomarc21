# gomarc21

[![GoDoc](https://godoc.org/github.com/jasonzou/gomarc21/pkg/marc21?status.svg)](https://godoc.org/github.com/jasonzou/gomarc21/pkg/marc21)
[![Go Report Card](https://goreportcard.com/badge/github.com/jasonzou/gomarc21)](https://goreportcard.com/report/github.com/jasonzou/gomarc21)

Another golang implementation of MARC 21.

## Functionality

- parse marc21 files
- convert marc21 into mrk format
- write MARCMaker files https://www.loc.gov/marc/makrbrkr.html
- Read MARC21 and MARCXML data.

## A to-do list
- Read MARCXML data.
- convert marc21 into marc21 xml format
- convert marc21 into marc21 json format
- validate marc21 json against avarm marc21 schema
- authority records 
- more tests
- Convert MARC-8 encoding to UTF-8
- Perform error checking on MARC records
- Read MARCMaker files https://www.loc.gov/marc/makrbrkr.html

## Revision History

- March 23, 2022 version 0.0.3, add files into cmd folder
Marc21Record.go is working...
- March 23, 2022 version 0.0.2, Marc21Record.go is working...
- March 21, 2022 Initial version (0.0.1) created based marcnet, other marc21 golang implementations
