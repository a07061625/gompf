package process

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取用户待审批数量
type todoNumGet struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    userId   string // 用户id
}

func (ig *todoNumGet) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ig.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不合法", nil))
    }
}

func (ig *todoNumGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ig.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不能为空", nil))
    }
    ig.ExtendData["userid"] = ig.userId

    ig.ReqUrl = dingtalk.UrlService + "/topapi/process/gettodonum?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(ig.corpId, ig.agentTag)

    reqBody := mpf.JSONMarshal(ig.ExtendData)
    client, req := ig.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewTodoNumGet(corpId, agentTag string) *todoNumGet {
    ig := &todoNumGet{dingtalk.NewCorp(), "", "", ""}
    ig.corpId = corpId
    ig.agentTag = agentTag
    ig.ReqContentType = project.HTTPContentTypeJSON
    ig.ReqMethod = fasthttp.MethodPost
    return ig
}
