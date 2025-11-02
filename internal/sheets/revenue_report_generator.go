package sheets

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type RevenueExportReport struct {
	*PageSetup
	*ReportFile
}

func NewRevenueExportReport(ctx context.Context) (*RevenueExportReport, error) {
	report := &RevenueExportReport{
		ReportFile: &ReportFile{
			File: nil,
		},
		PageSetup: &PageSetup{
			SheetName:   "Sheet1",
			PageSize:    9,
			Orientation: "portrait",
			Margins: MarginConfig{
				Top:    0.75,
				Bottom: 0.75,
				Left:   0.5,
				Right:  0.5,
				Header: 0.236220472440945,
				Footer: 0.511811023622047,
			},
			ColumnWidth: map[string]float64{
				"A": 5.0,  // STT
				"B": 12.0, // Ngày
				"C": 20.0, // Bác sĩ
				"D": 25.0, // Họ tên
				"E": 30.0, // Địa chỉ
				"F": 15.0, // Số điện thoại
				"G": 15.0, // Thành tiền
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

func (r *RevenueExportReport) Generate(ctx context.Context, data interface{}) (io.Reader, error) {
	reportData, ok := data.(*models.ReportResponse)
	if !ok {
		return nil, fmt.Errorf("invalid data type for revenue export report generation")
	}

	f := r.File
	defer f.Close()

	sm := NewStyleManager(ctx, f)
	log := logger.FromCtx(ctx)

	// Load Vietnam timezone for date display
	vietnamLocation, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Error("Failed to load Vietnam timezone", zap.Error(err))
		vietnamLocation = time.UTC // Fallback to UTC
	}

	var currentRow = 1

	// 1. Title Row
	titleCell := fmt.Sprintf("A%d", currentRow)
	f.MergeCell("Sheet1", titleCell, fmt.Sprintf("G%d", currentRow))
	f.SetCellValue("Sheet1", titleCell, "Báo Cáo Doanh Thu")
	f.SetCellStyle("Sheet1", titleCell, fmt.Sprintf("G%d", currentRow), sm.GetStyleV2(ReportNameStyle))
	f.SetRowHeight("Sheet1", currentRow, 25.0)

	currentRow++

	// 2. Date Range Row
	dateRangeCell := fmt.Sprintf("A%d", currentRow)
	f.MergeCell("Sheet1", dateRangeCell, fmt.Sprintf("G%d", currentRow))

	startDateStr := reportData.Summary.StartDate.In(vietnamLocation).Format("02/01/2006")
	endDateStr := reportData.Summary.EndDate.In(vietnamLocation).Format("02/01/2006")
	dateRangeText := fmt.Sprintf("Từ %s đến %s", startDateStr, endDateStr)

	f.SetCellValue("Sheet1", dateRangeCell, dateRangeText)
	f.SetCellStyle("Sheet1", dateRangeCell, fmt.Sprintf("G%d", currentRow), sm.GetStyleV2(ReportDateStyle))
	f.SetRowHeight("Sheet1", currentRow, 18.0)

	currentRow += 2 // Add spacing

	// 3. Table Headers
	tableHeaderRow := currentRow
	tableHeaderStartCell := fmt.Sprintf("A%d", tableHeaderRow)

	tableHeaders := []string{"STT", "Ngày", "Bác sĩ", "Họ tên", "Địa chỉ", "Số điện thoại", "Thành tiền"}
	headerErr := f.SetSheetRow("Sheet1", tableHeaderStartCell, &tableHeaders)
	if headerErr != nil {
		log.Error("Failed to set table headers", zap.Error(headerErr))
		return nil, headerErr
	}

	// Apply header styling to all header cells
	for col := 0; col < len(tableHeaders); col++ {
		colLetter := string(rune('A' + col))
		headerCell := fmt.Sprintf("%s%d", colLetter, tableHeaderRow)
		f.SetCellStyle("Sheet1", headerCell, headerCell, sm.GetStyleV2(TestTableHeaderStyle))
	}

	f.SetRowHeight("Sheet1", tableHeaderRow, 20.0)
	currentRow++

	// 4. Data Rows
	if reportData.Records == nil {
		log.Warn("No records in report data")
		reportData.Records = []*models.RecordWithTotal{}
	}

	totalRevenue := 0
	for idx, record := range reportData.Records {
		dataRow := currentRow
		rowNum := idx + 1

		// Format date with Vietnam timezone
		dateStr := ""
		if !record.CreatedAt.IsZero() {
			dateStr = record.CreatedAt.In(vietnamLocation).Format("02/01/2006")
		}

		// Handle missing doctor name - just show empty
		doctorName := record.DoctorName

		// Prepare row data
		rowData := []interface{}{
			rowNum,                         // STT
			dateStr,                        // Ngày
			doctorName,                     // Bác sĩ (empty if not set)
			record.Patient.Name,            // Họ tên
			record.Patient.Address,         // Địa chỉ
			record.Patient.Phone,           // Số điện thoại
			FormatPrice(record.TotalPrice), // Thành tiền
		}

		rowStartCell := fmt.Sprintf("A%d", dataRow)
		err := f.SetSheetRow("Sheet1", rowStartCell, &rowData)
		if err != nil {
			log.Error("Failed to set data row", zap.Error(err), zap.Int("row", dataRow))
			return nil, err
		}

		// Apply cell styling with borders
		for col := 0; col < len(rowData); col++ {
			colLetter := string(rune('A' + col))
			cellRef := fmt.Sprintf("%s%d", colLetter, dataRow)

			// Determine style based on column - all use consistent font12 with borders
			var styleName StyleName
			if col == 6 { // Thành tiền column (right aligned)
				styleName = RevenueTableDataRightStyle
			} else {
				styleName = RevenueTableDataStyle
			}

			f.SetCellStyle("Sheet1", cellRef, cellRef, sm.GetStyleV2(styleName))
		}

		f.SetRowHeight("Sheet1", dataRow, 18.0)
		totalRevenue += record.TotalPrice
		currentRow++
	}

	// 5. Total Row
	totalRow := currentRow
	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", totalRow), "")
	f.SetCellValue("Sheet1", fmt.Sprintf("B%d", totalRow), "")
	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", totalRow), "")
	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", totalRow), "")
	f.SetCellValue("Sheet1", fmt.Sprintf("E%d", totalRow), "")
	f.SetCellValue("Sheet1", fmt.Sprintf("F%d", totalRow), "Tổng cộng")
	f.SetCellValue("Sheet1", fmt.Sprintf("G%d", totalRow), FormatPrice(totalRevenue))

	// Apply total row styling
	for col := 0; col < 7; col++ {
		colLetter := string(rune('A' + col))
		cellRef := fmt.Sprintf("%s%d", colLetter, totalRow)
		f.SetCellStyle("Sheet1", cellRef, cellRef, sm.GetStyleV2(TotalPriceStyle))
	}

	f.SetRowHeight("Sheet1", totalRow, 20.0)

	// 6. Set print area
	lastRow := totalRow + 1
	printArea := fmt.Sprintf("$A$1:$G$%d", lastRow)
	err = r.ApplyPrintArea(ctx, f, printArea)
	if err != nil {
		log.Error("Failed to apply print area", zap.Error(err))
		return nil, err
	}

	return r.GetIOReader(ctx)
}
