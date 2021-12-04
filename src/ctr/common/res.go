package common

import (
	"example.com/m/v2/src/model"
	"github.com/kataras/iris/v12"
)

func ErrJson(msg string, ctx iris.Context) {
	ctx.JSON(iris.Map{
		"status": iris.StatusOK,
		"errNo":  model.ErrorCode.ERROR,
		"msg":    msg,
		"data":   iris.Map{},
	})
}
