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
- Read [MARCMaker files](https://www.loc.gov/marc/makrbrkr.html)

## Revision History
- version 0.0.8, cmd files using kong instead of flag. March 28, 2022
- version 0.0.8, cmd added. March 28,2022
- version 0.0.7, recordAsJson is working. This approach may need to be revisited later. March 27, 2022
- version 0.0.6, recordAsXml is working; reflect Tag. March 25, 2022
- version 0.0.5, more tests. March 24, 2022
- version 0.0.4, Marc21Record_test.go is working;worked on json. March 23, 2022
- version 0.0.3, add files into cmd folder Marc21Record.go is working... March 23,2022
- version 0.0.2, Marc21Record.go is working... March 23, 2022
- version 0.0.1, Initial version (0.0.1) created based marcnet, other marc21 golang implementations. March 21, 2022
