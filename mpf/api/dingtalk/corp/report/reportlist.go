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

// 获取用户日志数据
type reportList struct {
    dingtalk.BaseCorp
    corpId       string
    agentTag     string
    atType       string
    startTime    int64  // 开始时间
    endTime      int64  // 结束时间
    templateName string // 模板名称
    userId       string // 用户ID
}

func (rl *reportList) SetStartAndEndTime(startTime, endTime int) {
    if startTime < 946656000 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不合法", nil))
    } else if endTime < startTime {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能小于开始时间", nil))
    }
    rl.startTime = int64(1000 * startTime)
    rl.endTime = int64(1000 * endTime)
}

func (rl *reportList) SetTemplateName(templateName string) {
    if len(templateName) > 0 {
        rl.templateName = templateName
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "模板名称不合法", nil))
    }
}

func (rl *reportList) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        rl.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (rl *reportList) SetCursor(cursor int) {
    if cursor >= 0 {
        rl.ExtendData["cursor"] = cursor
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页游标不合法", nil))
    }
}

func (rl *reportList) SetSize(size int) {
    if (size > 0) && (size <= 20) {
        rl.ExtendData["size"] = size
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (rl *reportList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if rl.startTime <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不能为空", nil))
    }
    if len(rl.templateName) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "模板名称不能为空", nil))
    }
    if len(rl.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
    }
    rl.ExtendData["start_time"] = rl.startTime
    rl.ExtendData["end_time"] = rl.endTime
    rl.ExtendData["template_name"] = rl.templateName
    rl.ExtendData["userid"] = rl.userId

    rl.ReqUrl = dingtalk.UrlService + "/topapi/report/list?access_token=" + dingtalk.NewUtil().GetAccessToken(rl.corpId, rl.agentTag, rl.atType)

    reqBody := mpf.JSONMarshal(rl.ExtendData)
    client, req := rl.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewReportList(corpId, agentTag, atType string) *reportList {
    rl := &reportList{dingtalk.NewCorp(), "", "", "", 0, 0, "", ""}
    rl.corpId = corpId
    rl.agentTag = agentTag
    rl.atType = atType
    rl.ExtendData["cursor"] = 0
    rl.ExtendData["size"] = 10
    rl.ReqContentType = project.HTTPContentTypeJSON
    rl.ReqMethod = fasthttp.MethodPost
    return rl
}
