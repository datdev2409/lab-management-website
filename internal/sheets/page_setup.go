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
	PageSize    int    // e.g., excelize.PaperA4
	Orientation string // "portrait" or "landscape"
	Margins     MarginConfig
}

func (ps *PageSetup) ApplyPageSetup(ctx context.Context, f *excelize.File, sheetName string, printArea string) error {
	err := f.SetPageLayout(sheetName, &excelize.PageLayoutOptions{
		Size:        &ps.PageSize,
		Orientation: &ps.Orientation,
	})
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to set page layout", zap.String("sheetName", sheetName), zap.Error(err))
		return err
	}

	err = f.SetPageMargins(sheetName, &excelize.PageLayoutMarginsOptions{
		Top:    &ps.Margins.Top,
		Bottom: &ps.Margins.Bottom,
		Left:   &ps.Margins.Left,
		Right:  &ps.Margins.Right,
		Header: &ps.Margins.Header,
		Footer: &ps.Margins.Footer,
	})
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to set page margins", zap.String("sheetName", sheetName), zap.Error(err))
		return err
	}

	printAreaZone := &excelize.DefinedName{
		Name:     "_xlnm.Print_Area",
		RefersTo: "'" + sheetName + "'!" + printArea,
		Scope:    sheetName,
	}
	f.DeleteDefinedName(printAreaZone) // Remove existing print area if any
	err = f.SetDefinedName(printAreaZone)
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to set print area", zap.String("sheetName", sheetName), zap.Error(err))
		return err
	}

	return nil
}
