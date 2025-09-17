package sheets

import (
	"context"
	"fmt"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.uber.org/zap"
)

// CreateRecordBillingFile generates a billing report (invoice) for a single record
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

	// Get all common styles at once
	styles, err := styleManager.GetCommonStyles()
	if err != nil {
		return "", err
	}

	now := time.Now()
	f.SetCellValue("Sheet1", "B4", fmt.Sprintf("Ngày: %s", now.Format("02/01/2006")))
	f.SetCellStyle("Sheet1", "B4", "B4", styles.DateCenter)

	f.SetCellValue("Sheet1", "B6", record.Patient.Name)
	f.SetCellStyle("Sheet1", "B6", "B6", styles.PatientName)

	f.SetCellValue("Sheet1", "B7", record.Patient.Address)
	f.SetCellStyle("Sheet1", "B7", "B7", styles.PatientInfo)

	f.SetCellValue("Sheet1", "D6", record.Patient.YOB)
	f.SetCellStyle("Sheet1", "D6", "D6", styles.PatientInfo)

	startTestRow := 10
	for range len(record.TestResults) - 1 {
		f.DuplicateRow("Sheet1", startTestRow)
	}

	totalPrice := 0
	for i, testResult := range record.TestResults {
		row := startTestRow + i
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1)

		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), testResult.Name)
		f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styles.TestName)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), 1)
		f.SetCellStyle("Sheet1", fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styles.TestResult)

		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), FormatPrice(testResult.Price))
		f.SetCellStyle("Sheet1", fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), styles.PriceRight)

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), FormatPrice(testResult.Price))
		f.SetCellStyle("Sheet1", fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styles.PriceRight)

		// Set row height for better spacing (in points, default is usually ~15)
		f.SetRowHeight("Sheet1", row, 19.0)

		totalPrice += testResult.Price * 1
	}

	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", startTestRow+len(record.TestResults)), FormatPrice(totalPrice))
	f.SetCellStyle("Sheet1", fmt.Sprintf("E%d", startTestRow+len(record.TestResults)), fmt.Sprintf("E%d", startTestRow+len(record.TestResults)), styles.PriceRight)

	// Calculate print area based on content (A1 to E + last row with data)
	lastRow := startTestRow + len(record.TestResults) + 2 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$E$%d", lastRow)

	if err := billingReport.ApplyPageSetup(ctx, f, "Sheet1", printArea); err != nil {
		return "", err
	}

	return billingReport.Save(ctx, f)
}
