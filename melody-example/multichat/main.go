package main

import (
	"net/http"

	"github.com/olahol/melody"
)

func main() {
	m := melody.New()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}

		http.ServeFile(w, r, "chan.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
	})

	http.ListenAndServe(":5000", nil)
}