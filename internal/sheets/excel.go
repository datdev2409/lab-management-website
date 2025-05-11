package sheets

import (
	"fmt"
	"strings"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/xuri/excelize/v2"
)

func CreateRecordBillingFile(record *models.Record) (string, error) {
	f, err := OpenTemplate("record_billing")
	if err != nil {
		return "", err
	}
	defer f.Close()

	borderConfig := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	leftAlignConfig := &excelize.Alignment{
		Horizontal: "left",
		Vertical:   "center",
		Indent:     1,
	}

	centerAlignConfig := &excelize.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	}

	borderCenterStyle, _ := f.NewStyle(&excelize.Style{
		Border:    borderConfig,
		Alignment: centerAlignConfig,
	})

	borderLeftStyle, _ := f.NewStyle(&excelize.Style{
		Border:    borderConfig,
		Alignment: leftAlignConfig,
	})

	now := time.Now()
	f.SetCellValue("Sheet1", "A4", fmt.Sprintf("Ngày: %s", now.Format("02/01/2006")))
	f.SetCellValue("Sheet1", "B6", record.Patient.Name)
	f.SetCellValue("Sheet1", "B7", record.Patient.Address)
	f.SetCellValue("Sheet1", "D6", record.Patient.YOB)

	startTestRow := 10
	err = f.InsertRows("Sheet1", startTestRow, len(record.TestResults))
	if err != nil {
		return "", err
	}

	totalPrice := 0
	for i, testResult := range record.TestResults {
		row := startTestRow + i
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1)
		f.SetCellStyle("Sheet1", fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), borderCenterStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), testResult.Name)
		f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), borderLeftStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), 1)
		f.SetCellStyle("Sheet1", fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), borderCenterStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), testResult.Price)
		f.SetCellStyle("Sheet1", fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), borderCenterStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testResult.Price)
		f.SetCellStyle("Sheet1", fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), borderCenterStyle)

		totalPrice += testResult.Price * 1
	}

	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", startTestRow+len(record.TestResults)+1), totalPrice)

	pageSizeA4 := 9
	fitToWidth := 1
	f.SetPageLayout("Sheet1", &excelize.PageLayoutOptions{
		Size:       &pageSizeA4,
		FitToWidth: &fitToWidth,
	})

	filename := fmt.Sprintf("reports/%s-%s-hoa-don.xlsx", now.Format("20060102"), strings.ReplaceAll(record.Patient.Name, " ", "_"))
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}
	return filename, nil
}

func CreateRecordResultFile(record *models.Record) (string, error) {
	f, err := OpenTemplate("record_result")
	if err != nil {
		return "", err
	}
	defer f.Close()

	borderConfig := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	leftAlignConfig := &excelize.Alignment{
		Horizontal: "left",
		Vertical:   "center",
		Indent:     1,
	}

	centerAlignConfig := &excelize.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	}

	borderCenterStyle, _ := f.NewStyle(&excelize.Style{
		Border:    borderConfig,
		Alignment: centerAlignConfig,
	})

	borderLeftStyle, _ := f.NewStyle(&excelize.Style{
		Border:    borderConfig,
		Alignment: leftAlignConfig,
	})

	now := time.Now()
	f.SetCellValue("Sheet1", "A4", fmt.Sprintf("Ngày: %s", now.Format("02/01/2006")))
	f.SetCellValue("Sheet1", "B6", record.Patient.Name)
	f.SetCellValue("Sheet1", "B7", record.Patient.Address)
	f.SetCellValue("Sheet1", "D6", record.Patient.YOB)

	startTestRow := 10
	for i, testResult := range record.TestResults {
		row := startTestRow + i

		testFieldValue := testResult.Result
		if testResult.ResultText != "" {
			testFieldValue += " (" + testResult.ResultText + ")"
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1)
		f.SetCellStyle("Sheet1", fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), borderCenterStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), testResult.Name)
		f.SetCellStyle("Sheet1", fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), borderLeftStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), testFieldValue)
		f.SetCellStyle("Sheet1", fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), borderCenterStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), testResult.Unit)
		f.SetCellStyle("Sheet1", fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), borderCenterStyle)

		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), testResult.NormalValue)
		f.SetCellStyle("Sheet1", fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), borderLeftStyle)
	}

	pageSizeA4 := 9
	fitToWidth := 1
	f.SetPageLayout("Sheet1", &excelize.PageLayoutOptions{
		Size:       &pageSizeA4,
		FitToWidth: &fitToWidth,
	})

	// endTestRow := startTestRow + len(record.TestResults) - 1
	// startTestCell := fmt.Sprintf("A%d", startTestRow)
	// endTestCell := fmt.Sprintf("E%d", endTestRow)
	// f.SetCellStyle("Sheet1", startTestCell, endTestCell, borderStyle)

	filename := fmt.Sprintf("reports/%s-%s-ket-qua.xlsx", now.Format("20060102"), strings.ReplaceAll(record.Patient.Name, " ", "_"))
	if err := f.SaveAs(filename); err != nil {
		return "", err
	}
	return filename, nil
}
