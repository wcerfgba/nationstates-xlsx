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
		"resultFactory": util.ObviousStringTreeNode{}
	})

	outputter.Configure(util.Configuration{
		"outputFilename": outFileName,
	})

	res := provider.Get()
	err := outputter.Output(res)

	if err != nil {
		log.Println(err)
	}
}
