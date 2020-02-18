package attendance

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取用户考勤组
type userGroup struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    userId   string // 员工用户ID
}

func (ug *userGroup) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ug.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工用户ID不合法", nil))
    }
}

func (ug *userGroup) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ug.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工用户ID不能为空", nil))
    }
    ug.ExtendData["userid"] = ug.userId

    ug.ReqUrl = dingtalk.UrlService + "/topapi/attendance/getusergroup?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(ug.corpId, ug.agentTag)

    reqBody := mpf.JsonMarshal(ug.ExtendData)
    client, req := ug.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewUserGroup(corpId, agentTag string) *userGroup {
    ug := &userGroup{dingtalk.NewCorp(), "", "", ""}
    ug.corpId = corpId
    ug.agentTag = agentTag
    ug.ReqContentType = project.HTTPContentTypeJSON
    ug.ReqMethod = fasthttp.MethodPost
    return ug
}
