package sheets

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type MarginConfig struct {
	Top    float64 // in inches
	Bottom float64 // in inches
	Left   float64 // in inches
	Right  float64 // in inches
	Header float64 // in inches
	Footer float64 // in inches
}

type PageSetup struct {
	SheetName   string
	PageSize    int    // e.g., excelize.PaperA4
	Orientation string // "portrait" or "landscape"
	Margins     MarginConfig
	ColumnWidth map[string]float64 // column letter to width in characters
}

func (ps *PageSetup) ApplyColumnWidths(ctx context.Context, f *excelize.File) error {
	for col, width := range ps.ColumnWidth {
		err := f.SetColWidth(ps.SheetName, col, col, width+0.89)
		if err != nil {
			logger.FromCtx(ctx).Error("Failed to set column width", zap.String("sheetName", ps.SheetName), zap.String("column", col), zap.Float64("width", width), zap.Error(err))
			return err
		}
	}

	return nil
}

func (ps *PageSetup) ApplyPrintArea(ctx context.Context, f *excelize.File, printArea string) error {
	printAreaZone := &excelize.DefinedName{
		Name:     "_xlnm.Print_Area",
		RefersTo: "'" + ps.SheetName + "'!" + printArea,
		Scope:    ps.SheetName,
	}
	_ = f.DeleteDefinedName(printAreaZone)
	err := f.SetDefinedName(printAreaZone)
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to set print area", zap.String("sheetName", ps.SheetName), zap.Error(err))
		return err
	}
	return nil
}

func (ps *PageSetup) ApplyPageSetupV2(ctx context.Context, f *excelize.File) error {
	err := f.SetPageLayout(ps.SheetName, &excelize.PageLayoutOptions{
		Size:        &ps.PageSize,
		Orientation: &ps.Orientation,
	})
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to set page layout", zap.String("sheetName", ps.SheetName), zap.Error(err))
		return err
	}

	err = f.SetPageMargins(ps.SheetName, &excelize.PageLayoutMarginsOptions{
		Top:    &ps.Margins.Top,
		Bottom: &ps.Margins.Bottom,
		Left:   &ps.Margins.Left,
		Right:  &ps.Margins.Right,
		Header: &ps.Margins.Header,
		Footer: &ps.Margins.Footer,
	})
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to set page margins", zap.String("sheetName", ps.SheetName), zap.Error(err))
		return err
	}

	return nil
}
