package sso

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取应用管理员的身份信息
type userInfoGet struct {
    dingtalk.BaseCorp
    corpId   string
    authCode string // 授权码
}

func (uig *userInfoGet) SetAuthCode(authCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, authCode)
    if match {
        uig.authCode = authCode
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "授权码不合法", nil))
    }
}

func (uig *userInfoGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(uig.authCode) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "授权码不能为空", nil))
    }
    uig.ReqData["code"] = uig.authCode
    if len(uig.corpId) > 0 {
        uig.ReqData["access_token"] = dingtalk.NewUtil().GetCorpSsoToken(uig.corpId)
    } else {
        uig.ReqData["access_token"] = dingtalk.NewUtil().GetProviderSsoToken()
    }

    uig.ReqURI = dingtalk.UrlService + "/sso/getuserinfo?" + mpf.HTTPCreateParams(uig.ReqData, "none", 1)

    return uig.GetRequest()
}

func NewUserInfoGet(corpId string) *userInfoGet {
    uig := &userInfoGet{dingtalk.NewCorp(), "", ""}
    uig.corpId = corpId
    return uig
}
