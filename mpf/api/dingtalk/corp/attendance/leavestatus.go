package attendance

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

// 查询请假状态
type leaveStatus struct {
    dingtalk.BaseCorp
    corpId     string
    agentTag   string
    userIdList []string // 用户列表
    startTime  int      // 开始时间
    endTime    int      // 结束时间
}

func (ls *leaveStatus) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表不能为空", nil))
    } else if len(userList) > 100 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户不能超过100个", nil))
    }

    ls.userIdList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            ls.userIdList = append(ls.userIdList, v)
        }
    }

    if len(ls.userIdList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表不能为空", nil))
    }
}

func (ls *leaveStatus) SetStartAndEndTime(startTime, endTime int) {
    if startTime < 946656000 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不合法", nil))
    } else if endTime < startTime {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能小于开始时间", nil))
    } else if (endTime - startTime) > 15552000 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能超过开始时间180天", nil))
    }
    ls.startTime = startTime
    ls.endTime = endTime
}

func (ls *leaveStatus) SetOffset(offset int) {
    if offset >= 0 {
        ls.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (ls *leaveStatus) SetSize(size int) {
    if size > 0 {
        if size > 20 {
            ls.ExtendData["size"] = 20
        } else {
            ls.ExtendData["size"] = size
        }
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (ls *leaveStatus) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ls.userIdList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表不能为空", nil))
    }
    if ls.startTime <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不能为空", nil))
    }
    ls.ExtendData["userid_list"] = strings.Join(ls.userIdList, ",")
    ls.ExtendData["start_time"] = ls.startTime
    ls.ExtendData["end_time"] = ls.endTime

    ls.ReqUrl = dingtalk.UrlService + "/topapi/attendance/getleavestatus?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(ls.corpId, ls.agentTag)

    reqBody := mpf.JsonMarshal(ls.ExtendData)
    client, req := ls.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewLeaveStatus(corpId, agentTag string) *leaveStatus {
    ls := &leaveStatus{dingtalk.NewCorp(), "", "", make([]string, 0), 0, 0}
    ls.corpId = corpId
    ls.agentTag = agentTag
    ls.ExtendData["offset"] = 0
    ls.ExtendData["size"] = 10
    ls.ReqContentType = project.HTTPContentTypeJSON
    ls.ReqMethod = fasthttp.MethodPost
    return ls
}
