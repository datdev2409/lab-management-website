package sheets

import (
	"context"
	"fmt"
	"io"

	"github.com/datdev2409/lab-admin-go/internal/models"
)

type ResultOnlineReport struct {
	*BaseReportBuilder
}

func NewResultOnlineReport(ctx context.Context) (*ResultOnlineReport, error) {
	pageSetup := &PageSetup{
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
	}

	builder, err := NewBaseReportBuilder(ctx, pageSetup)
	if err != nil {
		return nil, err
	}

	if err := builder.InitializeFromTemplate(ctx, "templates/PhieuKetQuaChuKy.xlsx"); err != nil {
		return nil, err
	}

	return &ResultOnlineReport{BaseReportBuilder: builder}, nil
}

func (r ResultOnlineReport) Generate(ctx context.Context, data interface{}) (io.Reader, error) {
	record, ok := data.(*models.Record)
	if !ok {
		return nil, fmt.Errorf("invalid data type for result report generation")
	}

	f := r.File
	defer f.Close()

	sm := NewStyleManager(ctx, f)

	// Create and apply the patient info component
	patientTable := NewPatientInfoTable(f, sm, &record.Patient, 9, "B")
	if err := patientTable.Apply(ctx); err != nil {
		return nil, err
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
