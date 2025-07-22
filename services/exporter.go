package services

import (
	"fmt"
	"io"
	"jinovatka/entities"
	"net/url"

	"github.com/xuri/excelize/v2"
)

type ExporterService struct{}

func NewExporterService() *ExporterService {
	return &ExporterService{}
}

// Convert SeedsGroup to nice excel sheet for users to keep track of their submited seeds.
// The excel data will be written to the provided io.Writer.
func (service ExporterService) GroupToExcel(group *entities.SeedsGroup, w io.Writer, seedUrlPrefix *url.URL) error {
	header := []any{"URL", "Odkaz na detail", "Stav", "Odkaz do Webarchivu"}
	const sheet = "Semínka"

	f := excelize.NewFile()
	defer f.Close()

	defaultSheet := f.GetSheetName(f.GetActiveSheetIndex())
	err := f.SetSheetName(defaultSheet, sheet)
	if err != nil {
		return fmt.Errorf("ExporterService.GroupToExcel could not rename sheet: %w", err)
	}

	err = service.writeExcelRow(f, sheet, 1, header)
	if err != nil {
		return fmt.Errorf("ExporterService.GroupToExcel could not write header to sheet: %w", err)
	}

	for i, seed := range group.Seeds {
		rowIndex := i + 2 // This is excel, data starts at row 2 :)
		fmt.Println("string:", seedUrlPrefix.String(), "host:", seedUrlPrefix.Host, "path:", seedUrlPrefix.Path)
		detailLink := seedUrlPrefix.JoinPath("/" + seed.ShadowID)
		fmt.Println("string:", detailLink.String(), "host:", detailLink.Host, "path:", detailLink.Path)

		state := "Nesklizeno"
		if seed.State == entities.HarvestedSucessfully {
			state = "Sklizeno"
		}
		row := []any{
			possibleLink{IsLink: true, Value: seed.URL, Link: seed.URL},                 // URL
			possibleLink{IsLink: true, Value: seed.ShadowID, Link: detailLink.String()}, // Odkaz na detail
			possibleLink{IsLink: false, Value: state},                                   // Stav
			possibleLink{IsLink: true, Value: seed.ArchivalURL, Link: seed.ArchivalURL}, // Odkaz do Webarchivu
		}
		err = service.writeExcelRow(f, sheet, rowIndex, row)
		if err != nil {
			return fmt.Errorf("ExporterService.GroupToExcel could not write row to sheet: %w", err)
		}
	}

	_, err = f.WriteTo(w)
	if err != nil {
		return fmt.Errorf("ExporterService.GroupToExcel could not write sheet into writer: %w", err)
	}

	return nil
}

type possibleLink struct {
	IsLink bool
	Value  any
	Link   string
}

// Rows are indexed from 1!
func (service ExporterService) writeExcelRow(f *excelize.File, sheet string, rowIndex int, row []any) error {
	for i := range row {
		colIndex := i + 1
		cellName, err := excelize.CoordinatesToCellName(colIndex, rowIndex)
		if err != nil {
			return err
		}

		// There is definetly a way to do this more cleanly. I am certain of it. But who cares actually? If you do then please open an issue.
		if cellValue, ok := row[i].(possibleLink); ok {
			err = f.SetCellValue(sheet, cellName, cellValue.Value)
			if err != nil {
				return err
			}
			if cellValue.IsLink {
				err = f.SetCellHyperLink(sheet, cellName, cellValue.Link, "External")
				if err != nil {
					return err
				}
			}
		} else {
			err = f.SetCellValue(sheet, cellName, row[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
