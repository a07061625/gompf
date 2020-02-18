package process

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

// 批量获取审批实例id
type instanceIdList struct {
    dingtalk.BaseCorp
    corpId      string
    agentTag    string
    processCode string // 审批码
    startTime   int    // 开始时间
    endTime     int    // 结束时间
}

func (iil *instanceIdList) SetProcessCode(processCode string) {
    if len(processCode) > 0 {
        iil.processCode = processCode
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "审批码不合法", nil))
    }
}

func (iil *instanceIdList) SetStartAndEndTime(startTime, endTime int) {
    if startTime < 946656000 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不合法", nil))
    } else if endTime < startTime {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能小于开始时间", nil))
    }
    iil.startTime = startTime
    iil.endTime = endTime
}

func (iil *instanceIdList) SetCursor(cursor int) {
    if cursor >= 0 {
        iil.ExtendData["cursor"] = cursor
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页游标不合法", nil))
    }
}

func (iil *instanceIdList) SetSize(size int) {
    if (size > 0) && (size <= 10) {
        iil.ExtendData["size"] = size
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "每页记录数不合法", nil))
    }
}

func (iil *instanceIdList) SetUserList(userList []string) {
    if len(userList) > 10 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "发起人不能超过10个", nil))
    }

    users := make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            users = append(users, v)
        }
    }
    iil.ExtendData["userid_list"] = strings.Join(users, ",")
}

func (iil *instanceIdList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(iil.processCode) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "审批码不能为空", nil))
    }
    if iil.startTime <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不能为空", nil))
    }
    iil.ExtendData["process_code"] = iil.processCode
    iil.ExtendData["start_time"] = iil.startTime
    iil.ExtendData["end_time"] = iil.endTime

    iil.ReqUrl = dingtalk.UrlService + "/topapi/processinstance/listids?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(iil.corpId, iil.agentTag)

    reqBody := mpf.JsonMarshal(iil.ExtendData)
    client, req := iil.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewInstanceIdList(corpId, agentTag string) *instanceIdList {
    iil := &instanceIdList{dingtalk.NewCorp(), "", "", "", 0, 0}
    iil.corpId = corpId
    iil.agentTag = agentTag
    iil.ExtendData["cursor"] = 0
    iil.ExtendData["size"] = 10
    iil.ExtendData["userid_list"] = ""
    iil.ReqContentType = project.HTTPContentTypeJSON
    iil.ReqMethod = fasthttp.MethodPost
    return iil
}
