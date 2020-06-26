package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/video_server/api/session"
	"net/http"
)

type middlewareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middlewareHandler{}
	m.r = r
	return m
}

func (m middlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check session
	validateUserSession(r)
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)
	router.POST("/user/:username", Login)
	router.GET("/user/:username", GetUserInfo)

	router.POST("/user/:username/videos", AddNewVideo)
	router.GET("/user/:username/videos", ListAllVideos)
	router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)

	router.POST("/videos/:vid-id/comments", PostComment)
	router.GET("/videos/:vid-id/comments", ShowComments)

	return router
}

func Prepare()  {
	session.LoadSessionFromDB()
}

// main -> middleware -> modes(err, message) -> handlers -> db -> response
func main() {
	Prepare()
	router := RegisterHandlers()
	m := NewMiddleWareHandler(router)
	http.ListenAndServe(":8000", m)
}