package sheets

import (
	"errors"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/xuri/excelize/v2"
)

func OpenTemplate(name models.ReportType) (*excelize.File, error) {
	supportedTemplates := map[models.ReportType]string{
		models.BillingReport:        "templates/PhieuThu.xlsx",
		models.ResultsReport:        "templates/PhieuKetQua.xlsx",
		models.ResultsWithSignature: "templates/PhieuKetQuaChuKy.xlsx",
		models.TrackingReport:       "templates/PhieuTheoDoi.xlsx",
	}

	templatePath, ok := supportedTemplates[name]
	if !ok {
		return nil, errors.New("template is not supported")
	}

	return excelize.OpenFile(templatePath)
}
