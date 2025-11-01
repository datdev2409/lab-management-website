package sheets

import (
	"fmt"
	"strings"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
)

// NewReportWithMultipleRecords creates a report configuration for reports that handle multiple records
func NewReportWithMultipleRecords(reportType models.ReportType, records []*models.Record) (*BaseReport, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("cannot create report with empty records")
	}

	switch reportType {
	case models.TrackingReport:
		return &BaseReport{
			TemplateFilePath: "templates/PhieuTheoDoi.xlsx",
			OutputFilePath:   fmt.Sprintf("reports/%s-%s-theo-doi.xlsx", time.Now().Format("20060102"), strings.ReplaceAll(records[0].Patient.Name, " ", "_")),
			PageSetup: PageSetup{
				PageSize:    9,
				Orientation: "landscape",
				SheetName:   "Sheet1",
				Margins: MarginConfig{
					Top:    0.25,
					Bottom: 0.5,
					Left:   0.5,
					Right:  0.5,
					Header: 0.5,
					Footer: 0.5,
				},
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported multiple records report type: %s", reportType)
	}
}
