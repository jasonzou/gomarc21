package gomarc21

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func openTestMARC(t *testing.T) (data *os.File) {
	data, err := os.Open("data/test_1a.mrc")
	if err != nil {
		t.Fatal(err)
	}
	return
}

func TestGetRaw1(test *testing.T) {
}

// ControlNum returns the control number (tag 001) for the record.
func TestControlNum(test *testing.T) {
}

// ControlNum returns the control number (tag 001) for the record.
func TestTitle(test *testing.T) {
}

// parse directory of a marc record
func TestParseDirectory(test *testing.T) {

}

// parse leader of a marc record
func TestParseLeader(test *testing.T) {
}

// parse control fields from the raw MARC record bytes
func TestParseControlFields(test *testing.T) {

}

// extractDatafields extracts the data fields/sub-fields from the raw MARC record bytes
func TestParseDataFields(test *testing.T) {
}

// GetDatafields returns datafields for the record that match the
// specified comma separated list of tags. If no tags are specified
// (empty string) then all datafields are returned
func TestGetDatafields(test *testing.T) {
}

func TestGetMrk(test *testing.T) {
	const mrk = `=LDR  01805nam a2200385 i 4500
=001  ocm57175940
=005  20041206161421.0
=006  m        d f
=007  cr cn-
=008  041206s1976    dcua    sb   f000 0 eng c
=040  \\$aGPO$cGPO$dMvI$dMvI
=042  \\$apcc
=043  \\$an-us---
=074  \\$a0620-A (online)
=086  0\$aI 19.4/2:735
=100  1\$aSwanson, Vernon E.$q(Vernon Emmanuel),$d1922-1992.
=245  10$aGuidelines for sample collecting and analytical methods used in the U.S. Geological Survey for determining chemical composition of coal$h[electronic resource] /$cby Vernon E. Swanson and Claude Huffman, Jr.
=260  \\$a[Washington, D.C.] :$bU.S. Dept. of the Interior, U.S. Geological Survey,$c1976.
=336  \\$atext$2rdacontent.
=337  \\$acomputer$2rdamedia.
=338  \\$aonline resource$2rdacarrier.
=440  \0$aGeological Survey circular ;$v735.
=500  \\$aTitle from title screen (viewed on Dec. 06, 2004)
=504  \\$aIncludes bibliographical references.
=538  \\$aMode of access: Internet from the USGS Web site. Address as of 12/06/04: http://pubs.usgs.gov/circ/c735/index.htm; current access is available via PURL.
=650  \0$aCoal$xAnalysis.
=650  \0$aCoal$xSampling.
=700  1\$aHuffman, Claude.
=776  1\$aSwanson, Vernon Emanuel,$d1922-$tGuidelines for sample collecting and analytical methods used in the U.S. Geological Survey for determining chemical composition of coal$hiv, 11 p.$w(OCoLC)2331861.
=856  40$uhttp://purl.access.gpo.gov/GPO/LPS56007$zView online version
=907  \\$a.b37991760$b04-08-17$c07-26-05
=998  \\$aes001$b07-26-05$cm$da$e-$feng$gdcu$h0$i1
=910  \\$aMARCIVE
=910  \\$aHathi Trust report None
=945  \\$g0$j0$lesb  $on$p$0.00$q $r $s-$t255$u0$v0$w0$x0$y.i138993579$z07-26-05`
	data := openTestMARC(test)
	defer data.Close()

	rec, err := ReadRecord(data)
	if err != nil {
		fmt.Println(err)
		test.Error("failed to read a marc record")
	}
	if rec.GetMrk() == mrk {
		log.Print(rec.GetMrk())
		test.Error("not complete match")
	}
}
func TestParseRecord(test *testing.T) {
	const exp = `=LDR  00350cz  a2200157n  4500
	=001  cash10000\
	=003  CaOONL
	=005  20200519142027.8
	=008  850515\neanknnbabn\\\\\\\\\\\n\ana\\\\\\
	=040  \\$aCaOOP$beng$cCaOOP$dCaOONL
	=016  \\$a0200H4979
	=035  \\$a(CaOONL)182982
	=042  \\$anlc
	=151  \\$aBaffin Bay
	=667  \\$aCSH3.
	=751  \6$aBaffin, Baie de`
	data := openTestMARC(test)
	defer data.Close()

	rec, err := ReadRecord(data)
	if err != nil {
		fmt.Println(err)
		test.Error("failed to read a marc record")
	}
	if rec.Leader.BaseAddressOfData != 385 {
		test.Error("the base address of data is 332, but is ", rec.Leader.BaseAddressOfData)
	}

}

func TestNextRecord(test *testing.T) {

}

func TestParseNextRecord(test *testing.T) {

}

func TestRecordAsMARC(test *testing.T) {

}
