package main

import (
	"github.com/video_server/api/model"
	"github.com/video_server/api/session"
	"net/http"
)

var HEADER_SESSION = "X-Session-Id"
var HEADER_USERNAME = "X-User-Name"

func validateUserSession(r *http.Request) bool {
	sessionId := r.Header.Get(HEADER_SESSION)
	if len(sessionId) == 0 {
		return false
	}
	username, ok := session.IsSessionExpired(sessionId)
	if ok {
		return false
	}

	r.Header.Add(HEADER_USERNAME, username)
	return true
}

func validateUser(w http.ResponseWriter, r *http.Request) bool {
	username := r.Header.Get(HEADER_USERNAME)
	if len(username) == 0 {
		sendErrorResponse(w, model.ErrorUserAuthFailed)
		return false
	}
	return true
}