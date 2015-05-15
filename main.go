package main

import (
	"fmt"
	"github.com/yosupo06/runner/program/app"
	"github.com/yosupo06/runner/program/config"
	"github.com/yosupo06/runner/program/sample"
	"net/http"
)

func main() {
	http.Handle("/css/", http.StripPrefix("/css/",
		http.FileServer(http.Dir(config.BasePath+"views/css/"))))
	fmt.Println(config.BasePath + "views/css/")
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
	http.HandleFunc("/sample/ranking.html", sample.Ranking)
	http.ListenAndServe(":55001", nil)
}
