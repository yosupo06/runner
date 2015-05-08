package app

import (
	"fmt"
	"github.com/yosupo06/runner/program/auth"
	"github.com/yosupo06/runner/program/rank"
	"go/build"
	"html/template"
	"net/http"
	"time"
)

var (
	basePath = build.Default.GOPATH + "/src/github.com/yosupo06/runner/"
	viewPath = basePath + "views/app/"
)

var tp = make(map[string]*template.Template)

func init() {

}

func Index(rw http.ResponseWriter, req *http.Request) {
	tp, err := template.ParseFiles(viewPath + "index.html")
	if err != nil {
		panic(err)
	}
	u, _ := auth.GetCookie(req)
	err = tp.Execute(rw, u)
}

func Register(rw http.ResponseWriter, req *http.Request) {
	tp, _ := template.ParseFiles(viewPath + "register.html")
	if req.PostFormValue("submit") != "" {
		id := req.PostFormValue("id")
		pass := req.PostFormValue("pass")
		if id == "" {
			tp.Execute(rw, map[string]string{"Error": "IDが空"})
			return
		}
		if pass == "" {
			tp.Execute(rw, map[string]string{"Error": "Passが空"})
			return
		}
		err := auth.AddUser(id, pass)
		if err != nil {
			tp.Execute(rw, map[string]string{"Error": err.Error()})
			return
		}
		auth.SetCookie(rw, id)
		http.Redirect(rw, req, "/index.html", http.StatusFound)
	}
	tp.Execute(rw, nil)
}

func Login(rw http.ResponseWriter, req *http.Request) {
	tp, _ := template.ParseFiles(viewPath + "login.html")
	if req.PostFormValue("submit") != "" {
		id := req.PostFormValue("id")
		pass := req.PostFormValue("pass")
		if !auth.AuthPass(id, pass) {
			tp.Execute(rw, map[string]string{"Error": "認証失敗"})
			return
		}
		auth.SetCookie(rw, id)
		http.Redirect(rw, req, "/index.html", http.StatusFound)
	}
	tp.Execute(rw, nil)
}

func Logout(rw http.ResponseWriter, req *http.Request) {
	auth.DelCookie(rw)
	http.Redirect(rw, req, "/index.html", http.StatusFound)
}

var Start time.Time

func Problem(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Before(Start) {
		http.Error(rw, "まだコンテスト開始前です", http.StatusForbidden)
		return
	}
	tp, _ := template.ParseFiles(viewPath + "problem.html")
	u, _ := auth.GetCookie(req)
	tp.Execute(rw, u)
}

func Ranking(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Before(Start) {
		NotFound(rw, req)
		return
	}
	tp, _ := template.ParseFiles(viewPath + "ranking.html")
	r := rank.GetRanking()
	tp.Execute(rw, r)
}

func History(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Before(Start) {
		NotFound(rw, req)
		return
	}
	tp, _ := template.ParseFiles(viewPath + "history.html")
	h := GetHistory()
	tp.Execute(rw, h)
}

func NotFound(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "404")
}
