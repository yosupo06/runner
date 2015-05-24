package main

import (
	"github.com/yosupo06/runner/program/app"
	"github.com/yosupo06/runner/program/config"
	"github.com/yosupo06/runner/program/sample"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.Handle("/css/", http.StripPrefix("/css/",
		http.FileServer(http.Dir(config.BasePath+"views/css/"))))
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
	http.HandleFunc("/sample/problem.html", sample.Problem)
	http.HandleFunc("/sample/chat.html", sample.Chat)
	http.HandleFunc("/sample/ranking.html", sample.Ranking)
	http.HandleFunc("/sample/chat", sample.ChatApi)
	http.HandleFunc("/sample/comment", sample.CommentApi)
	http.ListenAndServe(":55001", nil)
}
