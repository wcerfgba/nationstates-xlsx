package data

import (
	"os"

	"github.com/tealeg/xlsx"
	"github.com/wcerfgba/nationstates-xlsx/util"
)

type XlsxAppendingOutputter struct {
	config struct {
		outputFilename string
	}
}

func (o *XlsxAppendingOutputter) Output(res Result) (err error) {

	f, err := getFile(o.config.outputFilename)
	if err != nil {
		return
	}

	for _, inSheet := range res.Children() {
		outSheet, err := getSheet(inSheet.Key(), f)
		if err != nil {
			return err
		}
		cells, err := getCells(inSheet, outSheet)
		if err != nil {
			return err
		}
		for cell, val := range cells {
			outSheet.Cell(cell.Row, cell.Col).SetValue(val)
		}
	}

	err = f.Save(o.config.outputFilename)

	return
}

func (o *XlsxAppendingOutputter) Configure(c util.Configuration) {
	o.config.outputFilename = c["outputFilename"].(string)
}

func getFile(name string) (*xlsx.File, error) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return xlsx.NewFile(), nil
	}
	return xlsx.OpenFile(name)
}

func getSheet(name string, f *xlsx.File) (*xlsx.Sheet, error) {
	for _, sheet := range f.Sheets {
		if sheet.Name == name {
			return sheet, nil
		}
	}
	return f.AddSheet(name)
}

func getCells(data Result, sheet *xlsx.Sheet) (cells map[Cell]string, err error) {
	cells = map[Cell]string{}
	headerRow := 1
	dataRow := getNextDataRow(sheet)
	for name, node := range data.ChildrenByKey() {
		headerCell := Cell{headerRow, 0}
		for sheet.Cell(headerCell.Row, headerCell.Col).Value != name &&
			(sheet.Cell(headerCell.Row, headerCell.Col).Value != "" ||
				cells[headerCell] != "") &&
			headerCell.Col <= len(data.Children()) {

			headerCell.Col = headerCell.Col + 1
		}
		cells[Cell{headerCell.Row, headerCell.Col}] = name
		dataCell := Cell{dataRow, headerCell.Col}
		cells[Cell{dataCell.Row, dataCell.Col}] = node.Value()
	}
	return cells, nil
}

type Cell struct {
	Row, Col int
}

func getNextDataRow(sheet *xlsx.Sheet) int {
	row := 2
rows:
	for ; row < sheet.MaxRow+1; row++ {
		for col := 0; col < sheet.MaxCol; col++ {
			if sheet.Cell(row, col).Value != "" {
				continue rows
			}
		}
		return row
	}
	return row
}

// func (s *T20170429) Create(data OutputData, filename string) (err error) {
// 	f := xlsx.NewFile()
// 	for sheet, sheetData := range data {
// 		sheet, err := getSheet(f, sheet)
// 		if err != nil {
// 			return err
// 		}
// 		for cell, v := range sheetData {
// 			sheet.Cell(cell.Row, cell.Col).SetValue(v.Contents)
// 		}
// 	}
// 	err = f.Save(filename)
// 	return
// }
