package xls

import (
	"bytes"
	"sort"

	"github.com/celsiainternet/elvis/et"
	"github.com/xuri/excelize/v2"
)

/**
* Sheet: Holds the data and header order for a single Excel worksheet.
**/
type Sheet struct {
	Name    string
	Headers []string
	Rows    []et.Json
}

/**
* Xls: Aggregates the sheets that make up an Excel workbook.
**/
type Xls struct {
	Sheets []*Sheet
}

/**
* NewXls: Builds an Xls workbook with an initial sheet from a list of Json rows and a sheet name.
* @param data []et.Json
* @param nameSheet string
* @return *Xls
**/
func NewXls(data []et.Json, nameSheet string) *Xls {
	result := &Xls{
		Sheets: []*Sheet{},
	}
	result.Add(data, nameSheet)

	return result
}

/**
* Add: Appends a new sheet to the workbook from a list of Json rows and a sheet name.
* @param data []et.Json
* @param nameSheet string
* @return *Xls
**/
func (s *Xls) Add(data []et.Json, nameSheet string) *Xls {
	headers := []string{}
	seen := map[string]bool{}
	for _, row := range data {
		for key := range row {
			if !seen[key] {
				seen[key] = true
				headers = append(headers, key)
			}
		}
	}
	sort.Strings(headers)

	s.Sheets = append(s.Sheets, &Sheet{
		Name:    nameSheet,
		Headers: headers,
		Rows:    data,
	})

	return s
}

/**
* ToFile: Exports the workbook to an xlsx file at the given path.
* @param path string
* @return error
**/
func (s *Xls) ToFile(path string) error {
	f := excelize.NewFile()
	defer f.Close()

	defaultSheet := f.GetSheetName(0)
	usedDefault := false

	for _, sheet := range s.Sheets {
		sheetName := sheet.Name
		if sheetName == "" {
			sheetName = "Sheet1"
		}

		if sheetName == defaultSheet && !usedDefault {
			usedDefault = true
		} else if _, err := f.NewSheet(sheetName); err != nil {
			return err
		}

		for col, header := range sheet.Headers {
			cell, err := excelize.CoordinatesToCellName(col+1, 1)
			if err != nil {
				return err
			}
			if err := f.SetCellValue(sheetName, cell, header); err != nil {
				return err
			}
		}

		for rowIdx, row := range sheet.Rows {
			for col, header := range sheet.Headers {
				cell, err := excelize.CoordinatesToCellName(col+1, rowIdx+2)
				if err != nil {
					return err
				}
				if err := f.SetCellValue(sheetName, cell, row[header]); err != nil {
					return err
				}
			}
		}
	}

	if !usedDefault && len(s.Sheets) > 0 {
		if err := f.DeleteSheet(defaultSheet); err != nil {
			return err
		}
	}

	f.SetActiveSheet(0)

	return f.SaveAs(path)
}

/**
* XlsReader: Wraps an opened Excel workbook for reading sheet data.
**/
type XlsReader struct {
	file *excelize.File
}

/**
* ReadXls: Opens an Excel workbook from raw bytes for reading.
* @param data []byte
* @return *XlsReader, error
**/
func ReadXls(data []byte) (*XlsReader, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return &XlsReader{file: f}, nil
}

/**
* Close: Releases the resources held by the opened workbook.
* @return error
**/
func (s *XlsReader) Close() error {
	return s.file.Close()
}

/**
* GetSheet: Returns the rows of a sheet as a list of Json objects. If columns is empty, all columns of the sheet are returned.
* @param nameSheet string
* @param columns []string
* @return []et.Json, error
**/
func (s *XlsReader) GetSheet(nameSheet string, columns []string) ([]et.Json, error) {
	rows, err := s.file.GetRows(nameSheet)
	if err != nil {
		return nil, err
	}

	result := []et.Json{}
	if len(rows) == 0 {
		return result, nil
	}

	headers := rows[0]
	selected := columns
	if len(selected) == 0 {
		selected = headers
	}

	for _, row := range rows[1:] {
		item := et.Json{}
		for _, col := range selected {
			idx := indexOf(headers, col)
			if idx == -1 || idx >= len(row) {
				continue
			}
			item[col] = row[idx]
		}
		result = append(result, item)
	}

	return result, nil
}

/**
* indexOf: Returns the index of value inside list, or -1 if not found.
* @param list []string
* @param value string
* @return int
**/
func indexOf(list []string, value string) int {
	for i, v := range list {
		if v == value {
			return i
		}
	}
	return -1
}
