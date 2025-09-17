package sheets

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// CreateRecordTrackingFile generates a tracking report comparing multiple records over time
func CreateRecordTrackingFile(ctx context.Context, records []*models.Record, testList []models.TestInfo) (string, error) {
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

	// Get all common styles at once
	styles, err := styleManager.GetCommonStyles()
	if err != nil {
		return "", err
	}

	startDate := records[0].CreatedAt.Format("02/01/2006")
	f.SetCellValue("Sheet1", "A4", fmt.Sprintf(" Từ: %s đến: %s", startDate, time.Now().Format("02/01/2006")))

	f.SetCellValue("Sheet1", "A5", fmt.Sprintf("Họ & Tên: %s", records[0].Patient.Name))

	// Apply styles
	f.SetCellStyle("Sheet1", "A4", "A4", styles.PatientInfo)

	startTestRow := 7
	startRecordCol := 'D'

	// Set up test rows and get mapping of test names to row numbers
	testNameToRowMap := setupTestRows(f, testList, startTestRow, styles.TestName)

	// Get table styles for record columns
	tableHeaderStyle, err := f.GetCellStyle("Sheet1", "A6")
	if err != nil {
		return "", err
	}
	tableCellStyle, _ := f.GetCellStyle("Sheet1", "A7")

	// Populate record columns with test results
	err = populateRecordColumns(f, records, testNameToRowMap, startRecordCol, startTestRow, testList,
		tableHeaderStyle, tableCellStyle, styles.TestResult, styles.Abnormal)
	if err != nil {
		return "", err
	}

	f.SetCellStyle("Sheet1", "A5", "A5", styles.PatientNameLargeCenter)

	// Calculate print area for tracking report (A1 to last column + last row with data)
	lastCol := string(rune('C') + rune(len(records))) // C + number of records
	lastRow := startTestRow + len(testList) + 1       // Add buffer row
	printArea := fmt.Sprintf("$A$1:$%s$%d", lastCol, lastRow)

	if err := trackingReport.ApplyPageSetup(ctx, f, "Sheet1", printArea); err != nil {
		return "", err
	}

	return trackingReport.Save(ctx, f)
}

// setupTestRows creates rows in the Excel sheet for each test and returns a mapping of test names to row numbers.
// Each test gets its own row with test name, normal value, and sequential numbering.
func setupTestRows(f *excelize.File, testList []models.TestInfo, startTestRow int, testNameStyle int) map[string]int {
	testNameToRowMap := make(map[string]int)

	for i, testInfo := range testList {
		currentRow := startTestRow + i
		testNameToRowMap[testInfo.Name] = currentRow

		// Duplicate the template row for this test
		f.DuplicateRow("Sheet1", currentRow)

		// Fill in test information
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", currentRow), i+1)                  // Sequential number
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", currentRow), testInfo.Name)        // Test name
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", currentRow), testInfo.NormalValue) // Normal value

		// Apply styling to test name
		f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", currentRow), fmt.Sprintf("B%d", currentRow), testNameStyle)
	}

	// Clean up: remove the extra duplicated row
	f.RemoveRow("Sheet1", startTestRow+len(testList))

	return testNameToRowMap
}

// populateRecordColumns fills in the test results for each record in separate columns.
// Records are sorted chronologically and each gets its own column with date header.
func populateRecordColumns(f *excelize.File, records []*models.Record, testNameToRowMap map[string]int,
	startRecordCol rune, startTestRow int, testList []models.TestInfo,
	tableHeaderStyle, tableCellStyle, testResultStyle, abnormalStyle int) error {

	// Sort records chronologically (oldest first)
	slices.SortFunc(records, func(a, b *models.Record) int {
		if a.CreatedAt.Before(b.CreatedAt) {
			return -1
		} else if a.CreatedAt.After(b.CreatedAt) {
			return 1
		}
		return 0
	})

	// Process each record in its own column
	for recordIndex, record := range records {
		columnLetter := string(startRecordCol + rune(recordIndex))
		headerCell := fmt.Sprintf("%s6", columnLetter)

		// Insert new column for all records except the last (template already has one column)
		if recordIndex != len(records)-1 {
			f.InsertCols("Sheet1", columnLetter, 1)
		}

		// Set column header with formatted date
		dateHeader := record.CreatedAt.Format("02/01/2006")
		f.SetCellValue("Sheet1", headerCell, dateHeader)
		f.SetCellStyle("Sheet1", headerCell, headerCell, tableHeaderStyle)

		// Apply base table styling to entire column
		if len(testList) > 0 {
			firstDataCell := fmt.Sprintf("%s%d", columnLetter, startTestRow)
			lastDataCell := fmt.Sprintf("%s%d", columnLetter, startTestRow+len(testList)-1)
			f.SetCellStyle("Sheet1", firstDataCell, lastDataCell, tableCellStyle)
		}

		// Fill in individual test results
		fillTestResultsForRecord(f, record, testNameToRowMap, columnLetter, testResultStyle, abnormalStyle)
	}

	return nil
}

// fillTestResultsForRecord fills in test results for a single record in the specified column
func fillTestResultsForRecord(f *excelize.File, record *models.Record,
	testNameToRowMap map[string]int, columnLetter string, testResultStyle, abnormalStyle int) {

	for _, testResult := range record.TestResults {
		rowNumber, testExists := testNameToRowMap[testResult.Name]
		if !testExists {
			// Skip tests that aren't in our tracking template
			// This can happen if tests were added/removed from combos
			continue
		}

		cellAddress := fmt.Sprintf("%s%d", columnLetter, rowNumber)

		// Set the result value (prefer text over numeric result)
		resultValue := testResult.Result
		if testResult.ResultText != "" {
			resultValue = testResult.ResultText
		}
		f.SetCellValue("Sheet1", cellAddress, resultValue)

		// Apply styling based on abnormal status
		styleToApply := testResultStyle
		if testResult.Abnormal {
			styleToApply = abnormalStyle
		}
		f.SetCellStyle("Sheet1", cellAddress, cellAddress, styleToApply)
	}
}
