package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/jasonzou/gomarc21"
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

		recxml, err := rec.RecordAsMrk()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(recxml)
		fmt.Println()
	}
}
