package scheduler

// user -> api -> video delete -> scheduler -> db delete record
// api -> scheduler -> dispatcher -> read db delete record -> chan -> executor -> delete video from db/folder

import (
	"github.com/julienschmidt/httprouter"
	"github.com/video_server/scheduler/taskrunner"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/video-delete-record/:vid-id", videoDelHandler)
	return router
}

func main() {
	go taskrunner.Start()
	r := RegisterHandlers()
	http.ListenAndServe(":9001", r)
}
