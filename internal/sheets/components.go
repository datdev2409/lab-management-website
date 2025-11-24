package sheets

import (
	"context"
	"fmt"

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

	col1 := t.startCol
	col2 := GetNextColumn(col1)
	col3 := GetNextColumn(col2)
	col4 := GetNextColumn(col3)
	col5 := GetNextColumn(col4)

	// Write the header row
	headerRow := []interface{}{"STT\n(No)", "Tên Xét Nghiệm\n(Test)", "Kết Quả\n(Results)", "Đơn Vị\n(Unit)", "Khoảng Tham Chiếu\n(Reference)"}
	endCol := t.startCol
	for range len(headerRow) - 1 {
		endCol = GetNextColumn(endCol)
	}

	startHeaderCell := fmt.Sprintf("%s%d", t.startCol, t.startRow)
	endHeaderCell := fmt.Sprintf("%s%d", endCol, t.startRow)

	t.file.SetSheetRow("Sheet1", startHeaderCell, &headerRow)
	t.file.SetRowHeight("Sheet1", t.startRow, 32)
	t.file.SetCellStyle("Sheet1", startHeaderCell, endHeaderCell, t.styleManager.GetStyleV2(TestResultTableHeader))

	// Write each test result
	for i, testResult := range t.testResults {
		row := t.startRow + 1 + i

		// Format the test result value
		testFieldValue := FormatResult(testResult.Result)
		if testResult.Result == "" && testResult.ResultText != "" {
			testFieldValue = testResult.ResultText
		}
		// Prepare the data row: [nil for formula, TestName, Result, Unit, NormalValue]
		// We'll set the formula separately using SetCellFormula
		dataRow := []interface{}{nil, testResult.Name, testFieldValue, testResult.Unit, testResult.NormalValue}

		// Write the data row starting from startCol (typically "B")
		startCell := fmt.Sprintf("%s%d", t.startCol, row)
		t.file.SetSheetRow("Sheet1", startCell, &dataRow)

		// Set the auto-increment formula for the index cell
		indexCell := fmt.Sprintf("%s%d", t.startCol, row)
		indexFormula := SetAutoIncrementIndexFormula(t.startRow + 1)
		if err := t.file.SetCellFormula("Sheet1", indexCell, indexFormula); err != nil {
			return err
		}

		// Apply styles to each cell
		// STT (Serial number) - Column B
		t.file.SetCellStyle("Sheet1", fmt.Sprintf("%s%d", col1, row), fmt.Sprintf("%s%d", col1, row), t.styleManager.GetStyleV2(TestIndexStyle))

		// Test Name - Column C
		t.file.SetCellStyle("Sheet1", fmt.Sprintf("%s%d", col2, row), fmt.Sprintf("%s%d", col2, row), t.styleManager.GetStyleV2(TestNameStyle))

		// Result - Column D (apply abnormal style if needed)
		resultCell := fmt.Sprintf("%s%d", col3, row)
		if testResult.Abnormal {
			t.file.SetCellStyle("Sheet1", resultCell, resultCell, t.styleManager.GetStyleV2(TestAbnormalResultStyle))
		} else {
			t.file.SetCellStyle("Sheet1", resultCell, resultCell, t.styleManager.GetStyleV2(TestResultStyle))
		}

		// Unit - Column E
		t.file.SetCellStyle("Sheet1", fmt.Sprintf("%s%d", col4, row), fmt.Sprintf("%s%d", col4, row), t.styleManager.GetStyleV2(TestUnitStyle))

		// Normal Range - Column F
		t.file.SetCellStyle("Sheet1", fmt.Sprintf("%s%d", col5, row), fmt.Sprintf("%s%d", col5, row), t.styleManager.GetStyleV2(TestNormalRangeStyle))

		// Set row height for better spacing
		t.file.SetRowHeight("Sheet1", row, 19.0)
	}

	// Update end row
	t.endRow = t.startRow + len(t.testResults) // Table header + data rows

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
	col1 := p.startCol
	col2 := GetNextColumn(col1)
	col3 := GetNextColumn(col2)
	col4 := GetNextColumn(col3)
	col5 := GetNextColumn(col4)

	sm := p.styleManager
	f := p.file

	// Patient info layout:
	// Row N: ["Họ tên", Patient.Name, "Năm sinh", Patient.YOB]
	// Row N+1: ["Địa chỉ", Patient.Address, "Giới tính", Patient.Gender]
	// Row N+2: ["Chẩn đoán", "", "Số điện thoại", Patient.Phone]

	row := p.startRow

	// Row 1: Name and YOB
	nameCell := fmt.Sprintf("%s%d", col1, row)
	f.SetCellValue("Sheet1", nameCell, "Họ tên:")
	f.SetCellStyle("Sheet1", nameCell, nameCell, sm.GetStyleV2(PatientInfoStyle))

	nameValueCell := fmt.Sprintf("%s%d", col2, row)
	f.SetCellValue("Sheet1", nameValueCell, p.patient.Name)
	f.SetCellStyle("Sheet1", nameValueCell, nameValueCell, sm.GetStyleV2(PatientNameStyle))

	yobLabelCell := fmt.Sprintf("%s%d", col3, row)
	f.SetCellValue("Sheet1", yobLabelCell, "Năm sinh:")
	f.SetCellStyle("Sheet1", yobLabelCell, yobLabelCell, sm.GetStyleV2(PatientInfoStyle))

	yobValueCell := fmt.Sprintf("%s%d", col4, row)
	yobValueEndCell := fmt.Sprintf("%s%d", col5, row)
	f.MergeCell("Sheet1", yobValueCell, yobValueEndCell)
	f.SetCellValue("Sheet1", yobValueCell, p.patient.YOB)
	f.SetCellStyle("Sheet1", yobValueCell, yobValueEndCell, sm.GetStyleV2(PatientInfoStyle))

	row++

	// Row 2: Address and Gender
	addressCell := fmt.Sprintf("%s%d", col1, row)
	f.SetCellValue("Sheet1", addressCell, "Địa chỉ:")
	f.SetCellStyle("Sheet1", addressCell, addressCell, sm.GetStyleV2(PatientInfoStyle))

	addressValueCell := fmt.Sprintf("%s%d", col2, row)
	f.SetCellValue("Sheet1", addressValueCell, p.patient.Address)
	f.SetCellStyle("Sheet1", addressValueCell, addressValueCell, sm.GetStyleV2(PatientInfoStyle))

	genderLabelCell := fmt.Sprintf("%s%d", col3, row)
	f.SetCellValue("Sheet1", genderLabelCell, "Giới tính:")
	f.SetCellStyle("Sheet1", genderLabelCell, genderLabelCell, sm.GetStyleV2(PatientInfoStyle))

	genderValueCell := fmt.Sprintf("%s%d", col4, row)
	genderValueEndCell := fmt.Sprintf("%s%d", col5, row)
	f.MergeCell("Sheet1", genderValueCell, genderValueEndCell)
	f.SetCellValue("Sheet1", genderValueCell, p.patient.Gender)
	f.SetCellStyle("Sheet1", genderValueCell, genderValueEndCell, sm.GetStyleV2(PatientInfoStyle))

	row++

	// Row 3: Diagnosis and Phone
	diagnosisLabelCell := fmt.Sprintf("%s%d", col1, row)
	f.SetCellValue("Sheet1", diagnosisLabelCell, "Chẩn đoán:")
	f.SetCellStyle("Sheet1", diagnosisLabelCell, diagnosisLabelCell, sm.GetStyleV2(PatientInfoStyle))

	diagnosisValueCell := fmt.Sprintf("%s%d", col2, row)
	f.SetCellValue("Sheet1", diagnosisValueCell, "")
	f.SetCellStyle("Sheet1", diagnosisValueCell, diagnosisValueCell, sm.GetStyleV2(PatientInfoStyle))

	phoneCell := fmt.Sprintf("%s%d", col3, row)
	f.SetCellValue("Sheet1", phoneCell, "Số điện thoại:")
	f.SetCellStyle("Sheet1", phoneCell, phoneCell, sm.GetStyleV2(PatientInfoStyle))

	phoneValueCell := fmt.Sprintf("%s%d", col4, row)
	phoneValueEndCell := fmt.Sprintf("%s%d", col5, row)
	f.MergeCell("Sheet1", phoneValueCell, phoneValueEndCell)
	f.SetCellValue("Sheet1", phoneValueCell, p.patient.Phone)
	f.SetCellStyle("Sheet1", phoneValueCell, phoneValueEndCell, sm.GetStyleV2(PatientInfoStyle))

	// Update end row (patient info takes 3 rows)
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
