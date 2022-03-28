package gomarc21

/*
func main() {

	var marcfile string
	if len(os.Args) > 1 {
		marcfile = os.Args[1]
	}

	if marcfile == "" {
		showHelp()
	}

	data, err := os.Open(marcfile)
	if err != nil {
		data.Close()
		fmt.Errorf("Failed to open the marc file: %s", marcfile)
		return
	}
	defer data.Close()

	for {
		rec, err := gomarc21.ParseNextRecord(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		recxml, err := rec.RecordAsJson()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(recxml)
	}
}
func main() {

	var recsPerFile int
	var marcFile string
	var dir string

	flag.IntVar(&recsPerFile, "c", 1000, "The number of MARC records per output file (defaults to 1000).")
	flag.StringVar(&marcFile, "m", "", "The file that contains the MARC records.")
	flag.StringVar(&dir, "d", "mark_split", "The directory to write the output files to (defaults to mark_split).")
	flag.Parse()

	fi, err := os.Open(marcFile)
	if err != nil {
		log.Fatal(fmt.Printf("File open failed: %q", err))
	}
	defer func() {
		if cerr := fi.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	recCount := 0
	fOut, fileCount := nextFile(dir, 0)

	for {
		rawRec, err := marc21.NextRecord(fi)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if _, err := fOut.Write(rawRec); err != nil {
			log.Fatal(err)
		}

		recCount++
		if recCount >= recsPerFile {
			closeFile(fOut)
			fOut, fileCount = nextFile(dir, fileCount)
			recCount = 0
		}
	}

	closeFile(fOut)
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatal(err)
	}
}


func nextFile(dir string, i int) (*os.File, int) {
	i++
	fileName := fmt.Sprintf("%s/%06d.mrc", dir, i)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		log.Fatal(fmt.Printf("File open failed: %q", err))
	}
	return f, i
}
*/
