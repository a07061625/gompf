/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 11:56
 */
package admin

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除app
type appDelete struct {
    mppush.BaseJPush
    appKey string // 应用标识
}

func (ad *appDelete) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        ad.appKey = appKey
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "应用标识不合法", nil))
    }
}

func (ad *appDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ad.appKey) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "应用标识不能为空", nil))
    }
    ad.ServiceUri = "/v1/app/" + ad.appKey + "/delete"

    ad.ReqUrl = ad.GetServiceUrl()

    return ad.GetRequest()
}

func NewAppDelete(key string) *appDelete {
    ad := &appDelete{mppush.NewBaseJPush(mppush.JPushServiceDomainAdmin, key, "dev"), ""}
    ad.ReqMethod = fasthttp.MethodPost
    return ad
}
