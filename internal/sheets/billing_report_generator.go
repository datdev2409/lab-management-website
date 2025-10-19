package sheets

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/xuri/excelize/v2"
)

type BillingReport struct {
	*PageSetup
	*ReportFile
}

func NewBillingReport(ctx context.Context) (*BillingReport, error) {
	report := &BillingReport{
		ReportFile: &ReportFile{
			File: nil,
		},
		PageSetup: &PageSetup{
			SheetName:   "Sheet1",
			PageSize:    9,
			Orientation: "portrait",
			Margins: MarginConfig{
				Top:    float64(0),
				Bottom: 0.511811023622047,
				Left:   0.4,
				Right:  0.4,
				Header: 0.236220472440945,
				Footer: 0.511811023622047,
			},
			ColumnWidth: map[string]float64{
				"A": 7.0,
				"B": 38.0,
				"C": 12.0,
				"D": 15.0,
				"E": 15.0,
			},
		},
	}

	report.File = excelize.NewFile()

	err := report.ApplyColumnWidths(ctx, report.File)
	if err != nil {
		return nil, err
	}

	err = report.ApplyPageSetupV2(ctx, report.File)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (b *BillingReport) Generate(ctx context.Context, data interface{}) (io.Reader, error) {
	record, ok := data.(*models.Record)
	if !ok {
		logger.FromCtx(ctx).Error("Invalid data type for billing report generation")
		return nil, fmt.Errorf("invalid data type for billing report generation")
	}

	f := b.File
	defer f.Close()

	// Create style manager
	sm := NewStyleManager(ctx, f)

	now := time.Now()

	_ = f.MergeCell("Sheet1", "B1", "E1")
	_ = f.MergeCell("Sheet1", "B2", "E2")
	_ = f.MergeCell("Sheet1", "B3", "E3")
	_ = f.MergeCell("Sheet1", "B4", "E4")

	rowHeight := map[int]float64{
		1: 23,
		2: 16,
		3: 27,
		4: 16,
		9: 21,
	}

	for row, height := range rowHeight {
		err := f.SetRowHeight("Sheet1", row, height)
		if err != nil {
			return nil, err
		}
	}

	err := f.AddPicture("Sheet1", "A1", "./assets/anhquanlab_logo.png", &excelize.GraphicOptions{
		ScaleX:          0.7,
		ScaleY:          0.6,
		LockAspectRatio: true,
	})
	if err != nil {
		return nil, err
	}

	cells := map[string]Cell{
		"B1": {
			value:     "PHÒNG XÉT NGHIỆM Y KHOA ANH QUÂN",
			styleName: GetStyleNamePtr(LabNameStyle),
		},
		"B2": {
			value:     "60 Đống Đa, Phường Cao Lãnh, Đồng Tháp",
			styleName: GetStyleNamePtr(LabAddressStyle),
		},
		"B3": {
			value:     "PHIẾU THU",
			styleName: GetStyleNamePtr(ReportNameStyle),
		},
		"B4": {
			value:     fmt.Sprintf("Ngày: %s", now.Format("02/01/2006")),
			styleName: GetStyleNamePtr(ReportDateStyle),
		},
		"A6": {
			value:     "Họ tên:",
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
		"B6": {
			value:     record.Patient.Name,
			styleName: GetStyleNamePtr(PatientNameStyle),
		},
		"A7": {
			value:     "Địa chỉ",
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
		"B7": {
			value:     record.Patient.Address,
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
		"C6": {
			value:     "Số điện thoại:",
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
		"D6": {
			value:     record.Patient.Phone,
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
		"C7": {
			value:     "Năm sinh:",
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
		"D7": {
			value:     record.Patient.YOB,
			styleName: GetStyleNamePtr(PatientInfoStyle),
		},
	}

	for cell, config := range cells {
		if err := f.SetCellValue("Sheet1", cell, config.value); err != nil {
			return nil, err
		}
		if config.styleName != nil {
			err := f.SetCellStyle("Sheet1", cell, cell, sm.GetStyleV2(*config.styleName))
			if err != nil {
				return nil, err
			}
		}
	}

	tableHeaderRow := 9
	tableHeaderStartCell := fmt.Sprintf("A%d", tableHeaderRow)
	tableHeaderEndCell := fmt.Sprintf("E%d", tableHeaderRow)
	tableHeaders := []string{"STT", "Tên hàng hóa, dịch vụ", "Số lượng", "Đơn giá", "Thành tiền"}
	err = f.SetSheetRow("Sheet1", tableHeaderStartCell, &tableHeaders)
	if err != nil {
		return nil, err
	}
	err = f.SetCellStyle("Sheet1", tableHeaderStartCell, tableHeaderEndCell, sm.GetStyleV2(TestTableHeaderStyle))
	if err != nil {
		return nil, err
	}

	startTestRow := 10
	totalPrice := 0
	for i, testResult := range record.TestResults {
		row := startTestRow + i
		rowData := []interface{}{i + 1, testResult.Name, 1, FormatPrice(testResult.Price), FormatPrice(testResult.Price)}
		rowStartCell := fmt.Sprintf("A%d", row)
		err := f.SetSheetRow("Sheet1", rowStartCell, &rowData)
		if err != nil {
			return nil, err
		}

		f.SetCellStyle("Sheet1", rowStartCell, rowStartCell, sm.GetStyleV2(TestIndexStyle))

		testNameCell := fmt.Sprintf("B%d", row)
		f.SetCellStyle("Sheet1", testNameCell, testNameCell, sm.GetStyleV2(TestNameStyle))

		testQuantityCell := fmt.Sprintf("C%d", row)
		f.SetCellStyle("Sheet1", testQuantityCell, testQuantityCell, sm.GetStyleV2(TestQuantityStyle))

		testPriceCell := fmt.Sprintf("D%d", row)
		testPriceTotalCell := fmt.Sprintf("E%d", row)
		f.SetCellStyle("Sheet1", testPriceCell, testPriceTotalCell, sm.GetStyleV2(TestPriceStyle))

		err = f.SetRowHeight("Sheet1", row, 19.0)
		if err != nil {
			return nil, err
		}

		totalPrice += testResult.Price * 1
	}

	totalRow := startTestRow + len(record.TestResults)
	startTotalCell := fmt.Sprintf("A%d", totalRow)
	endTotalCell := fmt.Sprintf("D%d", totalRow)
	err = f.MergeCell("Sheet1", startTotalCell, endTotalCell)
	if err != nil {
		return nil, err
	}
	f.SetRowHeight("Sheet1", totalRow, 19.0)
	f.SetCellValue("Sheet1", startTotalCell, "Tổng thành tiền")
	f.SetCellStyle("Sheet1", startTotalCell, endTotalCell, sm.GetStyleV2(TotalPriceLabelStyle))

	totalPriceCell := fmt.Sprintf("E%d", totalRow)
	f.SetCellValue("Sheet1", totalPriceCell, FormatPrice(totalPrice))
	f.SetCellStyle("Sheet1", totalPriceCell, totalPriceCell, sm.GetStyleV2(TotalPriceStyle))

	// Calculate print area based on content (A1 to E + last row with data)
	lastRow := startTestRow + len(record.TestResults) + 2 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$E$%d", lastRow)

	err = b.ApplyPrintArea(ctx, b.File, printArea)
	if err != nil {
		return nil, err
	}

	return b.GetIOReader(ctx)
}
