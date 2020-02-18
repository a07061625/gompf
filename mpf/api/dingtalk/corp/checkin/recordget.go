package checkin

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

// 获取用户签到记录
type recordGet struct {
    dingtalk.BaseCorp
    corpId    string
    agentTag  string
    atType    string
    userList  []string // 用户列表
    startTime int      // 开始时间
    endTime   int      // 结束时间
}

func (rg *recordGet) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表不能为空", nil))
    } else if len(userList) > 10 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户不能超过10个", nil))
    }

    rg.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            rg.userList = append(rg.userList, v)
        }
    }

    if len(rg.userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表不能为空", nil))
    }
}

func (rg *recordGet) SetStartAndEndTime(startTime, endTime int) {
    if startTime < 946656000 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不合法", nil))
    } else if endTime < startTime {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能小于开始时间", nil))
    } else if (endTime - startTime) > 864000 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能超过开始时间10天", nil))
    }
    rg.startTime = startTime
    rg.endTime = endTime
}

func (rg *recordGet) SetCursor(cursor int) {
    if cursor >= 0 {
        rg.ExtendData["cursor"] = cursor
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页游标不合法", nil))
    }
}

func (rg *recordGet) SetSize(size int) {
    if (size > 0) && (size <= 100) {
        rg.ExtendData["size"] = size
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (rg *recordGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rg.userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表不能为空", nil))
    }
    if rg.startTime <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不能为空", nil))
    }
    rg.ExtendData["userid_list"] = strings.Join(rg.userList, ",")
    rg.ExtendData["start_time"] = rg.startTime
    rg.ExtendData["end_time"] = rg.endTime

    rg.ReqURI = dingtalk.UrlService + "/topapi/checkin/record/get?access_token=" + dingtalk.NewUtil().GetAccessToken(rg.corpId, rg.agentTag, rg.atType)

    reqBody := mpf.JSONMarshal(rg.ExtendData)
    client, req := rg.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewRecordGet(corpId, agentTag, atType string) *recordGet {
    rg := &recordGet{dingtalk.NewCorp(), "", "", "", make([]string, 0), 0, 0}
    rg.corpId = corpId
    rg.agentTag = agentTag
    rg.atType = atType
    rg.ExtendData["cursor"] = 0
    rg.ExtendData["size"] = 10
    rg.ReqContentType = project.HTTPContentTypeJSON
    rg.ReqMethod = fasthttp.MethodPost
    return rg
}
