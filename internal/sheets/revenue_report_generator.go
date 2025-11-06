package sheets

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/datdev2409/lab-admin-go/internal/models"
	"go.uber.org/zap"
)

type RevenueExportReport struct {
	*BaseReportBuilder
}

func NewRevenueExportReport(ctx context.Context) (*RevenueExportReport, error) {
	pageSetup := &PageSetup{
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
	}

	builder, err := NewBaseReportBuilder(ctx, pageSetup)
	if err != nil {
		return nil, err
	}

	if err := builder.InitializeNewFile(ctx); err != nil {
		return nil, err
	}

	return &RevenueExportReport{BaseReportBuilder: builder}, nil
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
	_ = r.MergeCellsWithStyle("Sheet1", titleCell, fmt.Sprintf("G%d", currentRow), "Báo Cáo Doanh Thu", sm.GetStyleV2(ReportNameStyle))
	f.SetRowHeight("Sheet1", currentRow, 25.0)

	currentRow++

	// 2. Date Range Row
	dateRangeCell := fmt.Sprintf("A%d", currentRow)
	startDateStr := reportData.Summary.StartDate.In(vietnamLocation).Format("02/01/2006")
	endDateStr := reportData.Summary.EndDate.In(vietnamLocation).Format("02/01/2006")
	dateRangeText := fmt.Sprintf("Từ %s đến %s", startDateStr, endDateStr)
	_ = r.MergeCellsWithStyle("Sheet1", dateRangeCell, fmt.Sprintf("G%d", currentRow), dateRangeText, sm.GetStyleV2(ReportDateStyle))
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

	dataStartRow := currentRow
	for _, record := range reportData.Records {
		dataRow := currentRow

		// Format date with Vietnam timezone
		dateStr := ""
		if !record.CreatedAt.IsZero() {
			dateStr = record.CreatedAt.In(vietnamLocation).Format("02/01/2006")
		}

		// Handle missing doctor name - just show empty
		doctorName := record.DoctorName

		// Prepare row data with numeric price value for proper SUM formula calculation
		// Use nil for index column, will be set as formula
		rowData := []interface{}{
			nil,                    // STT (will be set as formula)
			dateStr,                // Ngày
			doctorName,             // Bác sĩ (empty if not set)
			record.Patient.Name,    // Họ tên
			record.Patient.Address, // Địa chỉ
			record.Patient.Phone,   // Số điện thoại
			record.TotalPrice,      // Thành tiền (numeric value, not formatted string)
		}

		rowStartCell := fmt.Sprintf("A%d", dataRow)
		err := f.SetSheetRow("Sheet1", rowStartCell, &rowData)
		if err != nil {
			log.Error("Failed to set data row", zap.Error(err), zap.Int("row", dataRow))
			return nil, err
		}

		// Set the auto-increment formula for the index cell (Column A)
		indexCell := fmt.Sprintf("A%d", dataRow)
		indexFormula := SetAutoIncrementIndexFormula(dataStartRow)
		err = f.SetCellFormula("Sheet1", indexCell, indexFormula)
		if err != nil {
			log.Error("Failed to set index formula", zap.Error(err), zap.Int("row", dataRow))
			return nil, err
		}

		// Apply cell styling with borders
		for col := 0; col < len(rowData); col++ {
			colLetter := string(rune('A' + col))
			cellRef := fmt.Sprintf("%s%d", colLetter, dataRow)

			// Determine style based on column - all use consistent font12 with borders
			var styleName StyleName
			if col == 0 { // STT column (center aligned)
				styleName = TestIndexStyle
			} else if col == 6 { // Thành tiền column (right aligned with number formatting)
				styleName = TestPriceStyle
			} else {
				styleName = RevenueTableDataStyle
			}

			f.SetCellStyle("Sheet1", cellRef, cellRef, sm.GetStyleV2(styleName))
		}

		f.SetRowHeight("Sheet1", dataRow, 18.0)
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

	// Use SUM formula to calculate total revenue from column G (Thành tiền)
	totalRevenueCell := fmt.Sprintf("G%d", totalRow)
	sumFormula := CreateSumFormula("G", dataStartRow, totalRow-1)
	err = f.SetCellFormula("Sheet1", totalRevenueCell, sumFormula)
	if err != nil {
		log.Error("Failed to set total revenue formula", zap.Error(err))
		return nil, err
	}

	// Apply total row styling with number formatting for the revenue cell
	for col := 0; col < 7; col++ {
		colLetter := string(rune('A' + col))
		cellRef := fmt.Sprintf("%s%d", colLetter, totalRow)

		// Use bold price style for the revenue total cell
		var styleName StyleName
		if col == 6 { // Column G - use bold price style
			styleName = TestPriceTotalStyle
		} else {
			styleName = TotalPriceStyle
		}
		f.SetCellStyle("Sheet1", cellRef, cellRef, sm.GetStyleV2(styleName))
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
