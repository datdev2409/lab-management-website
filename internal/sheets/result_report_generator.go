package sheets

import (
	"context"
	"fmt"
	"io"

	"github.com/datdev2409/lab-admin-go/internal/models"
)

type ResultReport struct {
	*BaseReportBuilder
}

func NewResultReport(ctx context.Context) (*ResultReport, error) {
	pageSetup := &PageSetup{
		SheetName:   "Sheet1",
		PageSize:    9,
		Orientation: "portrait",
		Margins: MarginConfig{
			Top:    1.9,
			Bottom: float64(0),
			Left:   0.65,
			Right:  0.65,
			Header: 0.236220472440945,
			Footer: float64(0),
		},
		ColumnWidth: map[string]float64{
			"A": 9.25,
			"B": 37.1,
			"C": 14.5,
			"D": 9.0,
			"E": 20.5,
		},
	}

	builder, err := NewBaseReportBuilder(ctx, pageSetup)
	if err != nil {
		return nil, err
	}

	if err := builder.InitializeNewFile(ctx); err != nil {
		return nil, err
	}

	return &ResultReport{BaseReportBuilder: builder}, nil
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
	patientTable := NewPatientInfoTable(f, sm, &record.Patient, 2, "A")
	if err := patientTable.Apply(ctx); err != nil {
		return nil, err
	}

	testTableStartRow := patientTable.GetEndRow() + 2
	// Create and apply the test result table component
	testTable := NewTestResultTable(f, sm, testTableStartRow, "A", record.TestResults)
	if err := testTable.Apply(ctx); err != nil {
		return nil, err
	}

	// Create and apply signature component with custom config for result report
	// Result report template already has signature name, so we need to override it
	startSignatureRow := testTable.GetEndRow() + 2
	signature := NewSignatureComponentWithConfig(f, sm, "Sheet1", startSignatureRow, 'C', 'E', SignatureConfig{
		IncludeDate:        true, // Result report doesn't include date in signature
		SignatureSpace:     5,    // 5 rows between lab dept and signature name
		WriteSignatureName: true, // Override the signature name in template
		Date:               record.CreatedAt,
	})
	if err := signature.Apply(ctx); err != nil {
		return nil, err
	}

	// Calculate print area based on content (A1 to F + last row with data)
	lastRow := testTable.GetEndRow() + 9 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$E$%d", lastRow)

	err := r.ApplyPrintArea(ctx, f, printArea)
	if err != nil {
		return nil, err
	}

	return r.GetIOReader(ctx)
}
