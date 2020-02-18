package report

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取用户可见的日志模板
type templateListByUser struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
}

func (tlu *templateListByUser) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        tlu.ExtendData["userid"] = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (tlu *templateListByUser) SetOffset(offset int) {
    if offset >= 0 {
        tlu.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (tlu *templateListByUser) SetSize(size int) {
    if (size > 0) && (size <= 100) {
        tlu.ExtendData["size"] = size
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (tlu *templateListByUser) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    tlu.ReqUrl = dingtalk.UrlService + "/topapi/report/template/listbyuserid?access_token=" + dingtalk.NewUtil().GetAccessToken(tlu.corpId, tlu.agentTag, tlu.atType)

    reqBody := mpf.JsonMarshal(tlu.ExtendData)
    client, req := tlu.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewTemplateListByUser(corpId, agentTag, atType string) *templateListByUser {
    tlu := &templateListByUser{dingtalk.NewCorp(), "", "", ""}
    tlu.corpId = corpId
    tlu.agentTag = agentTag
    tlu.atType = atType
    tlu.ExtendData["offset"] = 0
    tlu.ExtendData["size"] = 10
    tlu.ReqContentType = project.HTTPContentTypeJSON
    tlu.ReqMethod = fasthttp.MethodPost
    return tlu
}
