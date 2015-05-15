package config

import (
	"go/build"
	"time"
)

var Start, End time.Time
var BasePath = build.Default.GOPATH + "/src/github.com/yosupo06/runner/"

func init() {
	var loc, _ = time.LoadLocation("Asia/Tokyo")
	Start, _ = time.ParseInLocation(time.RFC3339,
		"2015-05-24T21:00:00Z",
		loc)
	End, _ = time.ParseInLocation(time.RFC3339,
		"2015-05-24T23:00:00Z",
		loc)
}
