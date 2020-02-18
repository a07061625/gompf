/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 12:05
 */
package admin

import (
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建app
type appCreate struct {
    mppush.BaseJPush
    appName        string // 应用名称
    androidPackage string // 应用Android包名
}

func (ac *appCreate) SetAppName(appName string) {
    if len(appName) > 0 {
        ac.appName = appName
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "应用名称不合法", nil))
    }
}

func (ac *appCreate) SetAndroidPackage(androidPackage string) {
    if len(androidPackage) > 0 {
        ac.androidPackage = androidPackage
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "应用Android包名不合法", nil))
    }
}

func (ac *appCreate) SetGroupName(groupName string) {
    ac.ExtendData["group_name"] = strings.TrimSpace(groupName)
}

func (ac *appCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ac.appName) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "应用名称不能为空", nil))
    }
    if len(ac.androidPackage) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "应用Android包名不能为空", nil))
    }
    ac.ExtendData["app_name"] = ac.appName
    ac.ExtendData["android_package"] = ac.androidPackage

    ac.ReqUrl = ac.GetServiceUrl()

    reqBody := mpf.JSONMarshal(ac.ExtendData)
    client, req := ac.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewAppCreate(key string) *appCreate {
    ac := &appCreate{mppush.NewBaseJPush(mppush.JPushServiceDomainAdmin, key, "dev"), "", ""}
    ac.ServiceUri = "/v1/app"
    ac.ExtendData["group_name"] = ""
    ac.ReqContentType = project.HTTPContentTypeJSON
    ac.ReqMethod = fasthttp.MethodPost
    return ac
}
