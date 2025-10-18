package sheets

import (
	"context"
	"fmt"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

const FontMyriadPro = "MyRIAD PRO"

// StyleName represents available style names
type StyleName int

const (
	StylePatientName StyleName = iota
	StylePatientInfo
	StyleDateCenter
	StyleTestResult
	StyleTestName
	StyleAbnormal
	StylePriceRight
	LabNameStyle
	LabAddressStyle
	ReportNameStyle
	ReportDateStyle
	PatientInfoStyle
	PatientNameStyle
	TestTableHeaderStyle
	TotalPriceLabelStyle
	TotalPriceStyle
	TestIndexStyle
	TestNameStyle
	TestQuantityStyle
	TestPriceStyle
	StylePatientNameLargeCenter
)

// CommonStyles holds commonly used style IDs to reduce repetitive style retrieval
type CommonStyles struct {
	PatientName            int
	PatientInfo            int
	DateCenter             int
	TestResult             int
	TestName               int
	Abnormal               int
	PriceRight             int
	PatientNameLargeCenter int
}

// GetPriceRightStyle returns style for price cells (same as testResultStyle but right aligned)
func (sm *StyleManager) GetPriceRightStyle() (int, error) {
	if styleID, exists := sm.cache[StylePriceRight]; exists {
		return styleID, nil
	}
	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 13, Family: FontMyriadPro},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
		Border: sm.getStandardBorder(),
	})
	if err != nil {
		logger.FromCtx(sm.ctx).Debug("Failed to create price right style", zap.Error(err))
		return 0, err
	}
	sm.cache[StylePriceRight] = styleID
	return styleID, nil
}

// StyleManager handles creation and caching of Excel styles
type StyleManager struct {
	file   *excelize.File
	ctx    context.Context
	styles map[StyleName]*excelize.Style // Cache created styles by name
	cache  map[StyleName]int
}

// NewStyleManager creates a new StyleManager instance
func NewStyleManager(ctx context.Context, file *excelize.File) *StyleManager {
	sm := &StyleManager{
		file:  file,
		ctx:   ctx,
		cache: make(map[StyleName]int),
	}

	alignCenter := &excelize.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	}

	alignLeft := &excelize.Alignment{
		Horizontal: "left",
		Vertical:   "center",
	}
	alignRight := &excelize.Alignment{
		Horizontal: "right",
		Vertical:   "center",
	}

	font12 := &excelize.Font{Size: 12, Family: FontMyriadPro}
	font11 := &excelize.Font{Size: 11, Family: FontMyriadPro}
	font11Bold := &excelize.Font{Size: 11, Family: FontMyriadPro, Bold: true}

	border := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
	}

	err := file.SetDefaultFont(FontMyriadPro)
	if err != nil {
		return nil
	}

	sm.styles = map[StyleName]*excelize.Style{
		LabNameStyle: {
			Font:      &excelize.Font{Size: 18, Family: FontMyriadPro},
			Alignment: alignCenter,
		},

		LabAddressStyle: {
			Font:      font12,
			Alignment: alignCenter,
		},
		ReportNameStyle: {
			Font:      &excelize.Font{Size: 18, Family: FontMyriadPro},
			Alignment: alignCenter,
		},
		ReportDateStyle: {
			Font:      font12,
			Alignment: alignCenter,
		},
		PatientInfoStyle: {
			Font:      font12,
			Alignment: alignLeft,
		},
		PatientNameStyle: {
			Font:      &excelize.Font{Size: 14, Family: FontMyriadPro, Bold: true},
			Alignment: alignLeft,
		},
		TestTableHeaderStyle: {
			Font:      font12,
			Alignment: alignCenter,
			Border:    border,
		},
		TotalPriceLabelStyle: {
			Font:      font11Bold,
			Alignment: alignCenter,
			Border:    border,
		},
		TotalPriceStyle: {
			Font:      font11Bold,
			Alignment: alignRight,
			Border:    border,
		},
		TestIndexStyle: {
			Font:      font11,
			Alignment: alignCenter,
			Border:    border,
		},
		TestNameStyle: {
			Font:      font11,
			Alignment: alignLeft,
			Border:    border,
		},
		TestQuantityStyle: {
			Font:      font11,
			Alignment: alignCenter,
			Border:    border,
		},
		TestPriceStyle: {
			Font:      font11,
			Alignment: alignRight,
			Border:    border,
		},
	}

	return sm
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

func (sm *StyleManager) GetStyleV2(styleName StyleName) int {
	if styleID, exists := sm.cache[styleName]; exists {
		return styleID
	}

	style, exists := sm.styles[styleName]
	if !exists {
		logger.FromCtx(sm.ctx).Error(fmt.Sprintf("style %s not found", styleName))
		return -1
	}

	styleId, err := sm.file.NewStyle(style)
	if err != nil {
		logger.FromCtx(sm.ctx).Error(fmt.Sprintf("can not create style %s", styleName))
		return -1
	}

	sm.cache[styleName] = styleId
	return styleId
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

// GetCommonStyles returns all commonly used styles in a single call to reduce code duplication
func (sm *StyleManager) GetCommonStyles() (*CommonStyles, error) {
	patientNameStyle, err := sm.GetPatientNameStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to get patient name style: %w", err)
	}

	patientInfoStyle, err := sm.GetPatientInfoStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to get patient info style: %w", err)
	}

	dateCenterStyle, err := sm.GetDateCenterStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to get date center style: %w", err)
	}

	testResultStyle, err := sm.GetTestResultStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to get test result style: %w", err)
	}

	testNameStyle, err := sm.GetTestNameStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to get test name style: %w", err)
	}

	abnormalStyle, err := sm.GetAbnormalStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to get abnormal style: %w", err)
	}

	priceRightStyle, err := sm.GetPriceRightStyle()
	if err != nil {
		return nil, fmt.Errorf("failed to get price right style: %w", err)
	}

	patientNameLargeCenterStyle, err := sm.GetPatientNameLargeCenter()
	if err != nil {
		return nil, fmt.Errorf("failed to get patient name large center style: %w", err)
	}

	return &CommonStyles{
		PatientName:            patientNameStyle,
		PatientInfo:            patientInfoStyle,
		DateCenter:             dateCenterStyle,
		TestResult:             testResultStyle,
		TestName:               testNameStyle,
		Abnormal:               abnormalStyle,
		PriceRight:             priceRightStyle,
		PatientNameLargeCenter: patientNameLargeCenterStyle,
	}, nil
}

// GetPatientNameStyle returns style for patient names (14pt, bold)
func (sm *StyleManager) GetPatientNameStyle() (int, error) {
	if styleID, exists := sm.cache[StylePatientName]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 14, Bold: true, Family: FontMyriadPro},
	})
	if err != nil {
		logger.FromCtx(sm.ctx).Debug("Failed to create patient name style", zap.Error(err))
		return 0, err
	}

	sm.cache[StylePatientName] = styleID
	return styleID, nil
}

// GetPatientInfoStyle returns style for patient information (12pt)
func (sm *StyleManager) GetPatientInfoStyle() (int, error) {
	if styleID, exists := sm.cache[StylePatientInfo]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 12, Family: FontMyriadPro},
	})
	if err != nil {
		logger.FromCtx(sm.ctx).Debug("Failed to create patient info style", zap.Error(err))
		return 0, err
	}

	sm.cache[StylePatientInfo] = styleID
	return styleID, nil
}

// GetDateCenterStyle returns style for centered date fields (12pt, center aligned)
func (sm *StyleManager) GetDateCenterStyle() (int, error) {
	if styleID, exists := sm.cache[StyleDateCenter]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 12, Family: FontMyriadPro},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		logger.FromCtx(sm.ctx).Debug("Failed to create date center style", zap.Error(err))
		return 0, err
	}

	sm.cache[StyleDateCenter] = styleID
	return styleID, nil
}

// GetTestResultStyle returns style for test results (13pt, center aligned with borders)
func (sm *StyleManager) GetTestResultStyle() (int, error) {
	if styleID, exists := sm.cache[StyleTestResult]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 13, Family: FontMyriadPro},
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

	sm.cache[StyleTestResult] = styleID
	return styleID, nil
}

// GetTestNameStyle returns style for test names (13pt, left aligned with borders)
func (sm *StyleManager) GetTestNameStyle() (int, error) {
	if styleID, exists := sm.cache[StyleTestName]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 13, Family: FontMyriadPro},
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

	sm.cache[StyleTestName] = styleID
	return styleID, nil
}

func (sm *StyleManager) GetPatientNameLargeCenter() (int, error) {
	if styleID, exists := sm.cache[StylePatientNameLargeCenter]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 15, Family: FontMyriadPro},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: sm.getStandardBorder(),
	})
	if err != nil {
		logger.FromCtx(sm.ctx).Debug("Failed to create test name style", zap.Error(err))
		return 0, err
	}

	sm.cache[StylePatientNameLargeCenter] = styleID
	return styleID, nil
}

// GetAbnormalStyle returns style for abnormal test results (13pt, bold, center aligned with borders)
func (sm *StyleManager) GetAbnormalStyle() (int, error) {
	if styleID, exists := sm.cache[StyleAbnormal]; exists {
		return styleID, nil
	}

	styleID, err := sm.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   13,
			Family: FontMyriadPro,
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

	sm.cache[StyleAbnormal] = styleID
	return styleID, nil
}
