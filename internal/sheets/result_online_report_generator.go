package sheets

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
)

type ResultOnlineReport struct {
	*PageSetup
	*ReportFile
}

func NewResultOnlineReport(ctx context.Context) (*ResultOnlineReport, error) {
	report := &ResultOnlineReport{
		ReportFile: &ReportFile{
			File: nil,
		},
		PageSetup: &PageSetup{
			SheetName:   "Sheet1",
			PageSize:    9,
			Orientation: "portrait",
			Margins: MarginConfig{
				Top:    0.31496063,
				Bottom: float64(0),
				Left:   0.4,
				Right:  0.4,
				Header: 0.511811023622047,
				Footer: float64(0),
			},
			ColumnWidth: map[string]float64{
				"A": 0,
				"B": 9.0,
				// "C": 37.17,
				"C": 36.83,
				// "D": 16.83,
				"D": 15.67,
				// "E": 13.0,
				"E": 12.0,
				"F": 27.0,
			},
		},
	}

	err := report.OpenTemplate(ctx, "templates/PhieuKetQuaChuKy.xlsx")
	if err != nil {
		return nil, err
	}

	err = report.ApplyColumnWidths(ctx, report.File)
	if err != nil {
		return nil, err
	}

	err = report.ApplyPageSetupV2(ctx, report.File)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r ResultOnlineReport) Generate(ctx context.Context, data interface{}) (io.Reader, error) {
	record, ok := data.(*models.Record)
	if !ok {
		return nil, fmt.Errorf("invalid data type for result report generation")
	}

	f := r.File
	defer f.Close()

	sm := NewStyleManager(ctx, f)

	cells := map[string]Cell{
		// "C4": {
		// 	value:     "                        Hotline/Zalo:       0919 663 747                                 Phone: 0833 657 774",
		// 	styleName: GetStyleNamePtr(LabContactStyle),
		// },
		"C9": {
			value:     time.Now().Format("02/01/2006"),
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
		"C10": {
			value:     record.Patient.Name,
			styleName: GetStyleNamePtr(PatientNameStyle),
		},
		"C11": {
			value:     record.Patient.Address,
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
		"E9": {
			value:     record.Patient.Phone,
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
		"E10": {
			value:     record.Patient.YOB,
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
		"E11": {
			value:     record.Patient.Gender,
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
		"D10": {
			value:     "Năm sinh",
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
	}

	for cell, config := range cells {
		if err := f.SetCellValue("Sheet1", cell, config.value); err != nil {
			return nil, err
		}
		if config.styleName != nil {
			err := f.SetCellStyle("Sheet1", cell, cell, sm.GetStyleV2(*config.styleName))
			if err != nil {
				return nil, err
			}
		}
	}

	// Create and apply the test result table component
	testTable := NewTestResultTable(f, sm, 15, "B", record.TestResults)
	if err := testTable.Apply(ctx); err != nil {
		return nil, err
	}

	// Calculate print area based on content (A1 to F + last row with data)
	lastRow := testTable.GetEndRow() + 2 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$F$%d", lastRow+7)

	err := r.ApplyPrintArea(ctx, f, printArea)
	if err != nil {
		return nil, err
	}

	return r.GetIOReader(ctx)
}
