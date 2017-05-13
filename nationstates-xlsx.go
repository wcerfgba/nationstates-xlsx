package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wcerfgba/nationstates-xlsx/input_spec"
	"github.com/wcerfgba/nationstates-xlsx/output_spec"
)

var (
	inputSpec  input_spec.InputSpec   = &input_spec.Nation_20170429{}
	outputSpec output_spec.OutputSpec = &output_spec.T20170429{}

	nation      string
	outFileName string
)

func init() {
	flag.StringVar(&nation, "n", "", "Name of nation to request data for.")
	flag.StringVar(&outFileName, "o", "", "Output file.")
}

func main() {
	flag.Parse()

	log.Printf("Fetching %v", nation)
	url := inputSpec.BuildRequestUrl(nation)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Reading response")
	raw, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Parsing response")
	input, extraInput, err := inputSpec.Parse(raw)
	if extraInput != nil && len(extraInput) > 0 {
		log.Println("Unimplemented fields:", extraInput)
	}
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("Input: %v", input)

	log.Println("Mapping")
	output, extraOutput, err := outputSpec.Parse(input)
	if extraOutput != nil {
		log.Printf("Unused fields: %v", extraOutput)
	}
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("Output: %v", output)

	log.Printf("Writing %v", outFileName)
	action, err := outputSpec.Write(output, outFileName)
	log.Printf("Mode used: %v", action)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
