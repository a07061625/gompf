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

// 发起待办
type workRecordUpdate struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    recordId string // 待办ID
    userId   string // 用户ID
}

func (wru *workRecordUpdate) SetRecordId(recordId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, recordId)
    if match {
        wru.recordId = recordId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "待办ID不合法", nil))
    }
}

func (wru *workRecordUpdate) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        wru.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (wru *workRecordUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(wru.recordId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "待办ID不能为空", nil))
    }
    if len(wru.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
    }
    wru.ExtendData["record_id"] = wru.recordId
    wru.ExtendData["userid"] = wru.userId

    wru.ReqUrl = dingtalk.UrlService + "/topapi/workrecord/update?access_token=" + dingtalk.NewUtil().GetAccessToken(wru.corpId, wru.agentTag, wru.atType)

    reqBody := mpf.JsonMarshal(wru.ExtendData)
    client, req := wru.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewWorkRecordUpdate(corpId, agentTag, atType string) *workRecordUpdate {
    wru := &workRecordUpdate{dingtalk.NewCorp(), "", "", "", "", ""}
    wru.corpId = corpId
    wru.agentTag = agentTag
    wru.atType = atType
    wru.ReqContentType = project.HTTPContentTypeJSON
    wru.ReqMethod = fasthttp.MethodPost
    return wru
}
