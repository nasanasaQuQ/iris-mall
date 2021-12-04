package main

import (
	"example.com/m/v2/src/config"
	"example.com/m/v2/src/model"
	"example.com/m/v2/src/route"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris/v12"
	"os"
	"strconv"
)

func init() {
	db, err := gorm.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	if config.DBConfig.SQLLog {
		db.LogMode(true)
	}
	db.DB().SetMaxIdleConns(config.DBConfig.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.DBConfig.MaxOpenConns)

	model.DB = db

}

func main() {

	app := iris.New()

	if config.ServerConfig.Debug {
		app.Use()
	}

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"status": iris.StatusOK,
			"errNo":  model.ErrorCode.NotFound,
			"msg":    "Not Found",
			"data":   iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"status": iris.StatusInternalServerError,
			"errNo":  model.ErrorCode.NotFound,
			"msg":    "Not Found",
			"data":   iris.Map{},
		})
	})

	app.Use()

	// 注册route到app
	route.Route(app)

	app.Run(iris.Addr(":"+strconv.Itoa(config.ServerConfig.Port)), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
		Charset:                           "UTF-8",
	}))

}
