package attendance

import (
    "regexp"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取请假时长
type leaveApproveDuration struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    userId   string // 员工用户ID
    fromDate string // 开始时间
    toDate   string // 结束时间
}

func (lad *leaveApproveDuration) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        lad.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工用户ID不合法", nil))
    }
}

func (lad *leaveApproveDuration) SetStartAndEndTime(startTime, endTime int) {
    if startTime < 946656000 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不合法", nil))
    } else if endTime < startTime {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能小于开始时间", nil))
    }

    st := time.Unix(int64(startTime), 0)
    et := time.Unix(int64(endTime), 0)
    lad.fromDate = st.Format("2006-01-02 03:04:05")
    lad.toDate = et.Format("2006-01-02 03:04:05")
}

func (lad *leaveApproveDuration) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(lad.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工用户ID不能为空", nil))
    }
    if len(lad.fromDate) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不能为空", nil))
    }
    lad.ExtendData["userid"] = lad.userId
    lad.ExtendData["from_date"] = lad.fromDate
    lad.ExtendData["to_date"] = lad.toDate

    lad.ReqUrl = dingtalk.UrlService + "/topapi/attendance/getleaveapproveduration?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(lad.corpId, lad.agentTag)

    reqBody := mpf.JsonMarshal(lad.ExtendData)
    client, req := lad.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewLeaveApproveDuration(corpId, agentTag string) *leaveApproveDuration {
    lad := &leaveApproveDuration{dingtalk.NewCorp(), "", "", "", "", ""}
    lad.corpId = corpId
    lad.agentTag = agentTag
    lad.ReqContentType = project.HTTPContentTypeJSON
    lad.ReqMethod = fasthttp.MethodPost
    return lad
}
