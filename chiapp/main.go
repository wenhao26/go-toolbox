package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Heartbeat("/ping"), middleware.RealIP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("chi app start！"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	r.Get("/article/{channel_id}", getArticle)
	r.Get("/post/{id}", getPost)

	_ = http.ListenAndServe(":5000", r)
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	channelIdStr := chi.URLParam(r, "channel_id")

	channelID, _ := strconv.Atoi(channelIdStr)
	if channelID != 1 {
		writeJSON(w, http.StatusOK, response(0, "no data found", nil))
		return
	}

	// mock
	data := []map[string]interface{}{
		{
			"article_id": 10001,
			"title":      "The Lycan Prince’s Puppy",
		},
		{
			"article_id": 10002,
			"title":      "Yes Daddy",
		},
		{
			"article_id": 10003,
			"title":      "Shattered Bonds",
		},
		{
			"article_id": 10004,
			"title":      "Accidental Surrogate for Alpha",
		},
	}
	writeJSON(w, http.StatusOK, response(0, "success", data))
}

func getPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, _ := strconv.Atoi(idStr)

	// mock
	data := map[string]interface{}{
		"article_id": id,
		"content":    "The Lycan Prince’s Puppy The Lycan Prince’s Puppy The Lycan Prince’s Puppy",
	}
	writeJSON(w, http.StatusOK, response(0, "success", data))
}

// writeJSON 输出JSON响应
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// respData 响应数据包
func response(code int, msg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	}
}

// Response 响应结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewResponse 响应构造函数
func NewResponse(opts ...func(*Response)) Response {
	r := Response{
		Code: 1,
		Msg:  "ok",
		Data: nil,
	}
	for _, opt := range opts {
		opt(&r)
	}
	return r
}
