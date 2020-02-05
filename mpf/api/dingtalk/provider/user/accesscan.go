package user

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type accessCan struct {
    dingtalk.BaseProvider
    corpId string
    appId  string // 应用ID
    userId string // 用户ID
}

func (ac *accessCan) SetAppId(appId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appId)
    if match {
        ac.appId = appId
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "应用ID不合法", nil))
    }
}

func (ac *accessCan) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ac.userId = userId
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "用户ID不合法", nil))
    }
}

func (ac *accessCan) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ac.appId) == 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "应用ID不能为空", nil))
    }
    if len(ac.userId) == 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "用户ID不能为空", nil))
    }
    ac.ReqData["appId"] = ac.appId
    ac.ReqData["userId"] = ac.userId
    ac.ReqData["access_token"] = dingtalk.NewUtil().GetProviderAuthorizeAccessToken(ac.corpId)
    ac.ReqUrl = dingtalk.UrlService + "/user/can_access_microapp?" + mpf.HttpCreateParams(ac.ReqData, "none", 1)

    return ac.GetRequest()
}

func NewAccessCan(corpId string) *accessCan {
    ac := &accessCan{dingtalk.NewProvider(), "", "", ""}
    ac.corpId = corpId
    return ac
}
