package main

import (
	"encoding/json"
	"github.com/video_server/api/model"
	"io"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, errReps model.ErrResponse) {
	w.WriteHeader(errReps.HttpSC)
	res, _ := json.Marshal(&errReps.Error)
	io.WriteString(w, string(res))
}

func sendNormalResponse(w http.ResponseWriter, reps string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, reps)
}
