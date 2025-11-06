package sheets

import (
	"fmt"
	"strconv"
)

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

func GetStyleNamePtr(styleName StyleName) *StyleName {
	return &styleName
}

// SetAutoIncrementIndexFormula sets an auto-increment formula for the index cell
// The formula uses ROW() to automatically calculate the index based on row position
// Parameters:
//   - startRow: The starting row number of the table data (first data row, not header)
//
// Returns: The formula string that calculates the relative index
// Example: For startRow=10, first data row will show 1, second row 2, etc.
func SetAutoIncrementIndexFormula(startRow int) string {
	return "=ROW()-" + strconv.Itoa(startRow-1)
}

// CreateSumFormula creates a SUM formula for a range of cells
// Parameters:
//   - column: The column letter (e.g., "E" or "G")
//   - startRow: The starting row number of the data (inclusive)
//   - endRow: The ending row number of the data (inclusive)
//
// Returns: The formula string (e.g., "=SUM(E10:E15)")
// Example: CreateSumFormula("E", 10, 15) returns "=SUM(E10:E15)"
func CreateSumFormula(column string, startRow, endRow int) string {
	return fmt.Sprintf("=SUM(%s%d:%s%d)", column, startRow, column, endRow)
}
