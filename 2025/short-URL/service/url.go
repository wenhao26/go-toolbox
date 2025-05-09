package service

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"toolbox/2025/short-URL/store"
	"toolbox/2025/short-URL/utils"
)

// CreateShortURL 创建短链URL
func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LongURL string `json:"long_url"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	// 生成短链
	shortURL := utils.GenerateShortURL(req.LongURL)

	// 存储
	err := store.SaveURL(shortURL, req.LongURL)
	if err != nil {
		http.Error(w, "保存URL失败", http.StatusInternalServerError)
		return
	}

	// 返回短链
	_ = json.NewEncoder(w).Encode(map[string]string{
		"short_url": "http://localhost:8080/" + shortURL,
	})
}

// RedirectURL 重定向URL
func RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shorturl")

	longURL, err := store.GetURL(shortURL)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}
