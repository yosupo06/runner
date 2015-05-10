package main

import (
	"github.com/yosupo06/runner/program/app"
	"go/build"
	"net/http"
	"time"
)

var basePath = build.Default.GOPATH + "/src/github.com/yosupo06/runner/"

func main() {
	var loc, _ = time.LoadLocation("Asia/Tokyo")
	app.Start, _ = time.ParseInLocation(time.RFC3339,
		"2015-05-24T21:00:00Z",
		loc)
	app.End, _ = time.ParseInLocation(time.RFC3339,
		"2015-05-24T23:00:00Z",
		loc)
	http.Handle("/css/", http.StripPrefix("/css/",
		http.FileServer(http.Dir(basePath+"views/css/"))))
	http.HandleFunc("/index.html", app.Index)
	http.HandleFunc("/register.html", app.Register)
	http.HandleFunc("/login.html", app.Login)
	http.HandleFunc("/logout.html", app.Logout)
	http.HandleFunc("/problem.html", app.Problem)
	http.HandleFunc("/ranking.html", app.Ranking)
	http.HandleFunc("/vote", app.VoteApi)
	http.HandleFunc("/comment", app.CommentApi)
	http.HandleFunc("/ranking", app.RankingApi)
	http.HandleFunc("/info", app.InfoApi)
	http.HandleFunc("/", app.NotFound)
	http.ListenAndServe(":55001", nil)
}
