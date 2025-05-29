package sheets

import (
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
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
	startRecordCol := 'D'

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
		log.Println("Error getting table header style:", err)
		return "", err
	}
	tableCellStyle, _ := f.GetCellStyle("Sheet1", "A7")
	for j, record := range records {
		col := string(startRecordCol + rune(j))
		headerCell := fmt.Sprintf("%s6", col)
		if j != 0 {
			f.InsertCols("Sheet1", col, 1)
		}
		f.SetCellValue("Sheet1", headerCell, "Ngày "+record.CreatedAt.Format("02/01/2006"))
		f.SetCellStyle("Sheet1", headerCell, headerCell, tableHeaderStyle)

		for _, testResult := range record.TestResults {
			row := rowMap[testResult.Name]
			cell := fmt.Sprintf("%s%d", col, row)
			f.SetCellValue("Sheet1", cell, testResult.Result)
		}

		if len(testMap) > 0 {
			f.SetCellStyle("Sheet1", fmt.Sprintf("%s%d", col, startTestRow), fmt.Sprintf("%s%d", col, startTestRow+len(testMap)-1), tableCellStyle)
		}
	}

	filename := fmt.Sprintf("reports/%s-%s-theo-doi.xlsx", now.Format("20060102"), ToLowerCaseNonAccentVietnamese(strings.ReplaceAll(records[0].Patient.Name, " ", "_")))
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}
	return filename, nil
}
