package sheets

import (
	"context"
	"fmt"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.uber.org/zap"
)

// CreateRecordResultFile generates a standard test result report for a single record
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

	// Get all common styles at once
	styles, err := styleManager.GetCommonStyles()
	if err != nil {
		return "", err
	}

	now := time.Now()
	f.SetCellValue("Sheet1", "C2", now.Format("02/01/2006"))
	f.SetCellStyle("Sheet1", "C2", "C2", styles.PatientInfo)

	f.SetCellValue("Sheet1", "C3", record.Patient.Name)
	f.SetCellStyle("Sheet1", "C3", "C3", styles.PatientName)

	f.SetCellValue("Sheet1", "C4", record.Patient.Address)
	f.SetCellStyle("Sheet1", "C4", "C4", styles.PatientInfo)

	f.SetCellValue("Sheet1", "E2", record.Patient.Phone)
	f.SetCellStyle("Sheet1", "E2", "E2", styles.PatientInfo)

	f.SetCellValue("Sheet1", "E3", record.Patient.YOB)
	f.SetCellStyle("Sheet1", "E3", "E3", styles.PatientInfo)

	f.SetCellValue("Sheet1", "E4", record.Patient.Gender)
	f.SetCellStyle("Sheet1", "E4", "E4", styles.PatientInfo)

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
		f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styles.TestResult)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), testResult.Name)
		f.SetCellStyle("Sheet1", fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styles.TestName)

		resultCell := fmt.Sprintf("D%d", row)
		f.SetCellValue("Sheet1", resultCell, testFieldValue)

		// Apply bold and underline style if result is abnormal, otherwise normal style
		// Manual override has higher priority than automatic detection
		if testResult.Abnormal {
			f.SetCellStyle("Sheet1", resultCell, resultCell, styles.Abnormal)
		} else {
			f.SetCellStyle("Sheet1", resultCell, resultCell, styles.TestResult)
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testResult.Unit)
		f.SetCellStyle("Sheet1", fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styles.TestResult)

		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), testResult.NormalValue)
		f.SetCellStyle("Sheet1", fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styles.TestResult)

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

// CreateRecordResultPDF generates a PDF version of test results with signature
func CreateRecordResultPDF(ctx context.Context, record *models.Record) (string, error) {
	return createRecordResultFile(ctx, record, models.ResultsWithSignaturePDF, "ket-qua-online.xlsx")
}

// CreateRecordResultWithSignatureFile generates test results with signature template
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

	// Get all common styles at once
	styles, err := styleManager.GetCommonStyles()
	if err != nil {
		return "", err
	}

	now := time.Now()
	f.SetCellValue("Sheet1", "C9", now.Format("02/01/2006"))
	f.SetCellStyle("Sheet1", "C9", "C9", styles.PatientInfo)

	f.SetCellValue("Sheet1", "C10", record.Patient.Name)
	f.SetCellStyle("Sheet1", "C10", "C10", styles.PatientName)

	f.SetCellValue("Sheet1", "C11", record.Patient.Address)
	f.SetCellStyle("Sheet1", "C11", "C11", styles.PatientInfo)

	f.SetCellValue("Sheet1", "E9", record.Patient.Phone)
	f.SetCellStyle("Sheet1", "E9", "E9", styles.PatientInfo)

	f.SetCellValue("Sheet1", "E10", record.Patient.YOB)
	f.SetCellStyle("Sheet1", "E10", "E10", styles.PatientInfo)

	f.SetCellValue("Sheet1", "E11", record.Patient.Gender)
	f.SetCellStyle("Sheet1", "E11", "E11", styles.PatientInfo)

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
		f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styles.TestResult)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), testResult.Name)
		f.SetCellStyle("Sheet1", fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styles.TestName)

		resultCell := fmt.Sprintf("D%d", row)
		f.SetCellValue("Sheet1", resultCell, testFieldValue)

		// Apply bold and underline style if result is abnormal, otherwise normal style
		// Manual override has higher priority than automatic detection
		if testResult.Abnormal {
			f.SetCellStyle("Sheet1", resultCell, resultCell, styles.Abnormal)
		} else {
			f.SetCellStyle("Sheet1", resultCell, resultCell, styles.TestResult)
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testResult.Unit)
		f.SetCellStyle("Sheet1", fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styles.TestResult)

		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), testResult.NormalValue)
		f.SetCellStyle("Sheet1", fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styles.TestResult)

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
