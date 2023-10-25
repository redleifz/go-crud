package excelutil

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

// CreateExcelFile creates an Excel file from data and columns and saves it with the given filename.
func CreateExcelFile(data []string, columns []string, filename string) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return err
	}

	// Add headers
	headerRow := sheet.AddRow()
	for _, colName := range columns {
		cell := headerRow.AddCell()
		cell.Value = colName
	}

	// Add data rows
	for _, rowData := range data {
		dataRow := sheet.AddRow()
		for _, cellValue := range rowData {
			cell := dataRow.AddCell()
			cell.Value = fmt.Sprintf("%v", cellValue)
		}
	}

	err = file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}
