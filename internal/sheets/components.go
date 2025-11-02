package sheets

import (
	"context"
	"fmt"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/xuri/excelize/v2"
)

// ReportComponent is the interface for all reusable report components
type ReportComponent interface {
	// Apply renders the component onto the Excel sheet
	Apply(ctx context.Context) error
	// GetEndRow returns the last row used by this component
	GetEndRow() int
}

// TestResultTable is a component that renders test results in a table format
type TestResultTable struct {
	file         *excelize.File
	styleManager *StyleManager
	startRow     int
	startCol     string              // Column letter (e.g., "B")
	testResults  []models.TestResult // Accept value type directly
	endRow       int                 // Will be set after Apply() is called
}

// NewTestResultTable creates a new TestResultTable component
func NewTestResultTable(
	file *excelize.File,
	styleManager *StyleManager,
	startRow int,
	startCol string,
	testResults []models.TestResult,
) *TestResultTable {
	return &TestResultTable{
		file:         file,
		styleManager: styleManager,
		startRow:     startRow,
		startCol:     startCol,
		testResults:  testResults,
		endRow:       startRow - 1, // Will be updated after Apply()
	}
}

// Apply renders the test result table to the Excel sheet
func (t *TestResultTable) Apply(ctx context.Context) error {
	// Duplicate rows for all test results except the first one
	if len(t.testResults) > 1 {
		for range len(t.testResults) - 1 {
			t.file.DuplicateRow("Sheet1", t.startRow)
		}
	}

	// Write each test result
	for i, testResult := range t.testResults {
		row := t.startRow + i

		// Format the test result value
		testFieldValue := FormatResult(testResult.Result)
		if testResult.ResultText != "" {
			testFieldValue += testResult.ResultText
		}

		// Prepare the data row: [Index, TestName, Result, Unit, NormalValue]
		dataRow := []interface{}{i + 1, testResult.Name, testFieldValue, testResult.Unit, testResult.NormalValue}

		// Write the data row starting from startCol (typically "B")
		startCell := fmt.Sprintf("%s%d", t.startCol, row)
		t.file.SetSheetRow("Sheet1", startCell, &dataRow)

		// Apply styles to each cell
		// STT (Serial number) - Column B
		t.file.SetCellStyle("Sheet1", fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), t.styleManager.GetStyleV2(TestIndexStyle))

		// Test Name - Column C
		t.file.SetCellStyle("Sheet1", fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), t.styleManager.GetStyleV2(TestNameStyle))

		// Result - Column D (apply abnormal style if needed)
		resultCell := fmt.Sprintf("D%d", row)
		if testResult.Abnormal {
			t.file.SetCellStyle("Sheet1", resultCell, resultCell, t.styleManager.GetStyleV2(TestAbnormalResultStyle))
		} else {
			t.file.SetCellStyle("Sheet1", resultCell, resultCell, t.styleManager.GetStyleV2(TestResultStyle))
		}

		// Unit - Column E
		t.file.SetCellStyle("Sheet1", fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), t.styleManager.GetStyleV2(TestUnitStyle))

		// Normal Range - Column F
		t.file.SetCellStyle("Sheet1", fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), t.styleManager.GetStyleV2(TestNormalRangeStyle))

		// Set row height for better spacing
		t.file.SetRowHeight("Sheet1", row, 19.0)
	}

	// Update end row
	t.endRow = t.startRow + len(t.testResults) - 1

	return nil
}

// GetEndRow returns the last row used by this component
func (t *TestResultTable) GetEndRow() int {
	return t.endRow
}

// PatientInfoTable is a component that renders patient information in a table format
type PatientInfoTable struct {
	file         *excelize.File
	styleManager *StyleManager
	patient      *models.Patient // Accept Patient object, not full Record
	startRow     int
	startCol     string // Column letter (e.g., "B")
	endRow       int    // Will be set after Apply() is called
}

// NewPatientInfoTable creates a new PatientInfoTable component
func NewPatientInfoTable(
	file *excelize.File,
	styleManager *StyleManager,
	patient *models.Patient,
	startRow int,
	startCol string,
) *PatientInfoTable {
	return &PatientInfoTable{
		file:         file,
		styleManager: styleManager,
		patient:      patient,
		startRow:     startRow,
		startCol:     startCol,
		endRow:       startRow - 1, // Will be updated after Apply()
	}
}

// Apply renders the patient info table to the Excel sheet
func (p *PatientInfoTable) Apply(ctx context.Context) error {
	sm := p.styleManager
	f := p.file

	// Patient info layout:
	// Row N: [blank, "Họ tên", Patient.Name, "Ngày khám", Today's date]
	// Row N+1: [blank, "Địa chỉ", Patient.Address, "Số điện thoại", Patient.Phone]
	// Row N+2: [blank, "Năm sinh", Patient.YOB, "Giới tính", Patient.Gender]

	row := p.startRow

	// Row 1: Name and date
	nameCell := fmt.Sprintf("%s%d", p.startCol, row)
	f.SetCellValue("Sheet1", nameCell, "Họ tên")
	f.SetCellStyle("Sheet1", nameCell, nameCell, sm.GetStyleV2(PatientInfoStyle))

	nameValueCell := fmt.Sprintf("%s%d", GetNextColumn(p.startCol), row)
	f.SetCellValue("Sheet1", nameValueCell, p.patient.Name)
	f.SetCellStyle("Sheet1", nameValueCell, nameValueCell, sm.GetStyleV2(PatientNameStyle))

	dateCell := fmt.Sprintf("%s%d", "D", row)
	f.SetCellValue("Sheet1", dateCell, "Ngày khám")
	f.SetCellStyle("Sheet1", dateCell, dateCell, sm.GetStyleV2(PatientInfoStyle))

	dateValueCell := fmt.Sprintf("%s%d", "E", row)
	today := fmt.Sprintf("%02d/%02d/%04d", time.Now().Day(), time.Now().Month(), time.Now().Year())
	f.SetCellValue("Sheet1", dateValueCell, today)
	f.SetCellStyle("Sheet1", dateValueCell, dateValueCell, sm.GetStyleV2(PatientInfoStyle))

	row++

	// Row 2: Address and phone
	addressCell := fmt.Sprintf("%s%d", p.startCol, row)
	f.SetCellValue("Sheet1", addressCell, "Địa chỉ")
	f.SetCellStyle("Sheet1", addressCell, addressCell, sm.GetStyleV2(PatientInfoStyle))

	addressValueCell := fmt.Sprintf("%s%d", GetNextColumn(p.startCol), row)
	f.SetCellValue("Sheet1", addressValueCell, p.patient.Address)
	f.SetCellStyle("Sheet1", addressValueCell, addressValueCell, sm.GetStyleV2(PatientInfoStyle))

	phoneCell := fmt.Sprintf("%s%d", "D", row)
	f.SetCellValue("Sheet1", phoneCell, "Số điện thoại")
	f.SetCellStyle("Sheet1", phoneCell, phoneCell, sm.GetStyleV2(PatientInfoStyle))

	phoneValueCell := fmt.Sprintf("%s%d", "E", row)
	f.SetCellValue("Sheet1", phoneValueCell, p.patient.Phone)
	f.SetCellStyle("Sheet1", phoneValueCell, phoneValueCell, sm.GetStyleV2(PatientInfoStyle))

	row++

	// Row 3: YOB and gender
	yobLabelCell := fmt.Sprintf("%s%d", p.startCol, row)
	f.SetCellValue("Sheet1", yobLabelCell, "Năm sinh")
	f.SetCellStyle("Sheet1", yobLabelCell, yobLabelCell, sm.GetStyleV2(PatientInfoStyle))

	yobValueCell := fmt.Sprintf("%s%d", GetNextColumn(p.startCol), row)
	f.SetCellValue("Sheet1", yobValueCell, p.patient.YOB)
	f.SetCellStyle("Sheet1", yobValueCell, yobValueCell, sm.GetStyleV2(PatientInfoStyle))

	genderLabelCell := fmt.Sprintf("%s%d", "D", row)
	f.SetCellValue("Sheet1", genderLabelCell, "Giới tính")
	f.SetCellStyle("Sheet1", genderLabelCell, genderLabelCell, sm.GetStyleV2(PatientInfoStyle))

	genderValueCell := fmt.Sprintf("%s%d", "E", row)
	f.SetCellValue("Sheet1", genderValueCell, p.patient.Gender)
	f.SetCellStyle("Sheet1", genderValueCell, genderValueCell, sm.GetStyleV2(PatientInfoStyle))

	// Row 4: Diagnosis (spanning columns B to C)
	row++
	diagnosisLabelCell := fmt.Sprintf("%s%d", p.startCol, row)
	f.SetCellValue("Sheet1", diagnosisLabelCell, "Chẩn đoán")
	f.SetCellStyle("Sheet1", diagnosisLabelCell, diagnosisLabelCell, sm.GetStyleV2(PatientInfoStyle))
	f.MergeCell("Sheet1", diagnosisLabelCell, fmt.Sprintf("%s%d", GetNextColumn(p.startCol), row))

	// Update end row (patient info takes 4 rows)
	p.endRow = row

	return nil
}

// GetEndRow returns the last row used by this component
func (p *PatientInfoTable) GetEndRow() int {
	return p.endRow
}

// Helper function to get the next column letter
func GetNextColumn(col string) string {
	if len(col) == 0 {
		return "B"
	}
	// Simple conversion: A->B, B->C, C->D, etc.
	return string([]byte{col[0] + 1})
}
