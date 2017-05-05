package main

import (
	"flag"
	"io/ioutil"
	"log"
	"nationstates-xlsx/input_spec"
	"nationstates-xlsx/output_spec"
	"net/http"
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
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Got response: %v", raw)

	log.Println("Parsing response")
	input, err := inputSpec.Parse(raw)
	//if extra != nil {
	//	log.Println("Unimplemented fields:", extra)
	//}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Input: %v", input)

	log.Println("Mapping")
	output := outputSpec.Parse(input)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Output: %v", output)

	log.Printf("Creating %v", outFileName)
	err = outputSpec.Create(output, outFileName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
