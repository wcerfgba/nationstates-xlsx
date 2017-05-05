package output_spec

import (
	"nationstates-xlsx/input_spec"

	"os"

	"github.com/tealeg/xlsx"
)

type T20170429 struct {
}

func (s *T20170429) Parse(in input_spec.InputData) (out OutputData) {
	out = OutputData{
		"Overview": SheetData{
			Cell{0, 0}: "Statistic",
			Cell{0, 1}: "Value",
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
		// "Deaths": buildSheetData(in, "Deaths", []string{
		// 	"Acts of God",
		// }),
	}
	for k, v := range buildPartialSheetData(in, Cell{1, 0}, []string{
		"GDP",
		"INCOME",
		"PUBLICSECTOR",
	}) {
		out["Overview"][k] = v
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
			sheet.Cell(cell.row, cell.col).Value = v
		}
	}
	err = f.Save(filename)
	return
}

func (s *T20170429) Append(data OutputData, filename string) (err error) {

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
		next.col++
		if in[k] == nil {
			data[next] = ""
		} else {
			data[next] = in[k].(string)
		}
		next.col--
		next.row++
	}
	return
}

func buildSheetData(in input_spec.InputData, sheet string, keys []string) (data SheetData) {
	data = buildPartialSheetData(in[sheet].(map[string]interface{}), Cell{0, 0}, keys)
	return
}
