package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"jinovatka/entities"
	"net/url"
	"path"

	"github.com/xuri/excelize/v2"
)

type ExporterService struct {
	ServerHost          string
	SeedDetailPath      string
	WaybackRedirectPath string
}

func NewExporterService(settings *ServiceSettings) *ExporterService {
	return &ExporterService{
		ServerHost:          settings.ServerHost,
		SeedDetailPath:      settings.SeedDetailPath,
		WaybackRedirectPath: settings.WaybackRedirectPath,
	}
}

// Convert SeedsGroup to nice excel sheet for users to keep track of their submited seeds.
// The excel data will be written to the provided io.Writer.
func (service ExporterService) GroupToExcel(group *entities.SeedsGroup, w io.Writer /*seedUrlPrefix *url.URL*/) error {
	header := []any{"URL", "Zkrácený odkaz do Webarchivu", "Odkaz na detail", "Stav", "Odkaz do Webarchivu"}
	const sheet = "Export"

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

		// fmt.Println("string:", seedUrlPrefix.String(), "host:", seedUrlPrefix.Host, "path:", seedUrlPrefix.Path)
		// detailLink := seedUrlPrefix.JoinPath("/" + seed.ShadowID)

		detailLink, err := url.JoinPath(service.ServerHost, path.Join(service.SeedDetailPath, seed.ShadowID))
		if err != nil {
			return fmt.Errorf("ExporterService.GroupToExcel could not create seed detail url: %w", err)
		}

		waybackShortLink, err := url.JoinPath(service.ServerHost, path.Join(service.WaybackRedirectPath, seed.ShadowID))
		if err != nil {
			return fmt.Errorf("ExporterService.GroupToExcel could not create wayback redirect url: %w", err)
		}

		// fmt.Println("string:", detailLink.String(), "host:", detailLink.Host, "path:", detailLink.Path)

		state := entities.PrettyPrintCaptureState(seed.State)

		row := []any{
			possibleLink{IsLink: true, Value: seed.URL, Link: seed.URL},                 // URL
			possibleLink{IsLink: true, Value: waybackShortLink, Link: waybackShortLink}, // Zkrácený odkaz do waybacku
			possibleLink{IsLink: true, Value: seed.ShadowID, Link: detailLink},          // Odkaz na detail
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

func (service *ExporterService) GroupToCsv(group *entities.SeedsGroup, writer io.Writer) error {
	header := []string{"URL", "Zkrácený odkaz do Webarchivu", "Odkaz na detail", "Stav", "Odkaz do Webarchivu"}

	csvWriter := csv.NewWriter(writer)

	err := csvWriter.Write(header)
	if err != nil {
		return fmt.Errorf("ExporterService.GroupToCsv could not write header: %w", err)
	}

	for _, seed := range group.Seeds {
		detailLink, err := url.JoinPath(service.ServerHost, path.Join(service.SeedDetailPath, seed.ShadowID))
		if err != nil {
			return fmt.Errorf("ExporterService.GroupToCsv could not create seed detail url: %w", err)
		}

		waybackShortLink, err := url.JoinPath(service.ServerHost, path.Join(service.WaybackRedirectPath, seed.ShadowID))
		if err != nil {
			return fmt.Errorf("ExporterService.GroupToCsv could not create wayback redirect url: %w", err)
		}

		row := []string{
			seed.URL,
			waybackShortLink,
			detailLink,
			entities.PrettyPrintCaptureState(seed.State),
			seed.ArchivalURL,
		}

		err = csvWriter.Write(row)
		if err != nil {
			return fmt.Errorf("ExporterService.GroupToCsv could not csv row: %w", err)
		}
	}

	csvWriter.Flush()
	if csvWriter.Error() != nil {
		return fmt.Errorf("ExporterService.GroupToCsv flush returned error: %w", err)
	}

	return nil
}
