package sheets

import (
	"context"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type BaseReport struct {
	TemplateFilePath string
	OutputFilePath   string
	PageSetup
}

func (br *BaseReport) Open(ctx context.Context) (*excelize.File, error) {
	f, err := excelize.OpenFile(br.TemplateFilePath)
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to open template file", zap.String("templateFilePath", br.TemplateFilePath), zap.Error(err))
		return nil, err
	}

	return f, nil
}

func (br *BaseReport) Save(ctx context.Context, f *excelize.File) (string, error) {
	err := f.SaveAs(br.OutputFilePath)
	if err != nil {
		logger.FromCtx(ctx).Error("Failed to save output file", zap.String("outputFilePath", br.OutputFilePath), zap.Error(err))
		return "", err
	}

	return br.OutputFilePath, nil
}

type Cell struct {
	value     interface{}
	styleName *StyleName
}
