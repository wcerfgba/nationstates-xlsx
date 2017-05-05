package output_spec

import (
	"github.com/wcerfgba/nationstates-xlsx/input_spec"
)

type OutputData map[string]SheetData
type SheetData map[Cell]string
type Cell struct {
	row,
	col int
}

type OutputSpec interface {
	Parse(in input_spec.InputData) (out OutputData)
	Write(data OutputData, filename string) (err error)
	Create(data OutputData, filename string) (err error)
	Append(data OutputData, filename string) (err error)
}
