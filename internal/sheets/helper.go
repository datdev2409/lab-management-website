package sheets

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/datdev2409/lab-admin-go/internal/models"
	"github.com/xuri/excelize/v2"
)

func OpenTemplate(name models.ReportType) (*excelize.File, error) {
	supportedTemplates := map[models.ReportType]string{
		models.BillingReport:           "templates/PhieuThu.xlsx",
		models.ResultsReport:           "templates/PhieuKetQua.xlsx",
		models.ResultsWithSignature:    "templates/PhieuKetQuaChuKy.xlsx",
		models.ResultsWithSignaturePDF: "templates/PhieuKetQuaOnlinePDF.xlsx",
		models.TrackingReport:          "templates/PhieuTheoDoi.xlsx",
	}

	templatePath, ok := supportedTemplates[name]
	if !ok {
		return nil, errors.New("template is not supported")
	}

	return excelize.OpenFile(templatePath)
}

func ToLowerCaseNonAccentVietnamese(str string) string {
	str = strings.ToLower(str)
	str = regexp.MustCompile(`[àáạảãâầấậẩẫăằắặẳẵ]`).ReplaceAllString(str, "a")
	str = regexp.MustCompile(`[èéẹẻẽêềếệểễ]`).ReplaceAllString(str, "e")
	str = regexp.MustCompile(`[ìíịỉĩ]`).ReplaceAllString(str, "i")
	str = regexp.MustCompile(`[òóọỏõôồốộổỗơờớợởỡ]`).ReplaceAllString(str, "o")
	str = regexp.MustCompile(`[ùúụủũưừứựửữ]`).ReplaceAllString(str, "u")
	str = regexp.MustCompile(`[ỳýỵỷỹ]`).ReplaceAllString(str, "y")
	str = regexp.MustCompile(`đ`).ReplaceAllString(str, "d")
	// Remove combining accent marks
	str = regexp.MustCompile(`[\\u0300\\u0301\\u0303\\u0309\\u0323]`).ReplaceAllString(str, "")
	// Remove Â, Ê, Ă, Ơ, Ư marks
	str = regexp.MustCompile(`[\\u02C6\\u0306\\u031B]`).ReplaceAllString(str, "")
	return str
}

// FormatPrice formats an integer price with comma separators
func FormatPrice(price int) string {
	if price == 0 {
		return "0"
	}
	
	// Convert to string
	priceStr := strconv.Itoa(price)
	
	// Add commas for thousands separator
	result := ""
	for i, char := range priceStr {
		if i > 0 && (len(priceStr)-i)%3 == 0 {
			result += ","
		}
		result += string(char)
	}
	
	return result
}

// FormatResult formats a result string to ensure at least 1 decimal point
// If the result is a whole number, adds .0
// If the result already has decimals, keeps the full value
func FormatResult(result string) string {
	if result == "" {
		return result
	}
	
	// Try to parse as float to check if it's a numeric value
	if val, err := strconv.ParseFloat(result, 64); err == nil {
		// If it's a whole number, format with 1 decimal place
		if val == float64(int64(val)) {
			return strconv.FormatFloat(val, 'f', 1, 64)
		}
		// If it already has decimals, keep the original precision
		return result
	}
	
	// If it's not a number, return as is
	return result
}
