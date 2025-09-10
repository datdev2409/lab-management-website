package sheets

import (
	"context"
	"fmt"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// StyleManager handles creation and caching of Excel styles
type StyleManager struct {
	file   *excelize.File
	ctx    context.Context
	styles map[StyleName]int // Cache created styles by name
}

// StyleName represents available style names
type StyleName string

const (
	StylePatientName StyleName = "patientName"
	StylePatientInfo StyleName = "patientInfo"
	StyleDateCenter  StyleName = "dateCenter"
	StyleTestResult  StyleName = "testResult"
	StyleTestName    StyleName = "testName"
	StyleAbnormal    StyleName = "abnormal"
)

// NewStyleManager creates a new StyleManager instance
func NewStyleManager(ctx context.Context, file *excelize.File) *StyleManager {
	return &StyleManager{
		file:   file,
		ctx:    ctx,
		styles: make(map[StyleName]int),
	}
}

// getStandardBorder returns the standard border configuration used across multiple styles
func (sm *StyleManager) getStandardBorder() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
	}
}

// GetStyle returns the style ID for the given style name
func (sm *StyleManager) GetStyle(styleName StyleName) (int, error) {
	switch styleName {
	case StylePatientName:
		return sm.GetPatientNameStyle()
	case StylePatientInfo:
		return sm.GetPatientInfoStyle()
	case StyleDateCenter:
		return sm.GetDateCenterStyle()
	case StyleTestResult:
		return sm.GetTestResultStyle()
	case StyleTestName:
		return sm.GetTestNameStyle()
	case StyleAbnormal:
		return sm.GetAbnormalStyle()
	default:
		return 0, fmt.Errorf("unknown style name: %s", styleName)
	}
}

// GetPatientNameStyle returns style for patient names (14pt, bold)
func (sm *StyleManager) GetPatientNameStyle() (int, error) {
	if styleID, exists := sm.styles[StylePatientName]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 14, Bold: true},
	})
	if err != nil {
		logger.FromCtx(sm.ctx).Debug("Failed to create patient name style", zap.Error(err))
		return 0, err
	}

	sm.styles[StylePatientName] = styleID
	return styleID, nil
}

// GetPatientInfoStyle returns style for patient information (12pt)
func (sm *StyleManager) GetPatientInfoStyle() (int, error) {
	if styleID, exists := sm.styles[StylePatientInfo]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 12},
	})
	if err != nil {
		logger.FromCtx(sm.ctx).Debug("Failed to create patient info style", zap.Error(err))
		return 0, err
	}

	sm.styles[StylePatientInfo] = styleID
	return styleID, nil
}

// GetDateCenterStyle returns style for centered date fields (12pt, center aligned)
func (sm *StyleManager) GetDateCenterStyle() (int, error) {
	if styleID, exists := sm.styles[StyleDateCenter]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 12},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		logger.FromCtx(sm.ctx).Debug("Failed to create date center style", zap.Error(err))
		return 0, err
	}

	sm.styles[StyleDateCenter] = styleID
	return styleID, nil
}

// GetTestResultStyle returns style for test results (13pt, center aligned with borders)
func (sm *StyleManager) GetTestResultStyle() (int, error) {
	if styleID, exists := sm.styles[StyleTestResult]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 13},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: sm.getStandardBorder(),
	})
	if err != nil {
		logger.FromCtx(sm.ctx).Debug("Failed to create test result style", zap.Error(err))
		return 0, err
	}

	sm.styles[StyleTestResult] = styleID
	return styleID, nil
}

// GetTestNameStyle returns style for test names (13pt, left aligned with borders)
func (sm *StyleManager) GetTestNameStyle() (int, error) {
	if styleID, exists := sm.styles[StyleTestName]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 13},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: sm.getStandardBorder(),
	})
	if err != nil {
		logger.FromCtx(sm.ctx).Debug("Failed to create test name style", zap.Error(err))
		return 0, err
	}

	sm.styles[StyleTestName] = styleID
	return styleID, nil
}

// GetAbnormalStyle returns style for abnormal test results (13pt, bold, center aligned with borders)
func (sm *StyleManager) GetAbnormalStyle() (int, error) {
	if styleID, exists := sm.styles[StyleAbnormal]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 13,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: sm.getStandardBorder(),
	})
	if err != nil {
		logger.FromCtx(sm.ctx).Debug("Failed to create abnormal style", zap.Error(err))
		return 0, err
	}

	sm.styles[StyleAbnormal] = styleID
	return styleID, nil
}
