package route

import (
	"example.com/m/v2/src/admin"
	"example.com/m/v2/src/config"
	"example.com/m/v2/src/ctr/user"

	"github.com/kataras/iris/v12"
)

func Route(app *iris.Application) {
	apiPrefix := config.APIConfig.Prefix

	router := app.Party(apiPrefix)
	{
		router.Get("/login", user.Login)
	}

	adminRouter := app.Party(apiPrefix+"/admin", admin.Authentication)
	{
		adminRouter.Get("/user/today", user.TodayRegisterUser)
		adminRouter.Get("/user/yesterday", user.YesterdayRegisterUser)
		adminRouter.Get("/user/latest/30", user.Latest30Day)
		adminRouter.Get("/user/analyze", user.Analyze)
	}
}
