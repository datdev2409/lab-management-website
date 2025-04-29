package sheets

import (
	"fmt"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/xuri/excelize/v2"
)

func CreateRecordBillingFile(record models.RecordWithDetails) (string, error) {
	f, err := OpenTemplate("record_billing")
	if err != nil {
		return "", err
	}
	defer f.Close()

	borderStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	f.SetCellValue("Sheet1", "B6", record.Patient.Name)
	f.SetCellValue("Sheet1", "B7", record.Patient.Address)
	f.SetCellValue("Sheet1", "D6", record.Patient.YOB)

	startTestRow := 10
	for i, testResult := range record.TestResults {
		row := startTestRow + i
		testInfo := record.TestInfoMap[testResult.Test.ID]
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), testInfo.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), 1)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), testInfo.Price)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testInfo.Price)
	}

	endTestRow := startTestRow + len(record.TestResults) - 1
	startTestCell := fmt.Sprintf("A%d", startTestRow)
	endTestCell := fmt.Sprintf("E%d", endTestRow)
	f.SetCellStyle("Sheet1", startTestCell, endTestCell, borderStyle)

	now := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("reports/%s-billing-%s.xlsx", record.ID, now)
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}
	return filename, nil
}

func CreateRecordResultFile(record models.RecordWithDetails) (string, error) {
	f, err := OpenTemplate("record_result")
	if err != nil {
		return "", err
	}
	defer f.Close()

	borderStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	f.SetCellValue("Sheet1", "B6", record.Patient.Name)
	f.SetCellValue("Sheet1", "B7", record.Patient.Address)
	f.SetCellValue("Sheet1", "D6", record.Patient.YOB)

	startTestRow := 11
	for i, testResult := range record.TestResults {
		row := startTestRow + i
		testInfo := record.TestInfoMap[testResult.Test.ID]

		testFieldValue := testResult.Result
		if testResult.ResultText != "" {
			testFieldValue += " (" + testResult.ResultText + ")"
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), testInfo.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), testFieldValue)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), testInfo.Unit)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testInfo.NormalValue)
	}

	endTestRow := startTestRow + len(record.TestResults) - 1
	startTestCell := fmt.Sprintf("A%d", startTestRow)
	endTestCell := fmt.Sprintf("E%d", endTestRow)
	f.SetCellStyle("Sheet1", startTestCell, endTestCell, borderStyle)

	now := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("reports/%s-result-%s.xlsx", record.ID, now)
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}
	return filename, nil
}
