package health

import (
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 批量获取钉钉运动数据
type stepInfoListByUser struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userList []string // 员工列表
    statDate string   // 时间
}

func (ilu *stepInfoListByUser) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    } else if len(userList) > 50 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工总数不能超过50个", nil))
    }

    ilu.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            ilu.userList = append(ilu.userList, v)
        }
    }
}

func (ilu *stepInfoListByUser) SetStatDate(statDate string) {
    match, _ := regexp.MatchString(`^2[0-9]{7}$`, statDate)
    if match {
        ilu.statDate = statDate
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "时间不合法", nil))
    }
}

func (ilu *stepInfoListByUser) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ilu.userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    }
    if len(ilu.statDate) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "时间不能为空", nil))
    }
    ilu.ExtendData["userids"] = strings.Join(ilu.userList, ",")
    ilu.ExtendData["stat_date"] = ilu.statDate

    ilu.ReqURI = dingtalk.UrlService + "/topapi/health/stepinfo/listbyuserid?access_token=" + dingtalk.NewUtil().GetAccessToken(ilu.corpId, ilu.agentTag, ilu.atType)

    reqBody := mpf.JSONMarshal(ilu.ExtendData)
    client, req := ilu.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewStepInfoListByUser(corpId, agentTag, atType string) *stepInfoListByUser {
    ilu := &stepInfoListByUser{dingtalk.NewCorp(), "", "", "", make([]string, 0), ""}
    ilu.corpId = corpId
    ilu.agentTag = agentTag
    ilu.atType = atType
    ilu.ReqContentType = project.HTTPContentTypeJSON
    ilu.ReqMethod = fasthttp.MethodPost
    return ilu
}
