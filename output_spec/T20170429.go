package output_spec

import (
	"math"
	"strconv"

	"github.com/wcerfgba/nationstates-xlsx/input_spec"

	"os"

	"github.com/tealeg/xlsx"

	"fmt"

	"regexp"

	"strings"

	. "github.com/ahmetb/go-linq"
)

type T20170429 struct {
}

func (s *T20170429) Parse(in input_spec.InputData) (out OutputData, extra input_spec.InputData, err error) {

	timestamp := in["_timestamp"].(string)

	// Causes of death
	causesOfDeath := buildSheet("Causes of death", simpleSheetSpec{
		"Old age",
		"Heart disease",
		"Murder",
		"Cancer",
		"Acts of God",
		"Capital Punishment",
		"Exposure",
	}, in["DEATHS"], 0, timestamp)

	// Government expenditure
	governmentExpenditure := buildSheet("Government expenditure", simpleSheetSpec{
		"Administration",
		"Defence",
		"Education",
		"Enviroment",
		"Healthcare",
		"Industry",
		"International aid",
		"Law and Order",
		"Public Transport",
		"Social Policy",
		"Spirituality",
		"Welfare",
	}, in["GOVT"], 0, timestamp)

	gdpInt, err := strconv.ParseInt(in["GDP"].(string), 10, 64)
	if err != nil {
		return
	}
	gdpBnsStr := fmt.Sprintf("%.3f", float64(gdpInt)/math.Pow10(9))
	gdpBnsFloat, err := strconv.ParseFloat(gdpBnsStr, 64)
	if err != nil {
		return
	}
	publicSectorFloat, err := strconv.ParseFloat(in["PUBLICSECTOR"].(string), 64)
	if err != nil {
		return
	}
	publicExpenditureBns := (publicSectorFloat / 100) * gdpBnsFloat
	publicExpenditureBnsStr := fmt.Sprintf("%.3f", publicExpenditureBns)

	From(SheetData{
		CellAddress{1, 1}: CellData{"Expenditure (billion)", StopIfNotEqual},
		CellAddress{1, 2}: CellData{"% of GDP", StopIfNotEqual},
		CellAddress{2, 1}: CellData{publicExpenditureBnsStr, IncrementRowUntilEmpty},
		CellAddress{2, 2}: CellData{in["PUBLICSECTOR"].(string), IncrementRowUntilEmpty},
	}).ToMap(&governmentExpenditure)

	// Economy
	economy := buildSheet("Economy", renameSheetSpec{
		"Government":           nil,
		"State-owned Industry": "PUBLIC",
		"Private Industry":     "INDUSTRY",
		"Black Market":         nil,
	}, in["SECTORS"], 2, timestamp)

	From(SheetData{
		CellAddress{1, 1}: CellData{"GDP (billion)", StopIfNotEqual},
		CellAddress{1, 2}: CellData{"Ave. wage", StopIfNotEqual},
		CellAddress{2, 1}: CellData{gdpBnsStr, IncrementRowUntilEmpty},
		CellAddress{2, 2}: CellData{in["INCOME"].(string), IncrementRowUntilEmpty},
	}).ToMap(&economy)

	// Rights
	rights := buildSheet("Rights", simpleSheetSpec{
		"Civil Rights",
		"Economy",
		"Political Freedom",
	}, in["FREEDOMSCORES"], 0, timestamp)

	// Output
	out = OutputData{
		"Causes of death":        causesOfDeath,
		"Government expenditure": governmentExpenditure,
		"Economy":                economy,
		"Rights":                 rights,
	}
	return
}

func (s *T20170429) Write(data OutputData, filename string) (err error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err = s.Create(data, filename)
	} else {
		err = s.Append(data, filename)
	}
	return
}

func (s *T20170429) Create(data OutputData, filename string) (err error) {
	f := xlsx.NewFile()
	for sheet, sheetData := range data {
		sheet, err := getSheet(f, sheet)
		if err != nil {
			return err
		}
		for cell, v := range sheetData {
			sheet.Cell(cell.Row, cell.Col).Value = v.Contents
		}
	}
	err = f.Save(filename)
	return
}

func (s *T20170429) Append(data OutputData, filename string) (err error) {
	return
}

func getSheet(f *xlsx.File, sheet string) (*xlsx.Sheet, error) {
	for _, s := range f.Sheets {
		if s.Name == sheet {
			return s, nil
		}
	}
	return f.AddSheet(sheet)
}

func buildSheet(title string, spec interface{}, data interface{}, colOffset int, timestamp string) (sheet SheetData) {
	fullSpec := map[string]interface{}{}
	switch specT := spec.(type) {
	case simpleSheetSpec:
		for _, v := range specT {
			fullSpec[v] = v
		}
	case renameSheetSpec:
		for k, v := range specT {
			if v == nil {
				fullSpec[k] = k
			} else {
				fullSpec[k] = v
			}
		}
	}

	fullData := map[string]string{}
	switch dataT := data.(type) {
	case map[string]interface{}:
		for k, v := range dataT {
			fullData[toDataName(k)] = v.(string)
		}
	case map[string]string:
		for k, v := range dataT {
			fullData[toDataName(k)] = v
		}
	}

	sheet = SheetData{
		CellAddress{0, 0}: CellData{title, StopIfNotEqual},
		CellAddress{1, 0}: CellData{"Timestamp", StopIfNotEqual},
		CellAddress{2, 0}: CellData{timestamp, IncrementRowUntilEmpty},
	}

	From(fullSpec).ForEachIndexed(func(i int, v interface{}) {
		header := CellAddress{1, i + 1 + colOffset}
		cell := CellAddress{2, i + 1 + colOffset}

		name := v.(KeyValue).Key.(string)
		dataName := toDataName(v.(KeyValue).Value.(string))

		sheet[header] = CellData{name, StopIfNotEqual}
		sheet[cell] = CellData{fullData[dataName], IncrementRowUntilEmpty}
	})

	return
}

type simpleSheetSpec []string
type renameSheetSpec map[string]interface{}

func toDataName(in string) (out string) {
	stripRE := regexp.MustCompile(`\W`)
	out = strings.ToUpper(stripRE.ReplaceAllString(in, ""))
	return
}
