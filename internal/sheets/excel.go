package sheets

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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
	for i := range len(record.TestResults) - 1 {
		fmt.Print("Duplicating row for test result...\n", i)
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

	startTestRow := 8
	for i := range len(record.TestResults) - 1 {
		fmt.Print("Duplicating row for test result...\n", i)
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

		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), testFieldValue)

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testResult.Unit)

		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), testResult.NormalValue)
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

	startTestRow := 15
	for i := range len(record.TestResults) - 1 {
		fmt.Print("Duplicating row for test result...\n", i)
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

		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), testFieldValue)

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testResult.Unit)

		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), testResult.NormalValue)
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

	now := time.Now()
	startTestRow := 7

	for i := range len(testMap) - 1 {
		fmt.Print("Duplicating row for test result...\n", i)
		f.DuplicateRow("Sheet1", startTestRow)
	}

	i := 0
	for testName, testInfo := range testMap {
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", startTestRow+i), i+1)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", startTestRow+i), testName)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", startTestRow+i), testInfo.NormalValue)
		i += 1
	}

	// Create columns for each record date
	tableHeaderStyle, err := f.GetCellStyle("Sheet1", "A6")
	var data map[string]map[bson.ObjectID]models.TestResult = make(map[string]map[bson.ObjectID]models.TestResult)
	if err != nil {
		return "", err
	}
	for i, record := range records {
		if i != 0 {
			// Insert a new column for each record after the first one
			f.InsertCols("Sheet1", "D", 1)
		}
		f.SetCellValue("Sheet1", "D6", "Ngày "+record.CreatedAt.Format("02/01/2006"))
		f.SetCellStyle("Sheet1", "D6", "D6", tableHeaderStyle)

		for _, testResult := range record.TestResults {
			if _, exists := data[testResult.Name]; !exists {
				data[testResult.Name] = make(map[bson.ObjectID]models.TestResult)
			}
			data[testResult.Name][record.ID] = testResult
		}
	}

	tableCellStyle, err := f.GetCellStyle("Sheet1", "A7")
	if err != nil {
		return "", err
	}

	row := 0
	slices.Reverse(records)
	for testName := range testMap {
		testData := data[testName]
		for col, record := range records {
			value, exists := testData[record.ID]
			if exists {
				cell := fmt.Sprintf("%s%d", string(rune('D'+col)), startTestRow+row)
				f.SetCellValue("Sheet1", cell, value.Result)
				f.SetCellStyle("Sheet1", cell, cell, tableCellStyle)
			} else {
				cell := fmt.Sprintf("%s%d", string(rune('D'+col)), startTestRow+row)
				f.SetCellValue("Sheet1", cell, "N/A") // or leave it empty
				f.SetCellStyle("Sheet1", cell, cell, tableCellStyle)
			}
		}
		row += 1
	}

	// // Fill dates in header
	// for i, recordID := range recordDates {
	// 	col := string('D' + i)
	// 	f.SetCellValue("Sheet1", fmt.Sprintf("%s7", col), recordID)
	// }

	// // Fill test names and results
	// row := startTestRow
	// i := 1
	// for testID := range testIDs {
	// 	testInfo := testMap[testID]
	// 	f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), i)
	// 	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), testInfo.Name)

	// 	// Fill results for each date
	// 	for j, recordID := range recordDates {
	// 		col := string('D' + j)
	// 		if result, exists := data[recordID][testID]; exists {
	// 			f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", col, row), result.Result)
	// 		}
	// 	}
	// 	row++
	// 	i++
	// }

	filename := fmt.Sprintf("reports/%s-theo-doi.xlsx", now.Format("20060102"))
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}
	return filename, nil
}
