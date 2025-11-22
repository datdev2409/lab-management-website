package sheets

import (
	"context"
	"fmt"
	"io"
	"sort"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/xuri/excelize/v2"
)

// TrackingReportData contains the data needed for tracking report generation
type TrackingReportData struct {
	Records  []*models.Record  `json:"records"`
	TestList []models.TestInfo `json:"test_list"`
}

type TrackingReport struct {
	*BaseReportBuilder
}

func NewTrackingReport(ctx context.Context) (*TrackingReport, error) {
	pageSetup := &PageSetup{
		SheetName:   "Sheet1",
		PageSize:    9,
		Orientation: "landscape",
		Margins: MarginConfig{
			Top:    0.25,
			Bottom: 0.5,
			Left:   0.5,
			Right:  0.5,
			Header: 0.5,
			Footer: 0.5,
		},
		ColumnWidth: map[string]float64{
			"A": 4.0,  // STT column
			"B": 38.0, // Test name column
			"C": 20.0, // Normal range column
		},
	}

	builder, err := NewBaseReportBuilder(ctx, pageSetup)
	if err != nil {
		return nil, err
	}

	if err := builder.InitializeNewFile(ctx); err != nil {
		return nil, err
	}

	return &TrackingReport{BaseReportBuilder: builder}, nil
}

func (r *TrackingReport) Generate(ctx context.Context, data interface{}) (io.Reader, error) {
	// Type assertion
	trackingData, ok := data.(*TrackingReportData)
	if !ok {
		return nil, fmt.Errorf("invalid data type for tracking report")
	}

	records := trackingData.Records
	testList := trackingData.TestList

	// Sort records from oldest to latest (first column = oldest record)
	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt.Before(records[j].CreatedAt)
	})

	f := r.File
	defer f.Close()

	// Create style manager
	sm := NewStyleManager(ctx, f)

	// Set dynamic column widths for record columns (D onwards)
	for i := range records {
		recordCol := string('D' + rune(i))
		err := f.SetColWidth("Sheet1", recordCol, recordCol, 16.0)
		if err != nil {
			return nil, fmt.Errorf("failed to set column width for %s: %v", recordCol, err)
		}
	}

	// Build header section
	err := r.buildReportHeader(f, sm, records)
	if err != nil {
		return nil, err
	}

	// Build test table
	err = r.buildTestTable(f, sm, records, testList)
	if err != nil {
		return nil, err
	}

	// Build signature section
	tableEndRow := 7 + len(testList) - 1 // Last row of test data
	signature := NewSignatureComponentWithConfig(f, sm, "Sheet1", tableEndRow+2, 'B', 'B', SignatureConfig{
		IncludeDate:        true,
		SignatureSpace:     5,
		WriteSignatureName: true,
		Date:               GetVietnamTime(),
	})
	if err := signature.Apply(ctx); err != nil {
		return nil, err
	}

	if err := r.ApplyPageSetupV2(ctx, r.File); err != nil {
		return nil, err
	}

	return r.GetIOReader(ctx)
}

// buildReportHeader creates the header section of the tracking report
func (r *TrackingReport) buildReportHeader(f *excelize.File, sm *StyleManager, records []*models.Record) error {
	// Calculate the last column based on number of records (A, B, C + records)
	lastCol := string('C' + rune(len(records))) // C + number of record columns

	// Set row heights for header rows
	_ = f.SetRowHeight("Sheet1", 1, 16.0) // Row 1 height: 16
	_ = f.SetRowHeight("Sheet1", 2, 16.0) // Row 2 height: 16
	_ = f.SetRowHeight("Sheet1", 3, 24.0) // Row 3 height: 20

	// Lab name - left aligned, no merge needed
	_ = f.SetCellValue("Sheet1", "A1", "PHÒNG XÉT NGHIỆM Y KHOA ANH QUÂN")
	_ = f.SetCellStyle("Sheet1", "A1", "A1", sm.GetStyleV2(LabNameLeftStyle))

	// Lab address - left aligned, no merge needed
	_ = f.SetCellValue("Sheet1", "A2", "60 Đống Đa, Phường Cao Lãnh, Đồng Tháp.")
	_ = f.SetCellStyle("Sheet1", "A2", "A2", sm.GetStyleV2(LabAddressLeftStyle))

	// Report title - merge across all columns and center
	_ = f.MergeCell("Sheet1", "A3", fmt.Sprintf("%s3", lastCol))
	_ = f.SetCellValue("Sheet1", "A3", "SỔ THEO DÕI KẾT QUẢ XÉT NGHIỆM")
	_ = f.SetCellStyle("Sheet1", "A3", "A3", sm.GetStyleV2(Font16BoldCenterStyle))

	// Patient name - merge across all columns and center
	_ = f.MergeCell("Sheet1", "A4", fmt.Sprintf("%s4", lastCol))
	_ = f.SetCellValue("Sheet1", "A4", fmt.Sprintf("Họ & Tên: %s", records[0].Patient.Name))
	_ = f.SetCellStyle("Sheet1", "A4", "A4", sm.GetStyleV2(PatientNameCenterStyle))

	return nil
}

// buildTestTable creates the test results table
func (r *TrackingReport) buildTestTable(f *excelize.File, sm *StyleManager, records []*models.Record, testList []models.TestInfo) error {
	// Table starts at row 6
	headerRow := 6
	dataStartRow := 7

	// Set table headers
	headers := []string{"STT", "Tên dịch vụ", "Khoảng tham chiếu"}

	// Add date headers for each record
	for _, record := range records {
		headers = append(headers, ToVietnamTime(record.CreatedAt).Format("02/01/2006"))
	}

	// Set header row
	for i, header := range headers {
		col := string('A' + rune(i))
		cell := fmt.Sprintf("%s%d", col, headerRow)
		_ = f.SetCellValue("Sheet1", cell, header)
		_ = f.SetCellStyle("Sheet1", cell, cell, sm.GetStyleV2(TrackingTableHeaderCyanStyle))
	}

	// Build test rows
	for i, testInfo := range testList {
		currentRow := dataStartRow + i

		// STT (sequential number)
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("A%d", currentRow), i+1)
		_ = f.SetCellStyle("Sheet1", fmt.Sprintf("A%d", currentRow), fmt.Sprintf("A%d", currentRow), sm.GetStyleV2(TestIndexStyle))

		// Test name
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("B%d", currentRow), testInfo.Name)
		_ = f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", currentRow), fmt.Sprintf("B%d", currentRow), sm.GetStyleV2(TestNameStyle))

		// Normal range
		_ = f.SetCellValue("Sheet1", fmt.Sprintf("C%d", currentRow), testInfo.NormalValue)
		_ = f.SetCellStyle("Sheet1", fmt.Sprintf("C%d", currentRow), fmt.Sprintf("C%d", currentRow), sm.GetStyleV2(TestNormalRangeStyle))

		// Set row height for test data row
		err := f.SetRowHeight("Sheet1", currentRow, 19.0)
		if err != nil {
			return err
		}

		// Fill in test results for each record
		for recordIndex, record := range records {
			resultCol := string('D' + rune(recordIndex))
			resultCell := fmt.Sprintf("%s%d", resultCol, currentRow)

			// Find the test result in this record
			var testResult *models.TestResult
			for _, result := range record.TestResults {
				if result.Name == testInfo.Name {
					testResult = &result
					break
				}
			}

			// Set cell value and style
			if testResult != nil {
				displayValue := formatTestResultDisplay(testResult)
				_ = f.SetCellValue("Sheet1", resultCell, displayValue)

				// Use abnormal style if result is abnormal, otherwise normal style
				if testResult.Abnormal {
					_ = f.SetCellStyle("Sheet1", resultCell, resultCell, sm.GetStyleV2(TestAbnormalResultStyle))
				} else {
					_ = f.SetCellStyle("Sheet1", resultCell, resultCell, sm.GetStyleV2(TestResultStyle))
				}
			} else {
				// Empty cell with normal style
				_ = f.SetCellStyle("Sheet1", resultCell, resultCell, sm.GetStyleV2(TestResultStyle))
			}
		}
	}

	return nil
}

// formatTestResultDisplay formats the test result for display based on Result and ResultText fields
func formatTestResultDisplay(testResult *models.TestResult) string {
	if testResult == nil {
		return ""
	}

	// Determine display value based on Result and ResultText
	if testResult.Result != "" && testResult.ResultText != "" {
		// Both exist: display "FormatResult(Result) (ResultText)"
		return fmt.Sprintf("%s (%s)", FormatResult(testResult.Result), testResult.ResultText)
	} else if testResult.Result != "" {
		// Only Result exists
		return FormatResult(testResult.Result)
	} else if testResult.ResultText != "" {
		// Only ResultText exists
		return testResult.ResultText
	} else {
		// Neither exists (should not happen, but handle gracefully)
		return ""
	}
}
