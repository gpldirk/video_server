package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"video_server/config"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

// 当visotor访问当前页面，跳转到register页面
// 当user访问当前页面，跳转到userhome页面
func homeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	username, err1 := r.Cookie("username")
	sessionId, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil { // 如果为visitor，跳转到home
		homePage := &HomePage{Name: "GepeiLu"}
		t, err := template.ParseFiles("./templates/home.html")
		if err != nil {
			log.Printf("Parse html file err: %s", err.Error())
			return
		}

		t.Execute(w, homePage)
		return
	}

	if len(username.Value) != 0 && len(sessionId.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound) // 如果为user，跳转到userhome
		return
	}

	// username and sessionId 是否匹配 -> 前端ajax
}

func userHomeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	username, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")

	// 如果为visitor，跳转到home
	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// 如果为user，跳转到userhome
	fname := r.FormValue("username")
	var userPage *UserPage
	// 先从cookie读取username，然后从提交的表单数据中读取username
	if len(username.Value) != 0 {
		userPage = &UserPage{Name: username.Value}
	} else if len(fname) != 0 {
		userPage = &UserPage{Name: fname}
	}

	t, err := template.ParseFiles("./templates/userhome.html")
	if err != nil {
		log.Printf("Parse html file err: %s", err.Error())
		return
	}

	t.Execute(w, userPage)
}

func apiHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method != http.MethodPost {
		reps, _ := json.Marshal(ErrorBadRequest)
		io.WriteString(w, string(reps))
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apiBody := &APIBody{}
	if err := json.Unmarshal(res, apiBody); err != nil {
		reps, _ := json.Marshal(ErrorReqBodyParseFailed)
		io.WriteString(w, string(reps))
		return
	}

	request(apiBody, w, r)
	defer r.Body.Close()
}

func proxyHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url, _ := url.Parse("http://" + config.GetLBAddr() + ":9000/")
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}

