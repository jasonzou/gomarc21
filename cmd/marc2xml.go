package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jasonzou/gomarc21"
)

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

	//xml header?
	for {
		rec, err := gomarc21.ParseNextRecord(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		recxml, err := rec.RecordAsXml()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(recxml)
		fmt.Println()
	}
}

func showHelp() {
	fmt.Println(os.Args[0])
	fmt.Println("   Converts a MARC file to MARCXML.")
	fmt.Printf("    Usage: %s <MARC file to convert>\n", os.Args[0])
	fmt.Println()
	os.Exit(0)
}
