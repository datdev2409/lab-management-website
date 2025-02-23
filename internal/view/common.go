package view

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type PageProps struct {
	Title       string
	Description string
}

func Page(props PageProps, children ...Node) Node {
	return HTML5(HTML5Props{
		Title:       props.Title,
		Description: props.Description,
		Head: []Node{
			Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css")),
			Script(Src("https://unpkg.com/htmx.org@2.0.4"), Defer()),
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
		},
		Body: []Node{
			Div(Class("container-fluid"),
				Div(Class("row"),
					Sidebar(),
					Div(Class("col-10 position-relative"), Style("height: 100vh; overflow-y: hidden"), Group(children)),
				),
			),
		},
	})
}

func SidebarLink(text string, href string, active bool) Node {
	return Li(Class("nav-item"),
		If(active, Class("active")),
		A(Class("nav-link"), Href(href), Text(text)),
	)
}

func Sidebar() Node {
	return Div(Class("col-2 border shadow-sm"), Style("height: 100vh; min-width: 180px; overflow-y: hidden"),
		Ul(Class("nav flex-column pt-3"),
			SidebarLink("Phiếu xét nghiệm", "/phieu-xet-nghiem", false),
			SidebarLink("Danh mục bệnh nhân", "/danh-muc-benh-nhan", false),
			SidebarLink("Danh mục bác sĩ", "/danh-muc-benh-nhan", false),
			SidebarLink("Danh mục xét nghiệm", "/danh-muc-xet-nghiem", false),
			SidebarLink("Danh mục gói xét nghiệm", "/danh-muc-goi-xet-nghiem", false),
			SidebarLink("Danh mục so sánh", "/danh-muc-so-sanh", false),
			SidebarLink("Sổ theo dõi kết quả", "/so-sanh-ket-qua", true),
		),
	)
}
