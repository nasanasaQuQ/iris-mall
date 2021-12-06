package user

import (
	"encoding/json"
	"example.com/m/v2/src/config"
	"example.com/m/v2/src/ctr/common"
	"example.com/m/v2/src/model"
	"example.com/m/v2/src/utils"
	"github.com/kataras/iris/v12"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"net/http"
	"strings"
	"time"
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

func SetAppUserInfo(ctx iris.Context) {
	errJson := common.ErrJson
	type EncryptedUser struct {
		EncryptedData string `json:"encryptedData"`
		IV            string `json:"iv"`
	}
	var appUser EncryptedUser

	if ctx.ReadJSON(&appUser) != nil {
		errJson("参数错误", ctx)
		return
	}

	sessions := model.MySession.Start(ctx)
	sessionKey := sessions.GetString("appSessionKey")
	if sessionKey == "" {
		errJson("session err", ctx)
		return
	}

	userInfoStr, err := utils.DecodeWeAppUserInfo(appUser.EncryptedData, sessionKey, appUser.IV)
	if err != nil {
		fmt.Println(err.Error())
		errJson("err", ctx)
		return
	}

	var user model.User
	if err := json.Unmarshal([]byte(userInfoStr), &user); err != nil {
		errJson("err", ctx)
		return
	}

	sessions.Set("appUser", user)
	ctx.JSON(iris.Map{
		"status": iris.StatusOK,
		"errNo":  model.ErrorCode.SUCCESS,
		"msg":    "success",
		"data":   iris.Map{},
	})

	return
}

func YesterdayRegisterUser(ctx iris.Context) {

	var user model.User
	count := user.YesterdayRegisterUser()
	ctx.JSON(iris.Map{
		"status": iris.StatusOK,
		"errNo":  model.ErrorCode.SUCCESS,
		"msg":    "success",
		"data": iris.Map{
			"count": count,
		},
	})
	return

}

func TodayRegisterUser(ctx iris.Context) {
	var users model.User
	count := users.TodayRegisterUser()

	ctx.JSON(iris.Map{
		"status": iris.StatusOK,
		"errNo":  model.ErrorCode.SUCCESS,
		"msg":    "success",
		"data": iris.Map{
			"count": count,
		},
	})

}

func Latest30Day(ctx iris.Context) {
	var users model.UserPerDay
	result := users.Latest30Day()
	var data iris.Map
	if result != nil {
		data = iris.Map{
			"users": [0]int{},
		}
	} else {
		data = iris.Map{
			"users": result,
		}
	}
	ctx.JSON(iris.Map{
		"status": iris.StatusOK,
		"errNo":  model.ErrorCode.SUCCESS,
		"msg":    "success",
		"data":   data,
	})

}

func Analyze(ctx iris.Context) {
	var user model.User
	now := time.Now()
	nowSec := now.Unix()
	yesterdaySec := nowSec - 24*60*60
	yesterday := time.Unix(yesterdaySec, 0)

	yesterdayCount := user.PurchaseUserByDate(yesterday)
	todayCount := user.PurchaseUserByDate(now)
	yesterdayRegisterCount := user.YesterdayRegisterUser()
	todayRegisterUser := user.TodayRegisterUser()

	data := iris.Map{
		"yesterdayCount":         yesterdayCount,
		"todayCount":             todayCount,
		"yesterdayRegisterCount": yesterdayRegisterCount,
		"todayRegisterUser":      todayRegisterUser,
	}

	ctx.JSON(iris.Map{
		"status": iris.StatusOK,
		"errNo":  model.ErrorCode.SUCCESS,
		"msg":    "success",
		"data":   data,
	})
	return

}
