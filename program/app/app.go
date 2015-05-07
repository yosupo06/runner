package app

import (
	"github.com/yosupo06/runner/program/auth"
	"github.com/yosupo06/runner/program/rank"
	"go/build"
	"html/template"
	"net/http"
)

var basePath = build.Default.GOPATH + "/src/github.com/yosupo06/runner/"
var viewPath = basePath + "views/app/"

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

func Problem(rw http.ResponseWriter, req *http.Request) {

}

func Ranking(rw http.ResponseWriter, req *http.Request) {
	tp, _ := template.ParseFiles(viewPath + "ranking.html")
	r := rank.GetRanking()
	tp.Execute(rw, r)
}
