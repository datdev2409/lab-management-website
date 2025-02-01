package handlers

import (
	"context"
	"github.com/datdev2409/lab-admin-go/internal/templates/pages"
	"net/http"
)

func HandleRecordPage(w http.ResponseWriter, r *http.Request) error {
	messages := map[string]string{
		"patient:create:success": "Thêm bệnh nhân thành công",
	}
	redirectCode := GetAndDeleteFlashCookie(w, r)
	return Render(context.Background(), w, pages.RecordPage(messages[redirectCode]))
}
