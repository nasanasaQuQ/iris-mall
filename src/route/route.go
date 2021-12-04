package route

import (
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
}
