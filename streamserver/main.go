package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type middlewareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func NewMiddleWareHandler(r *httprouter.Router, limit int) http.Handler {
	m := middlewareHandler{}
	m.r = r
	m.l = NewConnLimiter(limit)
	return m
}

func (m middlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.GetConn() {
		sendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
		return
	}
	defer m.l.ReleaseConn() // lose connection -> release token

	m.r.ServeHTTP(w, r) // request和response透传到内部router
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/videos/:vid-id", streamHandler)
	router.POST("/upload/:vid-id", uploadHandler)
	router.GET("/testpage", testPageHandler)
	return router
}

func main() {
	router := RegisterHandlers()
	m := NewMiddleWareHandler(router, 2)
	http.ListenAndServe(":9000", m)
}

