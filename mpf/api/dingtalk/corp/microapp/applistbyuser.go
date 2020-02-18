package microapp

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取员工可见的应用列表
type appListByUser struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户ID
}

func (alu *appListByUser) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        alu.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (alu *appListByUser) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(alu.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
    }
    alu.ReqData["userid"] = alu.userId
    alu.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(alu.corpId, alu.agentTag, alu.atType)
    alu.ReqUrl = dingtalk.UrlService + "/microapp/list_by_userid?" + mpf.HTTPCreateParams(alu.ReqData, "none", 1)

    return alu.GetRequest()
}

func NewAppListByUser(corpId, agentTag, atType string) *appListByUser {
    alu := &appListByUser{dingtalk.NewCorp(), "", "", "", ""}
    alu.corpId = corpId
    alu.agentTag = agentTag
    alu.atType = atType
    return alu
}
