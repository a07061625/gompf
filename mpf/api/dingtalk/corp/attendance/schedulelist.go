package attendance

import (
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 企业考勤排班详情
type scheduleList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    workDate string // 排班时间
}

func (sl *scheduleList) SetWorkDate(workTime int) {
    if workTime > 0 {
        wt := time.Unix(int64(workTime), 0)
        sl.workDate = wt.Format("2006-01-02 03:04:05")
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "排班时间不合法", nil))
    }
}

func (sl *scheduleList) SetOffset(offset int) {
    if offset >= 0 {
        sl.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (sl *scheduleList) SetSize(size int) {
    if size > 0 {
        if size > 200 {
            sl.ExtendData["size"] = 200
        } else {
            sl.ExtendData["size"] = size
        }
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (sl *scheduleList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sl.workDate) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "排班时间不能为空", nil))
    }
    sl.ExtendData["workDate"] = sl.workDate

    sl.ReqUrl = dingtalk.UrlService + "/topapi/attendance/listschedule?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(sl.corpId, sl.agentTag)

    reqBody := mpf.JsonMarshal(sl.ExtendData)
    client, req := sl.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewScheduleList(corpId, agentTag string) *scheduleList {
    sl := &scheduleList{dingtalk.NewCorp(), "", "", ""}
    sl.corpId = corpId
    sl.agentTag = agentTag
    sl.ExtendData["offset"] = 0
    sl.ExtendData["size"] = 10
    sl.ReqContentType = project.HTTPContentTypeJSON
    sl.ReqMethod = fasthttp.MethodPost
    return sl
}
