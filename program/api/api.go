package api

import (
	"fmt"
	"github.com/yosupo06/runner/program/auth"
	"github.com/yosupo06/runner/program/rank"
	"net/http"
	"time"
)

type TL struct {
	Id   string
	Vote time.Time //sampling
	Comm time.Time //comment
	Hist time.Time //submit
}

var tl = make(map[string]TL)

func getForm(req *http.Request) (string, string) {
	return req.FormValue("id"),
		req.FormValue("token")
}

func Comment(rw http.ResponseWriter, req *http.Request) {
	id, token := getForm(req)
	if !auth.AuthToken(id, token) {
		fmt.Fprintln(rw, "Invalid Token")
		return
	}
	n := time.Now()
	t := tl[id]
	if n.Before(t.Comm) {
		fmt.Fprintln(rw, "TLE")
	}
	rank.ChangeComment(id, req.FormValue("comment"))
	t.Comm = n.Add(time.Second)
	tl[id] = t
}
