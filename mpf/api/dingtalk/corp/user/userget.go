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

// 获取用户详情
type userGet struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户id
}

func (ug *userGet) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ug.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不合法", nil))
    }
}

func (ug *userGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ug.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不能为空", nil))
    }
    ug.ReqData["userid"] = ug.userId
    ug.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(ug.corpId, ug.agentTag, ug.atType)
    ug.ReqUrl = dingtalk.UrlService + "/user/get?" + mpf.HTTPCreateParams(ug.ReqData, "none", 1)

    return ug.GetRequest()
}

func NewUserGet(corpId, agentTag, atType string) *userGet {
    ug := &userGet{dingtalk.NewCorp(), "", "", "", ""}
    ug.corpId = corpId
    ug.agentTag = agentTag
    ug.atType = atType
    ug.ReqData["lang"] = "zh_CN"
    return ug
}
