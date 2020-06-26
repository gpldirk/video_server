package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"video_server/config"
)

// 请求转发：proxy / api
// 获取后端数据，渲染前端页面输出
// 使用http client代理转发http request

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

func request(apiBody *APIBody, w http.ResponseWriter, r *http.Request) {
	var reps *http.Response
	var err error

	url, _ := url.Parse(apiBody.URL)
	url.Host = config.GetLBAddr() + ":" + url.Port()
	newURL := url.String()

	switch apiBody.Method {
	case http.MethodGet:
		req, _ := http.NewRequest("GET", newURL, nil)
		req.Header = r.Header
		reps, err = httpClient.Do(req)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		normalResponse(w, reps)

	case http.MethodPost:
		req, _ := http.NewRequest("POST", newURL, bytes.NewBuffer([]byte(apiBody.ReqBody)))
		req.Header = r.Header
		reps, err = httpClient.Do(req)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		normalResponse(w, reps)

	case http.MethodDelete:
		req, _ := http.NewRequest("DELETE", newURL, nil)
		req.Header = r.Header
		reps, err := httpClient.Do(req)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		normalResponse(w, reps)
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad api request")
	}
}

func normalResponse(w http.ResponseWriter, r *http.Response) {
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		reps, _ := json.Marshal(ErrorInternalError)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, string(reps))
		return
	}

	w.WriteHeader(r.StatusCode)
	io.WriteString(w, string(res))
}


