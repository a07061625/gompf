/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 1:05
 */
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

// 获取用户基础信息
type userInfo struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    authCode string // 授权码
}

func (ui *userInfo) SetAuthCode(authCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, authCode)
    if match {
        ui.authCode = authCode
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "授权码不合法", nil))
    }
}

func (ui *userInfo) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ui.authCode) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "授权码不能为空", nil))
    }
    ui.ReqData["code"] = ui.authCode
    ui.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(ui.corpId, ui.agentTag, ui.atType)
    ui.ReqUrl = dingtalk.UrlService + "/user/getuserinfo?" + mpf.HttpCreateParams(ui.ReqData, "none", 1)

    return ui.GetRequest()
}

func NewUserInfo(corpId, agentTag, atType string) *userInfo {
    ui := &userInfo{dingtalk.NewCorp(), "", "", "", ""}
    ui.corpId = corpId
    ui.agentTag = agentTag
    ui.atType = atType
    return ui
}
