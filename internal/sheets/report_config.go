package sheets

import (
	"fmt"
	"strings"
	"time"

	"github.com/datdev2409/lab-admin-go/internal/models"
)

func NewReportWithSingleRecord(reportType models.ReportType, record models.Record) (*BaseReport, error) {
	switch reportType {
	case models.BillingReport:
		return &BaseReport{
			TemplateFilePath: "templates/PhieuThu.xlsx",
			OutputFilePath:   fmt.Sprintf("reports/%s-%s-hoa-don.xlsx", time.Now().Format("20060102"), strings.ReplaceAll(record.Patient.Name, " ", "_")),
			PageSetup: PageSetup{
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
			},
		}, nil

	case models.ResultsReport:
		return &BaseReport{
			TemplateFilePath: "templates/PhieuKetQua.xlsx",
			OutputFilePath:   fmt.Sprintf("reports/%s-%s-ket-qua.xlsx", time.Now().Format("20060102"), strings.ReplaceAll(record.Patient.Name, " ", "_")),
			PageSetup: PageSetup{
				PageSize:    9,
				Orientation: "portrait",
				Margins: MarginConfig{
					Top:    1.9,
					Bottom: float64(0),
					Left:   0.4,
					Right:  0.4,
					Header: 0.236220472440945,
					Footer: float64(0),
				},
			},
		}, nil

	case models.ResultsWithSignature:
		return &BaseReport{
			TemplateFilePath: "templates/PhieuKetQuaChuKy.xlsx",
			OutputFilePath:   fmt.Sprintf("reports/%s-%s-ket-qua-online.xlsx", time.Now().Format("20060102"), strings.ReplaceAll(record.Patient.Name, " ", "_")),
			PageSetup: PageSetup{
				PageSize:    9,
				Orientation: "portrait",
				Margins: MarginConfig{
					Top:    0.31496063,
					Bottom: float64(0),
					Left:   0.4,
					Right:  0.4,
					Header: 0.511811023622047,
					Footer: float64(0),
				},
			},
		}, nil

	case models.ResultsWithSignaturePDF:
		return &BaseReport{
			TemplateFilePath: "templates/PhieuKetQuaChuKy.xlsx", // Same template as ResultsWithSignature
			OutputFilePath:   fmt.Sprintf("reports/%s-%s-ket-qua-online.xlsx", time.Now().Format("20060102"), strings.ReplaceAll(record.Patient.Name, " ", "_")),
			PageSetup: PageSetup{
				PageSize:    9,
				Orientation: "portrait",
				Margins: MarginConfig{
					Top:    0.31496063,
					Bottom: float64(0),
					Left:   0.4,
					Right:  0.4,
					Header: 0.511811023622047,
					Footer: float64(0),
				},
			},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported report type: %s", reportType)
	}
}

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
