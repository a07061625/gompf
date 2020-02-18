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

// 获取管理员通讯录权限范围
type adminScope struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户ID
}

func (as *adminScope) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        as.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (as *adminScope) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(as.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
    }
    as.ExtendData["userid"] = as.userId

    as.ReqURI = dingtalk.UrlService + "/topapi/user/get_admin_scope?access_token=" + dingtalk.NewUtil().GetAccessToken(as.corpId, as.agentTag, as.atType)

    reqBody := mpf.JSONMarshal(as.ExtendData)
    client, req := as.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewAdminScope(corpId, agentTag, atType string) *adminScope {
    as := &adminScope{dingtalk.NewCorp(), "", "", "", ""}
    as.corpId = corpId
    as.agentTag = agentTag
    as.atType = atType
    as.ReqContentType = project.HTTPContentTypeJSON
    as.ReqMethod = fasthttp.MethodPost
    return as
}
