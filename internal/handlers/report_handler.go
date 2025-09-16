package handlers

import (
	"net/http"

	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
)

// HandleReportPage serves the revenue report page
func (h *Handler) HandleReportPage(w http.ResponseWriter, r *http.Request) error {
	return Render(r.Context(), w, pages.ReportPage())
}
