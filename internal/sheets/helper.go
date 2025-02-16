package sheets

import (
	"errors"

	"github.com/xuri/excelize/v2"
)

func OpenTemplate(name string) (*excelize.File, error) {
	supportedTemplates := map[string]string{
		"record_billing": "templates/record_billing.xlsx",
	}

	templatePath, ok := supportedTemplates[name]
	if !ok {
		return nil, errors.New("Template is not supported")
	}

	return excelize.OpenFile(templatePath)
}

// func WriteTable()
