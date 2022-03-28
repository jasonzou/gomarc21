package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jasonzou/gomarc21"
)

func main() {

	var recsPerFile int
	var marcfile string
	var dir string

	flag.IntVar(&recsPerFile, "c", 1000, "The number of MARC records per output file (defaults to 1000).")
	flag.StringVar(&marcfile, "m", "", "The file that contains the MARC records.")
	flag.StringVar(&dir, "d", "mark_split", "The directory to write the output files to (defaults to mark_split).")
	flag.Parse()

	fi, err := os.Open(marcfile)
	if err != nil {
		log.Fatal(fmt.Printf("File open failed: %q", err))
	}
	defer func() {
		if cerr := fi.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	for {
		rec, err := gomarc21.ParseNextRecord(fi)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		recxml, err := rec.AsXML()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(recxml)
	}

}
