package user

import (
	"encoding/json"
	"example.com/m/v2/src/config"
	"example.com/m/v2/src/ctr/common"
	"example.com/m/v2/src/model/"
	"github.com/kataras/iris/v12"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"net/http"
	"strings"
)

func Login(ctx iris.Context) {

	errJson := common.ErrJson
	code := ctx.FormValue("code")
	if code == "" {
		errJson("code is none", ctx)
		return
	}

	appID := config.WeAppConfig.AppID
	secret := config.WeAppConfig.Secret
	CodeToSessURL := config.WeAppConfig.CodeToSessURL
	CodeToSessURL = strings.Replace(CodeToSessURL, "{appid}", appID, -1)
	CodeToSessURL = strings.Replace(CodeToSessURL, "{secret}", secret, -1)
	CodeToSessURL = strings.Replace(CodeToSessURL, "{code}", code, -1)

	resp, err := http.Get(CodeToSessURL)

	if err != nil {
		fmt.Println(err.Error())
		errJson("error!", ctx)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errJson("error!", ctx)
		return
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
		errJson("error", ctx)
		return
	}

	if _, ok := data["session_key"]; !ok {
		fmt.Println("session_key is none")
		fmt.Println(data)
		errJson("error", ctx)
		return
	}

	var openID string
	var sessionKey string
	openID = data["openid"].(string)
	sessionKey = data["session_key"].(string)
	session := model.MySession.Start(ctx)
	session.Set("weAppOpenID", openID)
	session.Set("weAppSessionKey", sessionKey)

	res := iris.Map{}
	res[config.ServerConfig.SessionID] = session.ID()
	ctx.JSON(iris.Map{
		"status": iris.StatusOK,
		"errNo":  model.ErrorCode.SUCCESS,
		"msg":    "success",
		"data":   res,
	})

}
