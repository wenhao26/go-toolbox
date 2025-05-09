package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"toolbox/2025/short-URL/service"
)

func StartServer() {
	r := chi.NewRouter()
	r.Post("/shorten", service.CreateShortURL)
	r.Get("/{shorturl}", service.RedirectURL)
	http.ListenAndServe(":8080", r)
}
