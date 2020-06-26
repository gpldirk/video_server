package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/video_server/api/db"
	"github.com/video_server/api/model"
	"github.com/video_server/api/session"
	"github.com/video_server/api/utils"
	"io/ioutil"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	userBody := &model.UserCredential{}
	if err := json.Unmarshal(res, userBody); err != nil {
		sendErrorResponse(w, model.ErrorRequestBodyParseFailed)
		return
	}
	if err := db.AddUserCredential(userBody.Username, userBody.Password); err != nil {
		sendErrorResponse(w, model.ErrorDBError)
		return
	}

	sessionId := session.GenerateNewSession(userBody.Username)
	msg := &model.SignUp{
		Success:   true,
		SessionId: sessionId,
	}

	if reps, err := json.Marshal(msg); err != nil {
		sendErrorResponse(w, model.ErrorInternalError)
		return
	} else {
		sendNormalResponse(w, string(reps), 201)
	}
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	userBody := &model.UserCredential{}
	if err := json.Unmarshal(res, userBody); err != nil {
		sendErrorResponse(w, model.ErrorRequestBodyParseFailed)
		return
	}

	username := p.ByName("username")
	if username != userBody.Username {
		sendErrorResponse(w, model.ErrorUserAuthFailed)
		return
	}

	password, err := db.GetUserCredential(username)
	if err != nil || len(password) == 0 || password != userBody.Password {
		sendErrorResponse(w, model.ErrorUserAuthFailed)
		return
	}

	sessionId := session.GenerateNewSession(username)
	msg := &model.SignIn{
		Success:   true,
		SessionId: sessionId,
	}
	if reps, err := json.Marshal(msg); err != nil {
		sendErrorResponse(w, model.ErrorInternalError)
		return
	} else {
		sendNormalResponse(w, string(reps), 200)
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		return
	}

	username := p.ByName("username")
	user, err := db.GetUser(username)
	if err != nil {
		sendErrorResponse(w, model.ErrorDBError)
		return
	}

	userInfo := &model.UserInfo{Id: user.Id}
	if reps, err := json.Marshal(userInfo); err != nil {
		sendErrorResponse(w, model.ErrorInternalError)
		return
	} else {
		sendNormalResponse(w, string(reps), 200)
	}
}

func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	newVideoBody := &model.NewVideo{}
	if err := json.Unmarshal(res, newVideoBody); err != nil {
		sendErrorResponse(w, model.ErrorRequestBodyParseFailed)
		return
	}

	newVideo, err := db.AddNewVideo(newVideoBody.AuthorId, newVideoBody.Name)
	if err != nil {
		sendErrorResponse(w, model.ErrorDBError)
		return
	}

	if reps, err := json.Marshal(newVideo); err != nil {
		sendErrorResponse(w, model.ErrorInternalError)
		return
	} else {
		sendNormalResponse(w, string(reps), 201)
	}
}

func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		return
	}

	username := p.ByName("username")
	videos, err := db.ListVideoInfo(username, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		sendErrorResponse(w, model.ErrorDBError)
		return
	}

	videosInfo := &model.VideosInfo{Videos: videos}
	if reps, err := json.Marshal(videosInfo); err != nil {
		sendErrorResponse(w, model.ErrorInternalError)
		return
	} else {
		sendNormalResponse(w, string(reps), 200)
	}
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		return
	}

	videoId := p.ByName("vid-id")
	err := db.DeleteVideoInfo(videoId)
	if err != nil {
		sendErrorResponse(w, model.ErrorDBError)
		return
	}

	go utils.SendDeleteVideoRequest(videoId)
	sendNormalResponse(w, "", 204)
}

func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		return
	}
	
	res, _ := ioutil.ReadAll(r.Body)
	commentBody := &model.NewComment{}
	if err := json.Unmarshal(res, commentBody); err != nil {
		sendErrorResponse(w, model.ErrorRequestBodyParseFailed)
		return
	}

	videoId := p.ByName("vid-id")
	if err := db.AddNewComment(videoId, commentBody.AuthorId, commentBody.Content); err != nil {
		sendErrorResponse(w, model.ErrorDBError)
		return
	} else {
		sendNormalResponse(w, "ok", 201)
	}
}

func ShowComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(w, r) {
		return
	}
	videoId := p.ByName("vid-id")
	cms, err := db.ListComments(videoId, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		sendErrorResponse(w, model.ErrorDBError)
		return
	}

	comments := &model.Comments{Comments: cms}
	if reps, err := json.Marshal(comments); err != nil {
		sendErrorResponse(w, model.ErrorInternalError)
		return
	} else {
		sendNormalResponse(w, string(reps), 200)
	}
}



