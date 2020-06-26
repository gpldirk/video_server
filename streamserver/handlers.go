package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("./videos/upload.html")
	t.Execute(w, nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Enter stream handler")
	targetURL := "http://video-server-gepeilu.oss-us-west-1.aliyuncs.com/videos/" + p.ByName("vid-id")
	http.Redirect(w, r, targetURL, http.StatusMovedPermanently)
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		log.Printf("File parse err:%s", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, "File size is too large!")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("File form err: %s", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("File read err: %s", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, "File size is too large!")
		return
	}

	fileName := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR + "/" + fileName, data, 0666)
	if err != nil {
		log.Printf("File write err: %s", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	OSSfn := "videos/" + fileName
	path := VIDEO_DIR + "/" + fileName
	bucketName := "video-server-gepeilu"
	if success := UploadToOSS(OSSfn, path, bucketName); !success {
		sendErrorResponse(w, http.StatusInternalServerError, "upload to oss failed")
		return
	}

	os.Remove(path)
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded successfully!")
}
