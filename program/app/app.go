package app

import (
	"fmt"
	"github.com/yosupo06/runner/program/auth"
	"github.com/yosupo06/runner/program/config"
	"github.com/yosupo06/runner/program/rank"
	"html/template"
	"net/http"
	"time"
	"unicode/utf8"
)

var (
	viewPath = config.BasePath + "views/app/"
)

var tp = make(map[string]*template.Template)

func init() {
	f := [...]string{"index", "register", "login", "problem", "ranking"}
	for _, s := range f {
		t, err := template.ParseFiles(viewPath + s + ".html")
		if err != nil {
			panic(err)
		}
		tp[s] = t
	}
}
func Index(rw http.ResponseWriter, req *http.Request) {
	u, _ := auth.GetCookie(req)
	tp["index"].Execute(rw, u)
}

func Register(rw http.ResponseWriter, req *http.Request) {
	t := tp["register"]
	if req.PostFormValue("submit") != "" {
		id := req.PostFormValue("id")
		pass := req.PostFormValue("pass")
		if !utf8.ValidString(id) {
			t.Execute(rw, map[string]string{"Error": "IDはutf-8"})
		}
		if !utf8.ValidString(pass) {
			t.Execute(rw, map[string]string{"Error": "PASSはutf-8"})
		}
		if len(id) > auth.MaxLength {
			t.Execute(rw, map[string]string{"Error": "IDが長すぎます"})
		}
		if id == "" {
			t.Execute(rw, map[string]string{"Error": "IDが空"})
			return
		}
		if pass == "" {
			t.Execute(rw, map[string]string{"Error": "Passが空"})
			return
		}
		err := auth.AddUser(id, pass)
		if err != nil {
			t.Execute(rw, map[string]string{"Error": "原因不明エラー"})
			return
		}
		auth.SetCookie(rw, id)
		http.Redirect(rw, req, "/index.html", http.StatusFound)
	}
	t.Execute(rw, nil)
}

func Login(rw http.ResponseWriter, req *http.Request) {
	t := tp["login"]
	if req.PostFormValue("submit") != "" {
		id := req.PostFormValue("id")
		pass := req.PostFormValue("pass")
		if !auth.AuthPass(id, pass) {
			t.Execute(rw, map[string]string{"Error": "認証失敗"})
			return
		}
		auth.SetCookie(rw, id)
		http.Redirect(rw, req, "/index.html", http.StatusFound)
	}
	t.Execute(rw, nil)
}

func Logout(rw http.ResponseWriter, req *http.Request) {
	auth.DelCookie(rw)
	http.Redirect(rw, req, "/index.html", http.StatusFound)
}

func Problem(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Before(config.Start) {
		http.Error(rw, "まだコンテスト開始前です", http.StatusForbidden)
		return
	}
	u, _ := auth.GetCookie(req)
	tp["problem"].Execute(rw, u)
}

func Ranking(rw http.ResponseWriter, req *http.Request) {
	if time.Now().Before(config.Start) {
		NotFound(rw, req)
		return
	}
	r := rank.GetRanking()
	tp["ranking"].Execute(rw, r)
}

func NotFound(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "404")
}
