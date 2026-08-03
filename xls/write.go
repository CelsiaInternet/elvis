package xls

import (
	"io"
	"net/http"
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
* build: Creates the excelize.File populated with all the workbook's sheets.
* @return *excelize.File, error
**/
func (s *Xls) build() (*excelize.File, error) {
	f := excelize.NewFile()

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
			return nil, err
		}

		for col, header := range sheet.Headers {
			cell, err := excelize.CoordinatesToCellName(col+1, 1)
			if err != nil {
				return nil, err
			}
			if err := f.SetCellValue(sheetName, cell, header); err != nil {
				return nil, err
			}
		}

		for rowIdx, row := range sheet.Rows {
			for col, header := range sheet.Headers {
				cell, err := excelize.CoordinatesToCellName(col+1, rowIdx+2)
				if err != nil {
					return nil, err
				}
				if err := f.SetCellValue(sheetName, cell, row[header]); err != nil {
					return nil, err
				}
			}
		}
	}

	if !usedDefault && len(s.Sheets) > 0 {
		if err := f.DeleteSheet(defaultSheet); err != nil {
			return nil, err
		}
	}

	f.SetActiveSheet(0)

	return f, nil
}

/**
* ToFile: Exports the workbook to an xlsx file at the given path.
* @param path string
* @return error
**/
func (s *Xls) ToFile(path string) error {
	f, err := s.build()
	if err != nil {
		return err
	}
	defer f.Close()

	return f.SaveAs(path)
}

/**
* ToWriter: Writes the workbook as xlsx data into the given writer.
* @param w io.Writer
* @return error
**/
func (s *Xls) ToWriter(w io.Writer) error {
	f, err := s.build()
	if err != nil {
		return err
	}
	defer f.Close()

	return f.Write(w)
}

/**
* ToHttp: Writes the workbook as a downloadable xlsx attachment to the http response.
* @param w http.ResponseWriter
* @param filename string
* @return error
**/
func (s *Xls) ToHttp(w http.ResponseWriter, filename string) error {
	f, err := s.build()
	if err != nil {
		return err
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	return f.Write(w)
}
