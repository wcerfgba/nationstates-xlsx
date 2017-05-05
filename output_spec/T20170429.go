package output_spec

import (
	"github.com/wcerfgba/nationstates-xlsx/input_spec"

	"os"

	"github.com/tealeg/xlsx"
)

type T20170429 struct {
}

func (s *T20170429) Parse(in input_spec.InputData) (out OutputData) {
	out = OutputData{
		"Overview": SheetData{
			Cell{0, 0, Skip}: "Statistic",
			Cell{0, 1, Skip}: "Value",
		},
		"Government": buildSheetData(in, "GOVT", []string{
			"ADMINISTRATION",
			"DEFENCE",
			"EDUCATION",
			"ENVIRONMENT",
			"HEALTHCARE",
			"COMMERCE",
			"INTERNATIONALAID",
			"LAWANDORDER",
			"PUBLICTRANSPORT",
			"SOCIALEQUALITY",
			"SPIRITUALITY",
			"WELFARE",
		}),
		"Sectors": buildSheetData(in, "SECTORS", []string{
			"BLACKMARKET",
			"GOVERNMENT",
			"INDUSTRY",
			"PUBLIC",
		}),
		"Freedom Scores": buildSheetData(in, "FREEDOMSCORES", []string{
			"CIVILRIGHTS",
			"ECONOMY",
			"POLITICALFREEDOM",
		}),
		"Causes of death": SheetData{
			Cell{1, 0, StopIfNotEqual}:         "Timestamp",
			Cell{1, 1, StopIfNotEqual}:         "Old age",
			Cell{1, 2, StopIfNotEqual}:         "Heart disease",
			Cell{1, 3, StopIfNotEqual}:         "Murder",
			Cell{1, 4, StopIfNotEqual}:         "Cancer",
			Cell{1, 5, StopIfNotEqual}:         "Acts of God",
			Cell{1, 6, StopIfNotEqual}:         "Capital Punishment",
			Cell{1, 7, StopIfNotEqual}:         "Exposure",
			Cell{2, 0, IncrementRowUntilEmpty}: in["_timestamp"].(string),
			Cell{2, 1, IncrementRowUntilEmpty}: in["DEATHS"].(map[string]string)["Old Age"],
			Cell{2, 2, IncrementRowUntilEmpty}: in["DEATHS"].(map[string]string)["Heart Disease"],
			Cell{2, 3, IncrementRowUntilEmpty}: in["DEATHS"].(map[string]string)["Murder"],
			Cell{2, 4, IncrementRowUntilEmpty}: in["DEATHS"].(map[string]string)["Cancer"],
			Cell{2, 5, IncrementRowUntilEmpty}: in["DEATHS"].(map[string]string)["Acts of God"],
			Cell{2, 6, IncrementRowUntilEmpty}: in["DEATHS"].(map[string]string)["Capital Punishment"],
			Cell{2, 7, IncrementRowUntilEmpty}: in["DEATHS"].(map[string]string)["Exposure"],
		},
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
			sheet.Cell(cell.Row, cell.Col).Value = v
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

func buildPartialSheetData(in input_spec.InputData, start Cell, keys []string) (data SheetData) {
	data = SheetData{}
	next := start
	for _, k := range keys {
		data[next] = k
		next.Col++
		if in[k] == nil {
			data[next] = ""
		} else {
			data[next] = in[k].(string)
		}
		next.Col--
		next.Row++
	}
	return
}

func buildSheetData(in input_spec.InputData, sheet string, keys []string) (data SheetData) {
	data = buildPartialSheetData(in[sheet].(map[string]interface{}), Cell{0, 0, Skip}, keys)
	return
}
