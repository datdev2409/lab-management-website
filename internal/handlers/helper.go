package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	. "maragu.dev/gomponents"
)

func Render(ctx context.Context, w http.ResponseWriter, comp templ.Component) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return comp.Render(ctx, w)
}

func RenderMultiComponents(ctx context.Context, w http.ResponseWriter, comps []templ.Component) error {
	strBuffer := bytes.NewBufferString("")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, comp := range comps {
		comp.Render(ctx, strBuffer)
	}
	_, err := w.Write(strBuffer.Bytes())
	return err
}

func RenderOOB(ctx context.Context, w http.ResponseWriter, nodes []Node) error {
	return Group(nodes).Render(w)
}

type HandlerFuncReturnError = func(w http.ResponseWriter, r *http.Request) error

func Make(fn HandlerFuncReturnError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}
	}
}

func HTMXRedirect(w http.ResponseWriter, path string) {
	w.Header().Set("HX-Redirect", path)
	w.WriteHeader(http.StatusFound)
}

func SetFlashCookie(w http.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:   "flash",
		Value:  value,
		Path:   "/",
		MaxAge: 5,
	})
}

func GetAndDeleteFlashCookie(w http.ResponseWriter, r *http.Request) string {
	var value string
	if cookie, err := r.Cookie("flash"); err == nil {
		value = cookie.Value
		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "flash",
			Path:   "/",
			MaxAge: -1,
		})
	}
	return value
}

func ParseInputName(name string, sep string) (string, string) {
	parts := strings.SplitN(name, sep, 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func SafeAccessSliceIndex(slice []string, index int) string {
	if index < 0 || index >= len(slice) {
		return ""
	}
	return slice[index]
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
