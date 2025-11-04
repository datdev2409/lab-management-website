package sheets

import (
	"context"
	"fmt"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// HeaderComponent represents a reusable header section for reports
type HeaderComponent struct {
	file         *excelize.File
	styleManager *StyleManager
	sheetName    string
	startRow     int
	startCol     string
	includeDate  bool
	reportTitle  string
	dateValue    string
}

// NewHeaderComponent creates a new header component
func NewHeaderComponent(
	file *excelize.File,
	styleManager *StyleManager,
	sheetName string,
	startRow int,
	startCol string,
	reportTitle string,
	includeDate bool,
) *HeaderComponent {
	return &HeaderComponent{
		file:         file,
		styleManager: styleManager,
		sheetName:    sheetName,
		startRow:     startRow,
		startCol:     startCol,
		includeDate:  includeDate,
		reportTitle:  reportTitle,
	}
}

// Apply renders the header component
func (h *HeaderComponent) Apply(ctx context.Context) error {
	f := h.file
	sm := h.styleManager

	currentRow := h.startRow
	endCol := "E" // Default end column for merging

	// Lab name
	labNameCell := fmt.Sprintf("%s%d", h.startCol, currentRow)
	if err := f.MergeCell(h.sheetName, labNameCell, fmt.Sprintf("%s%d", endCol, currentRow)); err != nil {
		return err
	}
	if err := f.SetCellValue(h.sheetName, labNameCell, "PHÒNG XÉT NGHIỆM Y KHOA ANH QUÂN"); err != nil {
		return err
	}
	if err := f.SetCellStyle(h.sheetName, labNameCell, labNameCell, sm.GetStyleV2(LabNameStyle)); err != nil {
		return err
	}
	currentRow++

	// Lab address
	labAddressCell := fmt.Sprintf("%s%d", h.startCol, currentRow)
	if err := f.MergeCell(h.sheetName, labAddressCell, fmt.Sprintf("%s%d", endCol, currentRow)); err != nil {
		return err
	}
	if err := f.SetCellValue(h.sheetName, labAddressCell, "60 Đống Đa, Phường Cao Lãnh, Đồng Tháp"); err != nil {
		return err
	}
	if err := f.SetCellStyle(h.sheetName, labAddressCell, labAddressCell, sm.GetStyleV2(LabAddressStyle)); err != nil {
		return err
	}
	currentRow++

	// Report title
	reportTitleCell := fmt.Sprintf("%s%d", h.startCol, currentRow)
	if err := f.MergeCell(h.sheetName, reportTitleCell, fmt.Sprintf("%s%d", endCol, currentRow)); err != nil {
		return err
	}
	if err := f.SetCellValue(h.sheetName, reportTitleCell, h.reportTitle); err != nil {
		return err
	}
	if err := f.SetCellStyle(h.sheetName, reportTitleCell, reportTitleCell, sm.GetStyleV2(ReportNameStyle)); err != nil {
		return err
	}
	currentRow++

	// Date (if included)
	if h.includeDate {
		dateCell := fmt.Sprintf("%s%d", h.startCol, currentRow)
		if err := f.MergeCell(h.sheetName, dateCell, fmt.Sprintf("%s%d", endCol, currentRow)); err != nil {
			return err
		}
		dateText := h.dateValue
		if dateText == "" {
			now := time.Now()
			dateText = fmt.Sprintf("Ngày: %s", now.Format("02/01/2006"))
		}
		if err := f.SetCellValue(h.sheetName, dateCell, dateText); err != nil {
			return err
		}
		if err := f.SetCellStyle(h.sheetName, dateCell, dateCell, sm.GetStyleV2(ReportDateStyle)); err != nil {
			return err
		}
	}

	return nil
}

// SetDateValue sets a custom date value
func (h *HeaderComponent) SetDateValue(dateValue string) {
	h.dateValue = dateValue
}

// SignatureComponent represents a reusable signature section
type SignatureComponent struct {
	file         *excelize.File
	styleManager *StyleManager
	sheetName    string
	startRow     int
	startCol     rune
	endCol       rune
	endRow       int
}

// NewSignatureComponent creates a new signature component
func NewSignatureComponent(
	file *excelize.File,
	styleManager *StyleManager,
	sheetName string,
	startRow int,
	startCol, endCol rune,
) *SignatureComponent {
	return &SignatureComponent{
		file:         file,
		styleManager: styleManager,
		sheetName:    sheetName,
		startRow:     startRow,
		startCol:     startCol,
		endCol:       endCol,
	}
}

// Apply renders the signature section
func (s *SignatureComponent) Apply(ctx context.Context) error {
	f := s.file
	sm := s.styleManager
	signatureCol := string(s.startCol)

	// Location and date
	locationDateRow := s.startRow
	locationDateCell := fmt.Sprintf("%s%d", signatureCol, locationDateRow)
	now := time.Now()
	dateText := fmt.Sprintf("Cao Lãnh. Ngày %d tháng %d năm %d", now.Day(), int(now.Month()), now.Year())
	if err := f.SetCellValue(s.sheetName, locationDateCell, dateText); err != nil {
		return err
	}

	// Merge cells if needed
	if s.startCol != s.endCol {
		endLocationDateCell := fmt.Sprintf("%s%d", string(s.endCol), locationDateRow)
		if err := f.MergeCell(s.sheetName, locationDateCell, endLocationDateCell); err != nil {
			return err
		}
	}
	if err := f.SetCellStyle(s.sheetName, locationDateCell, locationDateCell, sm.GetStyleV2(LocationDateStyle)); err != nil {
		return err
	}

	// Lab department
	labDeptRow := locationDateRow + 1
	labDeptCell := fmt.Sprintf("%s%d", signatureCol, labDeptRow)
	if err := f.SetCellValue(s.sheetName, labDeptCell, "PHÒNG XÉT NGHIỆM"); err != nil {
		return err
	}

	// Merge cells if needed
	if s.startCol != s.endCol {
		endLabDeptCell := fmt.Sprintf("%s%d", string(s.endCol), labDeptRow)
		if err := f.MergeCell(s.sheetName, labDeptCell, endLabDeptCell); err != nil {
			return err
		}
	}
	if err := f.SetCellStyle(s.sheetName, labDeptCell, labDeptCell, sm.GetStyleV2(LabDepartmentStyle)); err != nil {
		return err
	}

	// Signature name (5 rows down for signature space)
	signatureNameRow := labDeptRow + 6
	signatureNameCell := fmt.Sprintf("%s%d", signatureCol, signatureNameRow)
	if err := f.SetCellValue(s.sheetName, signatureNameCell, "CKI.XN NGUYỄN CÔNG MẪN"); err != nil {
		return err
	}

	// Merge cells if needed
	if s.startCol != s.endCol {
		endSignatureNameCell := fmt.Sprintf("%s%d", string(s.endCol), signatureNameRow)
		if err := f.MergeCell(s.sheetName, signatureNameCell, endSignatureNameCell); err != nil {
			return err
		}
	}
	if err := f.SetCellStyle(s.sheetName, signatureNameCell, signatureNameCell, sm.GetStyleV2(SignatureNameStyle)); err != nil {
		return err
	}

	s.endRow = signatureNameRow
	return nil
}

// GetEndRow returns the last row used by the signature component
func (s *SignatureComponent) GetEndRow() int {
	return s.endRow
}

// AddLogoComponent adds the lab logo to the report
func AddLogoComponent(ctx context.Context, file *excelize.File, sheetName, cell, logoPath string, scaleX, scaleY float64) error {
	err := file.AddPicture(sheetName, cell, logoPath, &excelize.GraphicOptions{
		ScaleX:          scaleX,
		ScaleY:          scaleY,
		LockAspectRatio: true,
	})
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to add logo picture", zap.Error(err))
		return err
	}
	return nil
}
