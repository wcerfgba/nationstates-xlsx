package main

import (
	"flag"
	"log"

	"github.com/wcerfgba/nationstates-xlsx/data"
	"github.com/wcerfgba/nationstates-xlsx/util"
)

var (
	provider  data.Provider  = &data.Provider20170429{}
	outputter data.Outputter = &data.XlsxAppendingOutputter{}

	nation      string
	outFileName string
	debug       bool
)

func init() {
	flag.StringVar(&nation, "n", "", "Name of nation to request data for.")
	flag.StringVar(&outFileName, "o", "", "Output file.")
	flag.BoolVar(&debug, "debug", false, "Print debugging information.")
}

func main() {
	flag.Parse()

	provider.Configure(util.Configuration{
		"nation": nation,
	})

	outputter.Configure(util.Configuration{
		"outputFilename": outFileName,
	})

	res := provider.Get()
	err := outputter.Output(res)

	if err != nil {
		log.Println(err)
	}

	// log.Println("Parsing response")
	// input, extraInput, err := inputSpec.Parse(raw)
	// if extraInput != nil && len(extraInput) > 0 {
	// 	log.Println("Unimplemented fields:", extraInput)
	// }
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if debug {
	// 	log.Printf("Input: %v", input)
	// }

	// log.Println("Mapping")
	// output, extraOutput, err := outputSpec.Parse(input)
	// if extraOutput != nil {
	// 	log.Printf("Unused fields: %v", extraOutput)
	// }
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if debug {
	// 	log.Printf("Output: %v", output)
	// }

	// log.Printf("Writing %v", outFileName)
	// action, err := outputSpec.Write(output, outFileName)
	// log.Printf("Mode used: %v", action)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Done")
}
