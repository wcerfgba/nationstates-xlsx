package data

import (
	"os"

	"github.com/tealeg/xlsx"
	"github.com/wcerfgba/nationstates-xlsx/util"
)

// XlsxAppendingOutputter fulfils the Outputter interface and writes data to an
// XLSX spreadsheet. It is capable of reading an existing spreadsheet of a
// particular form, selecting existing sheets and columns from that spreadsheet,
// and appending a new data row with data in the appropriate columns.
type XlsxAppendingOutputter struct {
	config struct {
		outputFilename string
	}
}

// Output just writes the values to the cells for each sheet in the input
// Result. The real magic happens in getCells.
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

// getCells returns a map keyed by row and column whose values should be written
// to the appropriate cell in the spreadsheet. It assumes the data is laid out
// such that the second row of the given xlsx.Sheet is the header, with each
// cell in the header providing the name of the data in that columns, and that
// the data is appended in rows downwards from the header cell.
func getCells(data Result, sheet *xlsx.Sheet) (cells map[cell]string, err error) {
	cells = map[cell]string{}
	headerRow := 1
	// Make sure we don't overwrite any data, get the first empty row.
	dataRow := getNextDataRow(sheet)
	for name, node := range data.ChildrenByKey() {
		headerCell := cell{headerRow, 0}
		// For each data element in the Result, advance the column to the right
		// until we find one with a matching name, or a blank column.
		for sheet.Cell(headerCell.Row, headerCell.Col).Value != name &&
			(sheet.Cell(headerCell.Row, headerCell.Col).Value != "" ||
				cells[headerCell] != "") &&
			headerCell.Col <= len(data.Children()) {

			headerCell.Col = headerCell.Col + 1
		}
		// Store header and data.
		cells[cell{headerCell.Row, headerCell.Col}] = name
		dataCell := cell{dataRow, headerCell.Col}
		cells[cell{dataCell.Row, dataCell.Col}] = node.Value()
	}
	return cells, nil
}

// cell just stored a Row and Col index for use in the intermediate map.
type cell struct {
	Row, Col int
}

// getNextDataRow seeks downwards from the third row of the given xlsx.Sheet
// until it finds an empty row, and returns its index.
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
