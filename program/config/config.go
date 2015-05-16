package config

import (
	"github.com/BurntSushi/toml"
	"go/build"
	"time"
)

var BasePath = build.Default.GOPATH + "/src/github.com/yosupo06/runner/"
var configPath = BasePath + "program/config/"

var c struct {
	Salt  string
	Start time.Time
	End   time.Time
}

var Salt string

var Start, End time.Time

func init() {
	_, err := toml.DecodeFile(configPath+"config.toml", &c)
	if err != nil {
		panic(err)
	}
	Start = c.Start
	End = c.End
	Salt = c.Salt
}
