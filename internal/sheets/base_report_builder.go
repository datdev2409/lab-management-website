package sheets

import (
	"context"
	"fmt"

	"github.com/xuri/excelize/v2"
)

// BaseReportBuilder provides common initialization logic for all report generators
type BaseReportBuilder struct {
	*PageSetup
	*ReportFile
}

// NewBaseReportBuilder creates a new base report builder with the given configuration
func NewBaseReportBuilder(ctx context.Context, pageSetup *PageSetup) (*BaseReportBuilder, error) {
	builder := &BaseReportBuilder{
		PageSetup: pageSetup,
		ReportFile: &ReportFile{
			File: nil,
		},
	}
	return builder, nil
}

// InitializeNewFile creates a new Excel file and applies configuration
func (b *BaseReportBuilder) InitializeNewFile(ctx context.Context) error {
	b.File = excelize.NewFile()

	if err := b.ApplyColumnWidths(ctx, b.File); err != nil {
		return fmt.Errorf("failed to apply column widths: %w", err)
	}

	if err := b.ApplyPageSetupV2(ctx, b.File); err != nil {
		return fmt.Errorf("failed to apply page setup: %w", err)
	}

	return nil
}

// InitializeFromTemplate opens an Excel template file and applies configuration
func (b *BaseReportBuilder) InitializeFromTemplate(ctx context.Context, templatePath string) error {
	if err := b.OpenTemplate(ctx, templatePath); err != nil {
		return fmt.Errorf("failed to open template: %w", err)
	}

	if err := b.ApplyColumnWidths(ctx, b.File); err != nil {
		return fmt.Errorf("failed to apply column widths: %w", err)
	}

	if err := b.ApplyPageSetupV2(ctx, b.File); err != nil {
		return fmt.Errorf("failed to apply page setup: %w", err)
	}

	return nil
}

// SetCellWithStyle is a helper to set cell value and apply style in one call
func (b *BaseReportBuilder) SetCellWithStyle(sheetName, cell string, value interface{}, styleID int) error {
	if err := b.File.SetCellValue(sheetName, cell, value); err != nil {
		return err
	}
	if styleID >= 0 {
		if err := b.File.SetCellStyle(sheetName, cell, cell, styleID); err != nil {
			return err
		}
	}
	return nil
}

// MergeCellsWithStyle merges cells and applies style to the merged area
func (b *BaseReportBuilder) MergeCellsWithStyle(sheetName, startCell, endCell string, value interface{}, styleID int) error {
	if err := b.File.MergeCell(sheetName, startCell, endCell); err != nil {
		return err
	}
	if err := b.File.SetCellValue(sheetName, startCell, value); err != nil {
		return err
	}
	if styleID >= 0 {
		if err := b.File.SetCellStyle(sheetName, startCell, endCell, styleID); err != nil {
			return err
		}
	}
	return nil
}
