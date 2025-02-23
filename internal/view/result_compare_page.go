package view

import (
	"strconv"

	"github.com/datdev2409/lab-admin-go/internal/models"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func PatientSuggestionOption(p models.Patient) Node {
	return Div(Class("patient-suggestion-option"),
		Style("cursor: pointer; padding: 6px 12px; z-index: 100; color: black"),
		hx.Trigger("click"), hx.Get("/api/patients/"+p.ID),
		hx.Swap("outerHTML"), hx.Target("#patient-select-input"),
		Textf("%s - %s", p.Name, p.Phone),
	)
}

func PatientSuggestionList(patients []models.Patient, oob bool) Node {
	return Div(ID("cp_patient-suggestion-list"), Class("shadow bordered bg-white position-absolute start-0 end-0"),
		Style("width: 100%; display: flex; flex-direction: column"),
		If(oob, hx.SwapOOB("true")),
		Map(patients, func(p models.Patient) Node {
			return PatientSuggestionOption(p)
		}),
	)
}

func PatientInfo(p *models.Patient, oob bool) Node {
	if p == nil {
		return Div(ID("patient-info"))
	}

	return Div(ID("patient-info"), Class("mt-3"),
		If(oob, hx.SwapOOB("true")),
		Table(Class("table table-bordered"),
			TBody(
				Tr(
					Td(Text("Họ tên: ")),
					Td(Text(p.Name)),
				),
				Tr(
					Td(Text("Số điện thoại: ")),
					Td(Text(p.Phone)),
				),
				Tr(
					Td(Text("Địa chỉ: ")),
					Td(Text(p.Address)),
				),
				Tr(
					Td(Text("Ngày sinh: ")),
					Td(Text(p.YOB)),
				),
			),
		),
	)
}

func PatientSelectInput(patientName string, patientId string) Node {
	return Div(ID("patient-select-input"), Class("form-group"),
		Label(Class("form-label"), For("patient"), Text("Tên khách hàng")),
		Input(Type("text"), Class("form-control"), AutoComplete("off"),
			ID("patient"), Name("patient_name"), Value(patientName),
			hx.Trigger("keyup changed delay:500ms"), hx.Get("/api/patients"),
			hx.Swap("outerHTML"), hx.Target("#cp_patient-suggestion-list"),
		),
		Input(Type("hidden"), ID("patient_id"), Name("patient_id"), Value(patientId)),
	)
}

func RecordList(records []models.Record) Node {
	if len(records) == 0 {
		return Div(Class("alert alert-info"), Text("Không có kết quả xét nghiệm"))
	}

	return Table(Class("table table-bordered table-responsive"),
		THead(
			Tr(
				Th(Text("")),
				Th(Text("Ngày xét nghiệm")),
				Th(Text("Gói xét nghiệm")),
				Th(Text("Số luợng xét nghiệm")),
			),
			TBody(
				Map(records, func(r models.Record) Node {
					return Tr(
						Td(Input(Type("checkbox"), Name("record_id_"+r.ID), Value(r.ID), Checked())),
						Td(Text(r.CreatedAt.Format("2006-01-02"))),
						Td(Text(r.ComboName)),
						Td(Textf("%s xét nghiệm", strconv.Itoa(len(r.TestResults)))),
					)
				}),
			),
		),
	)
}

func CompareResultsPage(props PageProps) Node {
	props.Title = "Sổ theo dõi kết quả xét nghiệm"
	props.Description = "So sánh kết quả xét nghiệm"

	return Page(props,
		Div(
			H3(Class("pt-3"), Text(props.Title)),
			Div(
				Div(Class("row"),
					Div(Class("position-relative"),
						// Input(Type("text"), Class("form-control"), Placeholder("Chọn bệnh nhân"), AutoComplete("off"),
						// 	ID("patient"), Name("patient_name"),
						// 	hx.Trigger("keyup changed delay:500ms"), hx.Get("/api/patients"),
						// 	hx.Swap("outerHTML"), hx.Target("#cp_patient-suggestion-list"),
						// ),
						PatientSelectInput("", ""),
						PatientSuggestionList([]models.Patient{}, false),
					),
					PatientInfo(nil, false),
					// Div(Class("col"),
					// 	Input(Type("date"), ID("start_date"), Placeholder("Ngày bắt đầu"), Class("form-control")),
					// ),
					// Div(Class("col"),
					// 	Input(Type("date"), ID("end_date"), Placeholder("Ngày kết thúc"), Class("form-control")),
					// ),
					// Div(Class("col-auto"),
					// 	Button(Type("submit"), Class("btn btn-primary"), Text("Tìm kiếm")),
					// ),
				),
				Div(Class("form-group mt-3"),
					Label(Class("form-label"), For("date-select"), Text("Chọn thời gian")),
					Select(Class("form-select"), ID("test-select"),
						Option(Value("all_times"), Text("Tất cả")),
						Option(Value("last_2_times"), Text("2 lần gần nhất")),
						Option(Value("last_3_times"), Text("3 lần gần nhất")),
						Option(Value("last_5_times"), Text("5 lần gần nhất")),
						Option(Value("1_month"), Text("Trong 1 tháng")),
						Option(Value("3_months"), Text("Trong 3 tháng")),
						Option(Value("1_year"), Text("Trong 1 năm")),
						Option(Value("custom"), Text("Tùy chỉnh")),
					),
				),

				Div(ID("record-list"), Class("mt-3"),
					RecordList([]models.Record{}),
				),

				Div(
					Button(Class("btn btn-primary mt-3"),
						hx.Get("/api/records"), hx.Include("input#patient_id"), hx.Target("#record-list"),
						Text("Tìm kiếm kết quả xét nghiệm"),
					),
					Button(Class("btn btn-secondary mt-3 ms-2"),
						hx.Post("/api/tracking/export"), hx.Include("#record-list"),
						Text("Xuất báo cáo"),
					),
				),
			),
		),
	)
}
