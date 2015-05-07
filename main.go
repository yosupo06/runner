package main

import (
	//	"github.com/yosupo06/runner/program/api"
	"github.com/yosupo06/runner/program/app"
	"net/http"
)

func main() {
	http.HandleFunc("/index.html", app.Index)
	http.HandleFunc("/register.html", app.Register)
	http.HandleFunc("/login.html", app.Login)
	http.HandleFunc("/logout.html", app.Logout)
	//	http.HandleFunc("/comment", api.Comment)
	http.ListenAndServe(":55001", nil)
}
