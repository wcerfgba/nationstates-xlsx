package output_spec

import (
	"nationstates-xlsx/input_spec"
)

type OutputData map[string]SheetData
type SheetData map[Cell]string
type Cell struct {
	row,
	col int
}

type OutputSpec interface {
	Parse(in input_spec.InputData) (out OutputData)
	Create(data OutputData, filename string) error
	Append(data OutputData, filename string) error
}
