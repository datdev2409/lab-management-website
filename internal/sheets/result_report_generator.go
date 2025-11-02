package sheets

import (
	"context"
	"fmt"
	"io"

	"github.com/datdev2409/lab-admin-go/internal/models"
)

type ResultReport struct {
	*PageSetup
	*ReportFile
}

func NewResultReport(ctx context.Context) (*ResultReport, error) {
	report := &ResultReport{
		ReportFile: &ReportFile{
			File: nil,
		},
		PageSetup: &PageSetup{
			SheetName:   "Sheet1",
			PageSize:    9,
			Orientation: "portrait",
			Margins: MarginConfig{
				Top:    1.9,
				Bottom: float64(0),
				Left:   0.4,
				Right:  0.4,
				Header: 0.236220472440945,
				Footer: float64(0),
			},
			ColumnWidth: map[string]float64{
				"A": 0,
				"B": 9.0,
				"C": 36.83,
				"D": 15.67,
				"E": 12.0,
				"F": 28.0,
			},
		},
	}

	err := report.OpenTemplate(ctx, "templates/PhieuKetQua.xlsx")
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

func (r ResultReport) Generate(ctx context.Context, data interface{}) (io.Reader, error) {
	record, ok := data.(*models.Record)
	if !ok {
		return nil, fmt.Errorf("invalid data type for result report generation")
	}

	f := r.File
	defer f.Close()

	sm := NewStyleManager(ctx, f)

	// Create and apply the patient info component
	patientTable := NewPatientInfoTable(f, sm, &record.Patient, 2, "B")
	if err := patientTable.Apply(ctx); err != nil {
		return nil, err
	}

	f.MergeCell("Sheet1", "B5", "C5")

	// Create and apply the test result table component
	testTable := NewTestResultTable(f, sm, 8, "B", record.TestResults)
	if err := testTable.Apply(ctx); err != nil {
		return nil, err
	}

	// Calculate signature section position based on test table end row
	startSignatureRow := testTable.GetEndRow() + 2
	endSignatureRow := startSignatureRow + 5
	startSignatureCell := fmt.Sprintf("D%d", startSignatureRow)

	f.MergeCell("Sheet1", startSignatureCell, fmt.Sprintf("F%d", startSignatureRow))
	f.SetCellValue("Sheet1", startSignatureCell, "PHÒNG XÉT NGHIỆM")
	f.SetCellStyle("Sheet1", startSignatureCell, fmt.Sprintf("F%d", startSignatureRow), sm.GetStyleV2(SignatureStyle))

	f.MergeCell("Sheet1", fmt.Sprintf("D%d", endSignatureRow), fmt.Sprintf("F%d", endSignatureRow))
	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", endSignatureRow), "CKI.XN NGUYỄN CÔNG MẪN")
	f.SetCellStyle("Sheet1", fmt.Sprintf("D%d", endSignatureRow), fmt.Sprintf("F%d", endSignatureRow), sm.GetStyleV2(SignatureStyle))

	// Calculate print area based on content (A1 to F + last row with data)
	lastRow := testTable.GetEndRow() + 9 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$F$%d", lastRow)

	err := r.ApplyPrintArea(ctx, f, printArea)
	if err != nil {
		return nil, err
	}

	return r.GetIOReader(ctx)
}
