package route

import (
	"example.com/m/v2/src/config"
	"github.com/kataras/iris/v12"
)

func Route(app *iris.Application)  {
	apiPrefix := config.APIConfig.Prefix

	router := app.Party(apiPrefix)
	{

	}
}
