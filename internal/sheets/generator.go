package sheets

import (
	"context"
	"fmt"
	"io"

	"github.com/datdev2409/lab-admin-go/internal/models"
)

type ReportGenerator interface {
	Generate(ctx context.Context, data interface{}) (io.Reader, error)
}

func NewReportGenerator(ctx context.Context, reportType models.ReportType) (ReportGenerator, error) {
	switch reportType {
	case models.BillingReport:
		return NewBillingReport(ctx)
	default:
		return nil, fmt.Errorf("report type not support")
	}
}
