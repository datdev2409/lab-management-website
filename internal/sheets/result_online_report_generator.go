package sheets

import (
	"context"
	"fmt"
	_ "image/jpeg"
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
			// "A": 0,
			"A": 9.25,
			// "C": 37.17,
			"B": 37.1,
			// "D": 16.83,
			"C": 14.5,
			// "E": 13.0,
			"D": 9.0,
			"E": 20.5,
		},
	}

	builder, err := NewBaseReportBuilder(ctx, pageSetup)
	if err != nil {
		return nil, err
	}

	// if err := builder.InitializeFromTemplate(ctx, "templates/PhieuKetQuaChuKy.xlsx"); err != nil {
	// 	return nil, err
	// }
	if err := builder.InitializeNewFile(ctx); err != nil {
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

	f.MergeCell("Sheet1", "A1", "E1")
	f.MergeCell("Sheet1", "A2", "E2")

	r.MergeCellsWithStyle("Sheet1", "A1", "E1", "PHÒNG XÉT NGHIỆM Y KHOA ANH QUÂN", sm.GetStyleV2(ReportHeaderLabNameStyle))
	r.MergeCellsWithStyle("Sheet1", "A2", "E2", "60 Đống Đa, Phường Cao Lãnh, Đồng Tháp", sm.GetStyleV2(ReportHeaderLabAddressStyle))
	r.MergeCellsWithStyle("Sheet1", "A3", "E3", "Hotline/Zalo: 0919 663 747                                        Phone: 0833 657 774", sm.GetStyleV2(ReportHeaderLabInfoStyle))
	r.MergeCellsWithStyle("Sheet1", "A4", "E4", "         Email: laboanhquan@gmail.com                             Email: nguyencongman@gmail.com", sm.GetStyleV2(ReportHeaderLabInfoStyle))

	err := AddLogoComponent(ctx, f, "Sheet1", "A1", "./assets/anhquanlab_logo.png", 0.6, 0.5)
	if err != nil {
		return nil, err
	}

	r.MergeCellsWithStyle("Sheet1", "A6", "E6", "PHIẾU KẾT QUẢ XÉT NGHIỆM", sm.GetStyleV2(ResultReportReportNameStyle))
	f.SetRowHeight("Sheet1", 6, 30.0)

	// Create and apply the patient info component
	patientTable := NewPatientInfoTable(f, sm, &record.Patient, 7, "A")
	if err := patientTable.Apply(ctx); err != nil {
		return nil, err
	}

	resultTableStartRow := patientTable.GetEndRow() + 2

	// Create and apply the test result table component
	testTable := NewTestResultTable(f, sm, resultTableStartRow, "A", record.TestResults)
	if err := testTable.Apply(ctx); err != nil {
		return nil, err
	}

	startSignatureRow := testTable.GetEndRow() + 2
	signature := NewSignatureComponentWithConfig(f, sm, "Sheet1", startSignatureRow, 'C', 'E', SignatureConfig{
		IncludeDate:         true, // Result report doesn't include date in signature
		SignatureSpace:      7,    // 5 rows between lab dept and signature name
		WriteSignatureName:  true, // Override the signature name in template
		WriteSignatureImage: true,
		SignatureImagePath:  "./assets/signature.jpg",
	})
	if err := signature.Apply(ctx); err != nil {
		return nil, err
	}

	// Calculate print area based on content (A1 to F + last row with data)
	lastRow := signature.GetEndRow() + 3 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$E$%d", lastRow)

	err = r.ApplyPrintArea(ctx, f, printArea)
	if err != nil {
		return nil, err
	}

	return r.GetIOReader(ctx)
}
