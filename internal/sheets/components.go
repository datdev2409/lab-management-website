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
	file          *excelize.File
	styleManager  *StyleManager
	startRow      int
	startCol      string                 // Column letter (e.g., "B")
	testResults   []models.TestResult    // Accept value type directly
	endRow        int                    // Will be set after Apply() is called
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
