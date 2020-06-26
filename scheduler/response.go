package scheduler

import (
	"io"
	"net/http"
)

func sendResponse(w http.ResponseWriter, sc int, reps string) {
	w.WriteHeader(sc)
	io.WriteString(w, reps)
}
