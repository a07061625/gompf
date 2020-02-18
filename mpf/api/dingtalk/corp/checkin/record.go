package checkin

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取部门用户签到记录
type record struct {
    dingtalk.BaseCorp
    corpId       string
    agentTag     string
    atType       string
    departmentId int // 部门id
    startTime    int // 开始时间
    endTime      int // 结束时间
}

func (r *record) SetDepartmentId(departmentId int) {
    if departmentId > 0 {
        r.departmentId = departmentId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不合法", nil))
    }
}

func (r *record) SetStartAndEndTime(startTime, endTime int) {
    if startTime < 946656000 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不合法", nil))
    } else if endTime < startTime {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能小于开始时间", nil))
    } else if (endTime - startTime) > 3888000 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能超过开始时间45天", nil))
    }
    r.startTime = startTime
    r.endTime = endTime
}

func (r *record) SetOffset(offset int) {
    if offset >= 0 {
        r.ReqData["offset"] = strconv.Itoa(offset)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (r *record) SetSize(size int) {
    if (size > 0) && (size <= 100) {
        r.ReqData["size"] = strconv.Itoa(size)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (r *record) SetOrder(order string) {
    if (order == "desc") || (order == "asc") {
        r.ReqData["order"] = order
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "排序不合法", nil))
    }
}

func (r *record) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if r.departmentId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不能为空", nil))
    }
    if r.startTime <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不能为空", nil))
    }
    r.ReqData["department_id"] = strconv.Itoa(r.departmentId)
    r.ReqData["start_time"] = strconv.Itoa(r.startTime)
    r.ReqData["end_time"] = strconv.Itoa(r.endTime)
    r.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(r.corpId, r.agentTag, r.atType)

    r.ReqURI = dingtalk.UrlService + "/checkin/record?" + mpf.HTTPCreateParams(r.ReqData, "none", 1)

    return r.GetRequest()
}

func NewRecord(corpId, agentTag, atType string) *record {
    r := &record{dingtalk.NewCorp(), "", "", "", 0, 0, 0}
    r.corpId = corpId
    r.agentTag = agentTag
    r.atType = atType
    r.ReqData["offset"] = "0"
    r.ReqData["size"] = "10"
    r.ReqData["order"] = "desc"
    return r
}
