package sample

import (
	"github.com/yosupo06/runner/program/auth"
	"github.com/yosupo06/runner/program/config"
	"github.com/yosupo06/runner/program/sample/rank"
	"html/template"
	"net/http"
)

var (
	viewPath = config.BasePath + "views/app/sample/"
)

var tp = make(map[string]*template.Template)

func init() {
	f := [...]string{"problem", "chat", "ranking"}
	for _, s := range f {
		t, err := template.ParseFiles(viewPath + s + ".html")
		if err != nil {
			panic(err)
		}
		tp[s] = t
	}
}

func Problem(rw http.ResponseWriter, req *http.Request) {
	u, _ := auth.GetCookie(req)
	tp["problem"].Execute(rw, u)
}

func Chat(rw http.ResponseWriter, req *http.Request) {
	c := GetChat()
	tp["chat"].Execute(rw, c)
}

func Ranking(rw http.ResponseWriter, req *http.Request) {
	r := rank.GetRanking()
	tp["ranking"].Execute(rw, r)
}
