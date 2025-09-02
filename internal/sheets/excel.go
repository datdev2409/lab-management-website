package sheets

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/xuri/excelize/v2"
)

// MarginConfig holds margin settings for different templates
type MarginConfig struct {
	Top    float64
	Bottom float64
	Left   float64
	Right  float64
	Header float64
	Footer float64
}

// getTemplateMargins returns specific margin settings for each template
func getTemplateMargins(reportType models.ReportType) MarginConfig {
	switch reportType {
	case models.BillingReport: // Billing template
		return MarginConfig{
			Top:    float64(0),        // Smaller top margin for billing
			Bottom: 0.511811023622047, // Smaller bottom margin
			Left:   0.236220472440945, // Narrow left margin
			Right:  float64(0),        // Narrow right margin
			Header: 0.236220472440945, // Minimal header
			Footer: 0.511811023622047, // Minimal footer
		}
	case models.ResultsReport: // Result template
		return MarginConfig{
			Top:    1.8503937007874,   // Standard top margin
			Bottom: float64(0),        // Standard bottom margin
			Left:   0.118110236220472, // Medium left margin
			Right:  float64(0),        // Medium right margin
			Header: 0.236220472440945, // Standard header
			Footer: float64(0),        // Standard footer
		}
	case models.ResultsWithSignature: // Result with signature template
		return MarginConfig{
			Top:    0.31496063,        // Larger top margin for signature space
			Bottom: float64(0),        // Larger bottom margin for signature space
			Left:   0.118110236220472, // Wider left margin
			Right:  float64(0),        // Wider right margin
			Header: 0.511811023622047, // Standard header
			Footer: float64(0),        // Larger footer for signature
		}
	case models.TrackingReport: // Tracking template
		return MarginConfig{
			Top:    0.25, // Medium top margin
			Bottom: 0.5,  // Medium bottom margin
			Left:   0.5,  // Very narrow left margin (landscape needs more space)
			Right:  0.5,  // Very narrow right margin
			Header: 0.5,  // Compact header
			Footer: 0.5,  // Compact footer
		}
	default: // Default margins
		return MarginConfig{
			Top:    0.75,
			Bottom: 0.75,
			Left:   0.7,
			Right:  0.7,
			Header: 0.3,
			Footer: 0.3,
		}
	}
}

// setupPageLayoutWithCustomMargins configures the Excel sheet for A4 printing with template-specific margins
func setupPageLayoutWithCustomMargins(f *excelize.File, sheetName string, orientation string, reportType models.ReportType) error {
	// Set page layout for A4 paper size
	size := 9
	err := f.SetPageLayout(sheetName, &excelize.PageLayoutOptions{
		Size:        &size,
		Orientation: &orientation,
	})
	if err != nil {
		return err
	}

	// Get template-specific margins
	margins := getTemplateMargins(reportType)

	// Set page margins (in inches)
	err = f.SetPageMargins(sheetName, &excelize.PageLayoutMarginsOptions{
		Bottom: &margins.Bottom,
		Footer: &margins.Footer,
		Header: &margins.Header,
		Left:   &margins.Left,
		Right:  &margins.Right,
		Top:    &margins.Top,
	})
	if err != nil {
		return err
	}

	// Set print options to fit to page
	fitToPage := true
	err = f.SetSheetProps(sheetName, &excelize.SheetPropsOptions{
		FitToPage: &fitToPage,
	})
	return err
}

// setupPageLayoutWithCustomMarginsAndPrintArea configures the Excel sheet with custom margins and print area
func setupPageLayoutWithCustomMarginsAndPrintArea(f *excelize.File, sheetName string, orientation string, reportType models.ReportType, printArea string) error {
	// Set up page layout with custom margins first
	err := setupPageLayoutWithCustomMargins(f, sheetName, orientation, reportType)
	if err != nil {
		return err
	}

	// Set print area
	err = f.SetDefinedName(&excelize.DefinedName{
		Name:     "_xlnm.Print_Area",
		RefersTo: fmt.Sprintf("%s!%s", sheetName, printArea),
		Scope:    sheetName,
	})
	return err
}

func CreateRecordBillingFile(record *models.Record) (string, error) {
	f, err := OpenTemplate("phieu_thu")
	if err != nil {
		return "", err
	}
	defer f.Close()

	now := time.Now()
	f.SetCellValue("Sheet1", "B4", fmt.Sprintf("Ngày: %s", now.Format("02/01/2006")))
	f.SetCellValue("Sheet1", "B6", record.Patient.Name)
	f.SetCellValue("Sheet1", "B7", record.Patient.Address)
	f.SetCellValue("Sheet1", "D6", record.Patient.YOB)

	startTestRow := 10
	for range len(record.TestResults) - 1 {
		f.DuplicateRow("Sheet1", startTestRow)
	}

	totalPrice := 0
	for i, testResult := range record.TestResults {
		row := startTestRow + i
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1)

		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), testResult.Name)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), 1)

		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), testResult.Price)

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testResult.Price)

		totalPrice += testResult.Price * 1
	}

	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", startTestRow+len(record.TestResults)), totalPrice)

	// Calculate print area based on content (A1 to E + last row with data)
	lastRow := startTestRow + len(record.TestResults) + 2 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$E$%d", lastRow)

	// Setup page layout for A4 printing with template-specific margins and print area
	if err := setupPageLayoutWithCustomMarginsAndPrintArea(f, "Sheet1", "portrait", models.BillingReport, printArea); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("reports/%s-%s-hoa-don.xlsx", now.Format("20060102"), strings.ReplaceAll(record.Patient.Name, " ", "_"))
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}
	return filename, nil
}

func CreateRecordResultFile(record *models.Record) (string, error) {
	f, err := OpenTemplate("phieu_ket_qua")
	if err != nil {
		return "", err
	}
	defer f.Close()

	now := time.Now()
	f.SetCellValue("Sheet1", "C2", fmt.Sprintf("Ngày: %s", now.Format("02/01/2006")))
	f.SetCellValue("Sheet1", "C3", record.Patient.Name)
	f.SetCellValue("Sheet1", "C4", record.Patient.Address)
	f.SetCellValue("Sheet1", "E2", record.Patient.Phone)
	f.SetCellValue("Sheet1", "E3", record.Patient.YOB)
	f.SetCellValue("Sheet1", "E4", record.Patient.Gender)

	// Create style for abnormal results (bold + center + borders)
	abnormalStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", err
	}

	startTestRow := 8
	for range len(record.TestResults) - 1 {
		f.DuplicateRow("Sheet1", startTestRow)
	}

	for i, testResult := range record.TestResults {
		row := startTestRow + i

		testFieldValue := testResult.Result
		if testResult.ResultText != "" {
			testFieldValue += " (" + testResult.ResultText + ")"
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), i+1)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), testResult.Name)

		resultCell := fmt.Sprintf("D%d", row)
		f.SetCellValue("Sheet1", resultCell, testFieldValue)

		// Apply bold and underline style if result is abnormal
		// Manual override has higher priority than automatic detection
		if testResult.Abnormal {
			f.SetCellStyle("Sheet1", resultCell, resultCell, abnormalStyle)
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testResult.Unit)

		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), testResult.NormalValue)
	}

	// Calculate print area based on content (A1 to F + last row with data)
	lastRow := startTestRow + len(record.TestResults) + 2 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$F$%d", lastRow)

	// Setup page layout for A4 printing with template-specific margins and print area
	if err := setupPageLayoutWithCustomMarginsAndPrintArea(f, "Sheet1", "portrait", models.ResultsReport, printArea); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("reports/%s-%s-ket-qua.xlsx", now.Format("20060102"), strings.ReplaceAll(record.Patient.Name, " ", "_"))
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}
	return filename, nil
}

func CreateRecordResultWithSignatureFile(record *models.Record) (string, error) {
	f, err := OpenTemplate("phieu_ket_qua_chu_ky")
	if err != nil {
		return "", err
	}
	defer f.Close()

	now := time.Now()
	f.SetCellValue("Sheet1", "C9", fmt.Sprintf("Ngày: %s", now.Format("02/01/2006")))
	f.SetCellValue("Sheet1", "C10", record.Patient.Name)
	f.SetCellValue("Sheet1", "C11", record.Patient.Address)
	f.SetCellValue("Sheet1", "E9", record.Patient.Phone)
	f.SetCellValue("Sheet1", "E10", record.Patient.YOB)
	f.SetCellValue("Sheet1", "E11", record.Patient.Gender)

	// Create style for abnormal results (bold + center + borders)
	abnormalStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return "", err
	}

	startTestRow := 15
	for range len(record.TestResults) - 1 {
		f.DuplicateRow("Sheet1", startTestRow)
	}

	for i, testResult := range record.TestResults {
		row := startTestRow + i

		testFieldValue := testResult.Result
		if testResult.ResultText != "" {
			testFieldValue += " (" + testResult.ResultText + ")"
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), i+1)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), testResult.Name)

		resultCell := fmt.Sprintf("D%d", row)
		f.SetCellValue("Sheet1", resultCell, testFieldValue)

		// Apply bold and underline style if result is abnormal
		// Manual override has higher priority than automatic detection
		if testResult.Abnormal {
			f.SetCellStyle("Sheet1", resultCell, resultCell, abnormalStyle)
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testResult.Unit)

		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), testResult.NormalValue)
	}

	// Calculate print area based on content (A1 to F + last row with data)
	lastRow := startTestRow + len(record.TestResults) + 2 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$F$%d", lastRow+7)     // Add +7 to the last row to handle the signature image

	// Setup page layout for A4 printing with template-specific margins and print area
	if err := setupPageLayoutWithCustomMarginsAndPrintArea(f, "Sheet1", "portrait", models.ResultsWithSignature, printArea); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("reports/%s-%s-ket-qua-online.xlsx", now.Format("20060102"), strings.ReplaceAll(record.Patient.Name, " ", "_"))
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}
	return filename, nil
}

func CreateRecordTrackingFile(records []*models.Record, testMap map[string]models.TestInfo) (string, error) {
	f, err := OpenTemplate("phieu_theo_doi")
	if err != nil {
		return "", err
	}
	defer f.Close()

	startDate := records[0].CreatedAt.Format("02/01/2006")
	f.SetCellValue("Sheet1", "A4", fmt.Sprintf(" Từ ngày: %s đến ngày: %s", startDate, time.Now().Format("02/01/2006")))
	f.SetCellValue("Sheet1", "A5", fmt.Sprintf("Họ & Tên: %s", records[0].Patient.Name))

	now := time.Now()
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

	// Create style for abnormal results (bold + center + borders)
	abnormalStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
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

			// Apply abnormal style if result is abnormal
			// Manual override has higher priority than automatic detection
			if testResult.Abnormal {
				f.SetCellStyle("Sheet1", cell, cell, abnormalStyle)
			}
		}
	}

	// Calculate print area for tracking report (A1 to last column + last row with data)
	lastCol := string(rune('C') + rune(len(records))) // C + number of records
	lastRow := startTestRow + len(testMap) + 1        // Add buffer row
	printArea := fmt.Sprintf("$A$1:$%s$%d", lastCol, lastRow)

	// Setup page layout for A4 printing (landscape orientation for tracking reports) with template-specific margins and print area
	if err := setupPageLayoutWithCustomMarginsAndPrintArea(f, "Sheet1", "landscape", models.TrackingReport, printArea); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("reports/%s-%s-theo-doi.xlsx", now.Format("20060102"), ToLowerCaseNonAccentVietnamese(strings.ReplaceAll(records[0].Patient.Name, " ", "_")))
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}
	return filename, nil
}
