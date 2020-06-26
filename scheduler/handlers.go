package scheduler

import (
	"github.com/julienschmidt/httprouter"
	"github.com/video_server/scheduler/db"
	"log"
	"net/http"
)

func videoDelHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params)  {
	videoId := p.ByName("vid-id")
	if len(videoId) == 0 {
		sendResponse(w, http.StatusBadRequest, "video id should not be empty")
		return
	}

	err := db.AddVideoDeletionRecord(videoId)
	if err != nil {
		sendResponse(w, http.StatusInternalServerError, "Internal server error")
		log.Printf("db add video deletion record err: %s", err.Error())
		return
	}

	sendResponse(w, http.StatusOK, "")
	return
}
