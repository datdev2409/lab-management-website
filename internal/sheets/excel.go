package sheets

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.uber.org/zap"
)

func CreateRecordBillingFile(ctx context.Context, record *models.Record) (string, error) {
	billingReport, err := NewReportWithSingleRecord(models.BillingReport, *record)
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to create billing report config", zap.Error(err))
		return "", err
	}
	f, err := billingReport.Open(ctx)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Create style manager
	styleManager := NewStyleManager(ctx, f)

	// Get styles from style manager
	patientNameStyle, err := styleManager.GetPatientNameStyle()
	if err != nil {
		return "", err
	}

	patientInfoStyle, err := styleManager.GetPatientInfoStyle()
	if err != nil {
		return "", err
	}

	dateCenterStyle, err := styleManager.GetDateCenterStyle()
	if err != nil {
		return "", err
	}

	testResultStyle, err := styleManager.GetTestResultStyle()
	if err != nil {
		return "", err
	}

	testNameStyle, err := styleManager.GetTestNameStyle()
	if err != nil {
		return "", err
	}

	now := time.Now()
	f.SetCellValue("Sheet1", "B4", fmt.Sprintf("Ngày: %s", now.Format("02/01/2006")))
	f.SetCellStyle("Sheet1", "B4", "B4", dateCenterStyle)

	f.SetCellValue("Sheet1", "B6", record.Patient.Name)
	f.SetCellStyle("Sheet1", "B6", "B6", patientNameStyle)

	f.SetCellValue("Sheet1", "B7", record.Patient.Address)
	f.SetCellStyle("Sheet1", "B7", "B7", patientInfoStyle)

	f.SetCellValue("Sheet1", "D6", record.Patient.YOB)
	f.SetCellStyle("Sheet1", "D6", "D6", patientInfoStyle)

	startTestRow := 10
	for range len(record.TestResults) - 1 {
		f.DuplicateRow("Sheet1", startTestRow)
	}

	totalPrice := 0
	for i, testResult := range record.TestResults {
		row := startTestRow + i
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1)

		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), testResult.Name)
		f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), testNameStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), 1)
		f.SetCellStyle("Sheet1", fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), testResultStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), FormatPrice(testResult.Price))
		f.SetCellStyle("Sheet1", fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), testResultStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), FormatPrice(testResult.Price))
		f.SetCellStyle("Sheet1", fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), testResultStyle)

		// Set row height for better spacing (in points, default is usually ~15)
		f.SetRowHeight("Sheet1", row, 19.0)

		totalPrice += testResult.Price * 1
	}

	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", startTestRow+len(record.TestResults)), FormatPrice(totalPrice))

	// Calculate print area based on content (A1 to E + last row with data)
	lastRow := startTestRow + len(record.TestResults) + 2 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$E$%d", lastRow)

	if err := billingReport.ApplyPageSetup(ctx, f, "Sheet1", printArea); err != nil {
		return "", err
	}

	return billingReport.Save(ctx, f)
}

func CreateRecordResultFile(ctx context.Context, record *models.Record) (string, error) {
	resultReport, err := NewReportWithSingleRecord(models.ResultsReport, *record)
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to create result report config", zap.Error(err))
		return "", err
	}
	f, err := resultReport.Open(ctx)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Create style manager
	styleManager := NewStyleManager(ctx, f)

	// Get styles from style manager
	patientNameStyle, err := styleManager.GetPatientNameStyle()
	if err != nil {
		return "", err
	}

	patientInfoStyle, err := styleManager.GetPatientInfoStyle()
	if err != nil {
		return "", err
	}

	testResultStyle, err := styleManager.GetTestResultStyle()
	if err != nil {
		return "", err
	}

	testNameStyle, err := styleManager.GetTestNameStyle()
	if err != nil {
		return "", err
	}

	abnormalStyle, err := styleManager.GetAbnormalStyle()
	if err != nil {
		return "", err
	}

	now := time.Now()
	f.SetCellValue("Sheet1", "C2", now.Format("02/01/2006"))
	f.SetCellStyle("Sheet1", "C2", "C2", patientInfoStyle)

	f.SetCellValue("Sheet1", "C3", record.Patient.Name)
	f.SetCellStyle("Sheet1", "C3", "C3", patientNameStyle)

	f.SetCellValue("Sheet1", "C4", record.Patient.Address)
	f.SetCellStyle("Sheet1", "C4", "C4", patientInfoStyle)

	f.SetCellValue("Sheet1", "E2", record.Patient.Phone)
	f.SetCellStyle("Sheet1", "E2", "E2", patientInfoStyle)

	f.SetCellValue("Sheet1", "E3", record.Patient.YOB)
	f.SetCellStyle("Sheet1", "E3", "E3", patientInfoStyle)

	f.SetCellValue("Sheet1", "E4", record.Patient.Gender)
	f.SetCellStyle("Sheet1", "E4", "E4", patientInfoStyle)

	startTestRow := 8
	for range len(record.TestResults) - 1 {
		f.DuplicateRow("Sheet1", startTestRow)
	}

	for i, testResult := range record.TestResults {
		row := startTestRow + i

		testFieldValue := FormatResult(testResult.Result)
		if testResult.ResultText != "" {
			testFieldValue += testResult.ResultText
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), i+1)
		f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), testResultStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), testResult.Name)
		f.SetCellStyle("Sheet1", fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), testNameStyle)

		resultCell := fmt.Sprintf("D%d", row)
		f.SetCellValue("Sheet1", resultCell, testFieldValue)

		// Apply bold and underline style if result is abnormal, otherwise normal style
		// Manual override has higher priority than automatic detection
		if testResult.Abnormal {
			f.SetCellStyle("Sheet1", resultCell, resultCell, abnormalStyle)
		} else {
			f.SetCellStyle("Sheet1", resultCell, resultCell, testResultStyle)
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testResult.Unit)
		f.SetCellStyle("Sheet1", fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), testResultStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), testResult.NormalValue)
		f.SetCellStyle("Sheet1", fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), testResultStyle)

		// Set row height for better spacing (in points, default is usually ~15)
		f.SetRowHeight("Sheet1", row, 19.0)
	}

	// Calculate print area based on content (A1 to F + last row with data)
	lastRow := startTestRow + len(record.TestResults) + 2 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$F$%d", lastRow+7)

	if err := resultReport.ApplyPageSetup(ctx, f, "Sheet1", printArea); err != nil {
		return "", err
	}

	return resultReport.Save(ctx, f)
}

func CreateRecordResultPDF(ctx context.Context, record *models.Record) (string, error) {
	return createRecordResultFile(ctx, record, models.ResultsWithSignaturePDF, "ket-qua-online.xlsx")
}

func CreateRecordResultWithSignatureFile(ctx context.Context, record *models.Record) (string, error) {
	return createRecordResultFile(ctx, record, models.ResultsWithSignature, "ket-qua-online.xlsx")
}

// createRecordResultFile is a common helper function for creating result files with different templates
func createRecordResultFile(ctx context.Context, record *models.Record, templateType models.ReportType, filenameSuffix string) (string, error) {
	resultReport, err := NewReportWithSingleRecord(templateType, *record)
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to create result report config", zap.Error(err), zap.String("templateType", string(templateType)))
		return "", err
	}
	f, err := resultReport.Open(ctx)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Create style manager
	styleManager := NewStyleManager(ctx, f)

	// Get styles from style manager
	patientNameStyle, err := styleManager.GetPatientNameStyle()
	if err != nil {
		return "", err
	}

	patientInfoStyle, err := styleManager.GetPatientInfoStyle()
	if err != nil {
		return "", err
	}

	testResultStyle, err := styleManager.GetTestResultStyle()
	if err != nil {
		return "", err
	}

	testNameStyle, err := styleManager.GetTestNameStyle()
	if err != nil {
		return "", err
	}

	abnormalStyle, err := styleManager.GetAbnormalStyle()
	if err != nil {
		return "", err
	}

	now := time.Now()
	f.SetCellValue("Sheet1", "C9", now.Format("02/01/2006"))
	f.SetCellStyle("Sheet1", "C9", "C9", patientInfoStyle)

	f.SetCellValue("Sheet1", "C10", record.Patient.Name)
	f.SetCellStyle("Sheet1", "C10", "C10", patientNameStyle)

	f.SetCellValue("Sheet1", "C11", record.Patient.Address)
	f.SetCellStyle("Sheet1", "C11", "C11", patientInfoStyle)

	f.SetCellValue("Sheet1", "E9", record.Patient.Phone)
	f.SetCellStyle("Sheet1", "E9", "E9", patientInfoStyle)

	f.SetCellValue("Sheet1", "E10", record.Patient.YOB)
	f.SetCellStyle("Sheet1", "E10", "E10", patientInfoStyle)

	f.SetCellValue("Sheet1", "E11", record.Patient.Gender)
	f.SetCellStyle("Sheet1", "E11", "E11", patientInfoStyle)

	startTestRow := 15
	for range len(record.TestResults) - 1 {
		f.DuplicateRow("Sheet1", startTestRow)
	}

	for i, testResult := range record.TestResults {
		row := startTestRow + i

		testFieldValue := FormatResult(testResult.Result)
		if testResult.ResultText != "" {
			testFieldValue += testResult.ResultText
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), i+1)
		f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), testResultStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), testResult.Name)
		f.SetCellStyle("Sheet1", fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), testNameStyle)

		resultCell := fmt.Sprintf("D%d", row)
		f.SetCellValue("Sheet1", resultCell, testFieldValue)

		// Apply bold and underline style if result is abnormal, otherwise normal style
		// Manual override has higher priority than automatic detection
		if testResult.Abnormal {
			f.SetCellStyle("Sheet1", resultCell, resultCell, abnormalStyle)
		} else {
			f.SetCellStyle("Sheet1", resultCell, resultCell, testResultStyle)
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testResult.Unit)
		f.SetCellStyle("Sheet1", fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), testResultStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), testResult.NormalValue)
		f.SetCellStyle("Sheet1", fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), testResultStyle)

		// Set row height for better spacing (in points, default is usually ~15)
		f.SetRowHeight("Sheet1", row, 19.0)
	}

	// Calculate print area based on content (A1 to F + last row with data)
	lastRow := startTestRow + len(record.TestResults) + 2 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$F$%d", lastRow+7)     // Add +7 to the last row to handle the signature image

	if err := resultReport.ApplyPageSetup(ctx, f, "Sheet1", printArea); err != nil {
		return "", err
	}

	return resultReport.Save(ctx, f)
}

func CreateRecordTrackingFile(ctx context.Context, records []*models.Record, testMap map[string]models.TestInfo) (string, error) {
	trackingReport, err := NewReportWithMultipleRecords(models.TrackingReport, records)
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to create tracking report config", zap.Error(err))
		return "", err
	}
	f, err := trackingReport.Open(ctx)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Create style manager
	styleManager := NewStyleManager(ctx, f)

	// Get styles from style manager
	patientNameStyle, err := styleManager.GetPatientNameStyle()
	if err != nil {
		return "", err
	}

	patientInfoStyle, err := styleManager.GetPatientInfoStyle()
	if err != nil {
		return "", err
	}

	testResultStyle, err := styleManager.GetTestResultStyle()
	if err != nil {
		return "", err
	}

	testNameStyle, err := styleManager.GetTestNameStyle()
	if err != nil {
		return "", err
	}

	abnormalStyle, err := styleManager.GetAbnormalStyle()
	if err != nil {
		return "", err
	}

	startDate := records[0].CreatedAt.Format("02/01/2006")
	f.SetCellValue("Sheet1", "A4", fmt.Sprintf(" Từ ngày: %s đến ngày: %s", startDate, time.Now().Format("02/01/2006")))

	f.SetCellValue("Sheet1", "A5", fmt.Sprintf("Họ & Tên: %s", records[0].Patient.Name))

	// Apply styles
	f.SetCellStyle("Sheet1", "A4", "A4", patientInfoStyle)
	f.SetCellStyle("Sheet1", "A5", "A5", patientNameStyle)

	startTestRow := 7
	startRecordCol := 'D'

	// If testMap is empty or incomplete, build it from all test results in records
	if len(testMap) == 0 {
		testMap = make(map[string]models.TestInfo)
		for _, record := range records {
			for _, test := range record.TestResults {
				if _, exists := testMap[test.Name]; !exists {
					testMap[test.Name] = models.TestInfo{
						Name:        test.Name,
						NormalValue: test.NormalValue,
						Unit:        test.Unit,
					}
				}
			}
		}
	}

	rowMap := make(map[string]int)
	i := 0
	for testName, testInfo := range testMap {
		rowMap[testName] = startTestRow + i
		f.DuplicateRow("Sheet1", startTestRow+i)
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", startTestRow+i), i+1)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", startTestRow+i), testName)
		f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", startTestRow+i), fmt.Sprintf("B%d", startTestRow+i), testNameStyle)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", startTestRow+i), testInfo.NormalValue)
		i += 1
	}
	// Remove the last duplicated row since we don't need it
	f.RemoveRow("Sheet1", startTestRow+i)

	// Sort records by CreatedAt in increasing order
	slices.SortFunc(records, func(a, b *models.Record) int {
		if a.CreatedAt.Before(b.CreatedAt) {
			return -1
		} else if a.CreatedAt.After(b.CreatedAt) {
			return 1
		}
		return 0
	})

	tableHeaderStyle, err := f.GetCellStyle("Sheet1", "A6")
	if err != nil {
		return "", err
	}

	tableCellStyle, _ := f.GetCellStyle("Sheet1", "A7")
	for j, record := range records {
		col := string(startRecordCol + rune(j))
		headerCell := fmt.Sprintf("%s6", col)
		if j != len(records)-1 {
			f.InsertCols("Sheet1", col, 1)
		}
		f.SetCellValue("Sheet1", headerCell, "Ngày "+record.CreatedAt.Format("02/01/2006"))
		f.SetCellStyle("Sheet1", headerCell, headerCell, tableHeaderStyle)

		// Apply table cell style to the entire column first
		if len(testMap) > 0 {
			f.SetCellStyle("Sheet1", fmt.Sprintf("%s%d", col, startTestRow), fmt.Sprintf("%s%d", col, startTestRow+len(testMap)-1), tableCellStyle)
		}

		// Then set values and apply abnormal styles (which will override table style where needed)
		for _, testResult := range record.TestResults {
			row, exists := rowMap[testResult.Name]
			if !exists {
				// Skip test results that don't exist in testMap
				// This can happen if a test was removed from the combo or is no longer tracked
				continue
			}
			cell := fmt.Sprintf("%s%d", col, row)
			f.SetCellValue("Sheet1", cell, testResult.Result)

			// Apply abnormal style if result is abnormal, otherwise normal style
			// Manual override has higher priority than automatic detection
			if testResult.Abnormal {
				f.SetCellStyle("Sheet1", cell, cell, abnormalStyle)
			} else {
				f.SetCellStyle("Sheet1", cell, cell, testResultStyle)
			}
		}
	}

	// Calculate print area for tracking report (A1 to last column + last row with data)
	lastCol := string(rune('C') + rune(len(records))) // C + number of records
	lastRow := startTestRow + len(testMap) + 1        // Add buffer row
	printArea := fmt.Sprintf("$A$1:$%s$%d", lastCol, lastRow)

	if err := trackingReport.ApplyPageSetup(ctx, f, "Sheet1", printArea); err != nil {
		return "", err
	}

	return trackingReport.Save(ctx, f)
}
