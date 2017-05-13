package output_spec

import (
	"github.com/wcerfgba/nationstates-xlsx/input_spec"
)

type OutputData map[string]SheetData
type SheetData map[CellAddress]CellData
type CellAddress struct {
	Row,
	Col int
}
type CellData struct {
	Contents          string
	NotEmptyBehaviour NotEmptyBehaviour
}

type NotEmptyBehaviour int

const (
	Skip NotEmptyBehaviour = iota
	Replace
	StopIfNotEqual
	IncrementColUntilEmpty
	IncrementRowUntilEmpty
	DecrementColUntilEmptyOrSkip
	DecrementRowUntilEmptyOrSkip
	DecrementColUntilEmptyOrReplace
	DecrementRowUntilEmptyOrReplace
)

type OutputSpec interface {
	Parse(in input_spec.InputData) (out OutputData, err error)
	Write(data OutputData, filename string) (action string, err error)
	Create(data OutputData, filename string) (err error)
	Append(data OutputData, filename string) (err error)
}
