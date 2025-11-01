package sheets

import (
	"context"
	"fmt"

	"github.com/datdev2409/lab-admin-go/internal/logger"
	"github.com/xuri/excelize/v2"
)

const FontMyriadPro = "MyRIAD PRO"

// StyleName represents available style names
type StyleName int

const (
	LabNameStyle StyleName = iota
	LabAddressStyle
	LabContactStyle
	ReportNameStyle
	ReportDateStyle
	PatientInfoStyle
	PatientNameStyle
	PatientNameTrackingPageStyle
	TestTableHeaderStyle
	TotalPriceLabelStyle
	TotalPriceStyle
	TestIndexStyle
	TestNameStyle
	TestQuantityStyle
	TestPriceStyle
	TestResultStyle
	TestAbnormalResultStyle
	TestUnitStyle
	TestNormalRangeStyle
	SignatureStyle
	Font12BoldStyle
	Font16BoldCenterStyle
	TrackingTableHeaderStyle
	LocationDateStyle
	LabDepartmentStyle
	SignatureNameStyle
	LabNameLeftStyle
	LabAddressLeftStyle
	PatientNameCenterStyle
	TrackingTableHeaderCyanStyle
)



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

	font11 := &excelize.Font{Size: 11, Family: FontMyriadPro}
	font12 := &excelize.Font{Size: 12, Family: FontMyriadPro}
	font13 := &excelize.Font{Size: 13, Family: FontMyriadPro}
	font11Bold := &excelize.Font{Size: 11, Family: FontMyriadPro, Bold: true}
	font12Bold := &excelize.Font{Size: 12, Family: FontMyriadPro, Bold: true}

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
		PatientNameTrackingPageStyle: {
			Font:      font12Bold,
			Alignment: alignCenter,
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
		TestResultStyle: {
			Font:      font12,
			Alignment: alignCenter,
			Border:    border,
		},
		TestAbnormalResultStyle: {
			Font:      font12Bold,
			Alignment: alignCenter,
			Border:    border,
		},
		TestUnitStyle: {
			Font:      font11,
			Alignment: alignCenter,
			Border:    border,
		},
		TestNormalRangeStyle: {
			Font:      font11,
			Alignment: alignCenter,
			Border:    border,
		},
		LabContactStyle: {
			Font:      font13,
			Alignment: alignLeft,
		},
		SignatureStyle: {
			Font:      font12Bold,
			Alignment: alignCenter,
		},
		Font12BoldStyle: {
			Font:      font12Bold,
			Alignment: alignCenter,
		},
		Font16BoldCenterStyle: {
			Font:      &excelize.Font{Size: 16, Family: FontMyriadPro, Bold: true, Color: "3366FF"},
			Alignment: alignCenter,
		},
		TrackingTableHeaderStyle: {
			Font:      font12Bold,
			Alignment: alignCenter,
			Border:    border,
		},
		LocationDateStyle: {
			Font:      &excelize.Font{Size: 12, Family: FontMyriadPro, Italic: true},
			Alignment: alignCenter,
		},
		LabDepartmentStyle: {
			Font:      font12Bold,
			Alignment: alignCenter,
		},
		SignatureNameStyle: {
			Font:      font12Bold,
			Alignment: alignCenter,
		},
		LabNameLeftStyle: {
			Font:      &excelize.Font{Size: 10, Family: FontMyriadPro, Bold: true, Color: "3366FF"},
			Alignment: alignLeft,
		},
		LabAddressLeftStyle: {
			Font:      &excelize.Font{Size: 10, Family: FontMyriadPro, Color: "3366FF"},
			Alignment: alignLeft,
		},
		PatientNameCenterStyle: {
			Font:      &excelize.Font{Size: 14, Family: FontMyriadPro, Bold: true},
			Alignment: alignCenter,
		},
		TrackingTableHeaderCyanStyle: {
			Font:      &excelize.Font{Size: 11, Family: FontMyriadPro, Bold: true},
			Alignment: alignCenter,
			Border:    border,
			Fill: excelize.Fill{
				Type:    "pattern",
				Pattern: 1,
				Color:   []string{"CCFFFF"}, // Light cyan color
			},
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
		logger.FromCtx(sm.ctx).Error(fmt.Sprintf("style %v not found", styleName))
		return -1
	}

	styleId, err := sm.file.NewStyle(style)
	if err != nil {
		logger.FromCtx(sm.ctx).Error(fmt.Sprintf("can not create style %v", styleName))
		return -1
	}

	sm.cache[styleName] = styleId
	return styleId
}


