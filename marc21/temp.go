
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