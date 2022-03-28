package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kong"
)

var CLI struct {
	InputFile  string `short:"i" name:"input" help:"The file contains MARC records." type:"existingfile"`
	OutputFile string `short:"o" name:"output" help:"The file will contain Json records converted from the input MARC records." type:"file"`
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("marc2json"),
		kong.Description("Convert MARC records into Json records."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))
	fmt.Print(ctx.Command())
	fmt.Print(CLI.InputFile)
	fmt.Print(CLI.OutputFile)

	var marcfile = CLI.InputFile
	//var output = CLI.OutputFile

	fi, err := os.Open(marcfile)
	if err != nil {
		log.Fatal(fmt.Printf("File open failed: %q", err))
	}
	defer func() {
		if cerr := fi.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	fmt.Print(marc21.CollectionXMLHeader)

	for {
		rec, err := marc21.ParseNextRecord(fi)
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

	fmt.Print(marc21.CollectionXMLFooter)
}

func showHelp() {
	fmt.Println(os.Args[0])
	fmt.Println("   Converts a MARC file to MARCXML.")
	fmt.Printf("    Usage: %s <MARC file to convert>\n", os.Args[0])
	fmt.Println()
	os.Exit(0)
}
