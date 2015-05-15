package sample

import (
	"fmt"
	"github.com/yosupo06/runner/program/auth"
	"github.com/yosupo06/runner/program/sample/rank"
	"net/http"
	"sync"
	"time"
	"unicode/utf8"
)

type ChatD struct {
	Id   string
	Text string
}

var (
	sm = new(sync.Mutex) //chat
	st = make([]ChatD, 0)
)

type TL struct {
	Chat time.Time //chat
	Comm time.Time //comment
	Rank time.Time //ranking
}

var (
	tm = new(sync.Mutex) //time limit
	tl = make(map[string]TL)
)

func init() {
}

func getForm(req *http.Request) (string, string) {
	return req.FormValue("id"),
		req.FormValue("token")
}

func errorApi(rw http.ResponseWriter, req *http.Request, mes string) {
	fmt.Fprintln(rw, "Error")
	fmt.Fprintln(rw, mes)
}

func CommentApi(rw http.ResponseWriter, req *http.Request) {
	n := time.Now()
	id, token := getForm(req)
	if !auth.AuthToken(id, token) {
		errorApi(rw, req, "Invalid Token")
		return
	}
	c := req.FormValue("comment")
	const ML = 1000
	if len(c) > ML {
		errorApi(rw, req, "Text too large")
		return
	}
	if !utf8.ValidString(c) {
		errorApi(rw, req, "Comment must be UTF-8")
		return
	}

	tm.Lock()
	defer tm.Unlock()

	t := tl[id]
	if n.Before(t.Comm) {
		errorApi(rw, req, "TLE")
		return
	}
	rank.ChangeComment(id, req.FormValue("comment"))
	t.Comm = n.Add(time.Second)
	tl[id] = t
	fmt.Fprintln(rw, "Success")
}

func ChatApi(rw http.ResponseWriter, req *http.Request) {
	n := time.Now()
	id, token := getForm(req)
	if !auth.AuthToken(id, token) {
		errorApi(rw, req, "Invalid Token")
		return
	}
	c := req.FormValue("text")
	const ML = 1000
	if len(c) > ML {
		errorApi(rw, req, "Text too large")
		return
	}
	if !utf8.ValidString(c) {
		errorApi(rw, req, "Text must be UTF-8")
		return
	}

	tm.Lock()
	defer tm.Unlock()

	t := tl[id]
	if n.Before(t.Chat) {
		errorApi(rw, req, "TLE")
		return
	}

	sm.Lock()
	st = append(st, ChatD{id, c})
	rank.AddPoint(id, 1)
	sm.Unlock()
	t.Chat = n.Add(time.Second)
	tl[id] = t
	fmt.Fprintln(rw, "Success")
}

func GetChat() []ChatD {
	sm.Lock()
	const ML = 100
	ml := ML
	if len(st) < ml {
		ml = len(st)
	}
	s := make([]ChatD, len(st))
	for i := 0; i < ml; i++ {
		s[i] = st[len(st)-1-i]
	}
	sm.Unlock()
	return s
}
