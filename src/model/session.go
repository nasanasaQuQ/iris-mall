package model

import (
	"example.com/m/v2/src/config"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

var (
	MySession = *sessions.New(sessions.Config{Cookie: config.ServerConfig.SessionID, Expires: time.Minute * 20})
)
