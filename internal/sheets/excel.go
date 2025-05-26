package sheets

import (
	"fmt"
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
