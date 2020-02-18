/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 1:05
 */
package workrecord

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取用户待办事项
type listByUser struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户ID
    status   int    // 待办状态 0:未完成 1:完成
}

func (lu *listByUser) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        lu.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (lu *listByUser) SetStatus(status int) {
    if (status == 0) || (status == 1) {
        lu.status = status
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "待办状态不合法", nil))
    }
}

func (lu *listByUser) SetOffset(offset int) {
    if offset >= 0 {
        lu.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (lu *listByUser) SetLimit(limit int) {
    if (limit > 0) && (limit <= 50) {
        lu.ExtendData["limit"] = limit
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (lu *listByUser) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(lu.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
    }
    if lu.status < 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "待办状态不能为空", nil))
    }
    lu.ExtendData["userid"] = lu.userId
    lu.ExtendData["status"] = lu.status

    lu.ReqUrl = dingtalk.UrlService + "/topapi/workrecord/getbyuserid?access_token=" + dingtalk.NewUtil().GetAccessToken(lu.corpId, lu.agentTag, lu.atType)

    reqBody := mpf.JSONMarshal(lu.ExtendData)
    client, req := lu.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewListByUser(corpId, agentTag, atType string) *listByUser {
    lu := &listByUser{dingtalk.NewCorp(), "", "", "", "", -1}
    lu.corpId = corpId
    lu.agentTag = agentTag
    lu.atType = atType
    lu.status = -1
    lu.ExtendData["offset"] = 0
    lu.ExtendData["limit"] = 10
    lu.ReqContentType = project.HTTPContentTypeJSON
    lu.ReqMethod = fasthttp.MethodPost
    return lu
}
