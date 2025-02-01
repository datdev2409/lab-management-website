package handlers

import (
	"context"
	"github.com/a-h/templ"
	"log"
	"net/http"
)

func Render(ctx context.Context, w http.ResponseWriter, comp templ.Component) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return comp.Render(ctx, w)
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
