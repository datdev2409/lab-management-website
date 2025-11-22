package sheets

import (
	"context"
	"fmt"
	"io"

	_ "image/png"

	"github.com/datdev2409/lab-admin-go/internal/models"
)

type BillingReport struct {
	*BaseReportBuilder
}

func NewBillingReport(ctx context.Context) (*BillingReport, error) {
	pageSetup := &PageSetup{
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
	}

	builder, err := NewBaseReportBuilder(ctx, pageSetup)
	if err != nil {
		return nil, err
	}

	if err := builder.InitializeNewFile(ctx); err != nil {
		return nil, err
	}

	return &BillingReport{BaseReportBuilder: builder}, nil
}

func (r *BillingReport) Generate(ctx context.Context, data interface{}) (io.Reader, error) {
	record, ok := data.(*models.Record)
	if !ok {
		return nil, fmt.Errorf("invalid data type for billing report generation")
	}

	f := r.File
	defer f.Close()

	// Create style manager
	sm := NewStyleManager(ctx, f)

	now := GetVietnamTime()

	// Set row heights
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

	// Add logo
	err := AddLogoComponent(ctx, f, "Sheet1", "A1", "./assets/anhquanlab_logo.png", 0.7, 0.6)
	if err != nil {
		return nil, err
	}

	// Lab name, address, and title merged cells
	_ = r.MergeCellsWithStyle("Sheet1", "B1", "E1", "PHÒNG XÉT NGHIỆM Y KHOA ANH QUÂN", sm.GetStyleV2(LabNameStyle))
	_ = r.MergeCellsWithStyle("Sheet1", "B2", "E2", "60 Đống Đa, Phường Cao Lãnh, Đồng Tháp", sm.GetStyleV2(LabAddressStyle))
	_ = r.MergeCellsWithStyle("Sheet1", "B3", "E3", "PHIẾU THU", sm.GetStyleV2(ReportNameStyle))
	_ = r.MergeCellsWithStyle("Sheet1", "B4", "E4", fmt.Sprintf("Ngày: %s", now.Format("02/01/2006")), sm.GetStyleV2(ReportDateStyle))

	// Patient information
	cells := map[string]Cell{
		"A6": {value: "Họ tên:", styleName: GetStyleNamePtr(PatientInfoStyle)},
		"B6": {value: record.Patient.Name, styleName: GetStyleNamePtr(PatientNameStyle)},
		"A7": {value: "Địa chỉ", styleName: GetStyleNamePtr(PatientInfoStyle)},
		"B7": {value: record.Patient.Address, styleName: GetStyleNamePtr(PatientInfoStyle)},
		"C6": {value: "Số điện thoại:", styleName: GetStyleNamePtr(PatientInfoStyle)},
		"D6": {value: record.Patient.Phone, styleName: GetStyleNamePtr(PatientInfoStyle)},
		"C7": {value: "Năm sinh:", styleName: GetStyleNamePtr(PatientInfoStyle)},
		"D7": {value: record.Patient.YOB, styleName: GetStyleNamePtr(PatientInfoStyle)},
	}

	for cell, config := range cells {
		if err := r.SetCellWithStyle("Sheet1", cell, config.value, sm.GetStyleV2(*config.styleName)); err != nil {
			return nil, err
		}
	}

	// Table header
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

	// Test results rows
	startTestRow := 10
	for i, testResult := range record.TestResults {
		row := startTestRow + i
		// Use formula to calculate index: ROW() - startTestRow + 1
		// This ensures auto-increment even if rows are deleted or updated
		// Use numeric values for prices so Excel SUM formula works correctly
		rowData := []interface{}{nil, testResult.Name, 1, testResult.Price, testResult.Price}
		rowStartCell := fmt.Sprintf("A%d", row)
		err := f.SetSheetRow("Sheet1", rowStartCell, &rowData)
		if err != nil {
			return nil, err
		}

		// Set the auto-increment formula for the index cell
		indexCell := fmt.Sprintf("A%d", row)
		indexFormula := SetAutoIncrementIndexFormula(startTestRow)
		err = f.SetCellFormula("Sheet1", indexCell, indexFormula)
		if err != nil {
			return nil, err
		}
		f.SetCellStyle("Sheet1", indexCell, indexCell, sm.GetStyleV2(TestIndexStyle))

		testNameCell := fmt.Sprintf("B%d", row)
		f.SetCellStyle("Sheet1", testNameCell, testNameCell, sm.GetStyleV2(TestNameStyle))

		testQuantityCell := fmt.Sprintf("C%d", row)
		f.SetCellStyle("Sheet1", testQuantityCell, testQuantityCell, sm.GetStyleV2(TestQuantityStyle))

		testPriceCell := fmt.Sprintf("D%d", row)
		testPriceTotalCell := fmt.Sprintf("E%d", row)

		// Apply price style with number formatting for thousand separators
		f.SetCellStyle("Sheet1", testPriceCell, testPriceCell, sm.GetStyleV2(TestPriceStyle))
		f.SetCellStyle("Sheet1", testPriceTotalCell, testPriceTotalCell, sm.GetStyleV2(TestPriceStyle))

		err = f.SetRowHeight("Sheet1", row, 19.0)
		if err != nil {
			return nil, err
		}
	}

	// Total row - calculate sum of E column using Excel formula
	totalRow := startTestRow + len(record.TestResults)
	startTotalCell := fmt.Sprintf("A%d", totalRow)
	endTotalCell := fmt.Sprintf("D%d", totalRow)
	_ = r.MergeCellsWithStyle("Sheet1", startTotalCell, endTotalCell, "Tổng thành tiền", sm.GetStyleV2(TotalPriceLabelStyle))
	f.SetRowHeight("Sheet1", totalRow, 19.0)

	// Use SUM formula to calculate total price from all test rows
	totalPriceCell := fmt.Sprintf("E%d", totalRow)
	sumFormula := CreateSumFormula("E", startTestRow, totalRow-1)
	err = f.SetCellFormula("Sheet1", totalPriceCell, sumFormula)
	if err != nil {
		return nil, err
	}
	// Apply bold price style with number formatting for the total
	_ = f.SetCellStyle("Sheet1", totalPriceCell, totalPriceCell, sm.GetStyleV2(TestPriceTotalStyle))

	// Calculate print area based on content (A1 to E + last row with data)
	lastRow := startTestRow + len(record.TestResults) + 2 // Add buffer rows
	printArea := fmt.Sprintf("$A$1:$E$%d", lastRow)

	err = r.ApplyPrintArea(ctx, r.File, printArea)
	if err != nil {
		return nil, err
	}

	return r.GetIOReader(ctx)
}
