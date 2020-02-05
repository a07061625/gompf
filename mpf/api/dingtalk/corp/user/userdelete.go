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

// 删除用户
type userDelete struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户id
}

func (ud *userDelete) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ud.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不合法", nil))
    }
}

func (ud *userDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ud.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不能为空", nil))
    }
    ud.ReqData["userid"] = ud.userId
    ud.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(ud.corpId, ud.agentTag, ud.atType)
    ud.ReqUrl = dingtalk.UrlService + "/user/delete?" + mpf.HttpCreateParams(ud.ReqData, "none", 1)

    return ud.GetRequest()
}

func NewUserDelete(corpId, agentTag, atType string) *userDelete {
    ud := &userDelete{dingtalk.NewCorp(), "", "", "", ""}
    ud.corpId = corpId
    ud.agentTag = agentTag
    ud.atType = atType
    return ud
}
