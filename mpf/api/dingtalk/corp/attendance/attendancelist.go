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

// 获取打卡结果
type attendanceList struct {
    dingtalk.BaseCorp
    corpId     string
    agentTag   string
    userIdList []string // 员工列表
    fromDate   string   // 开始时间
    toDate     string   // 结束时间
}

func (al *attendanceList) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    } else if len(userList) > 50 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工不能超过50个", nil))
    }

    al.userIdList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            al.userIdList = append(al.userIdList, v)
        }
    }

    if len(al.userIdList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    }
}

func (al *attendanceList) SetStartAndEndTime(startTime, endTime int) {
    if startTime < 946656000 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不合法", nil))
    } else if endTime < startTime {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能小于开始时间", nil))
    } else if (endTime - startTime) > 604800 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能超过开始时间7天", nil))
    }

    st := time.Unix(int64(startTime), 0)
    et := time.Unix(int64(endTime), 0)
    al.fromDate = st.Format("2006-01-02 03:04:05")
    al.toDate = et.Format("2006-01-02 03:04:05")
}

func (al *attendanceList) SetOffset(offset int) {
    if offset >= 0 {
        al.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (al *attendanceList) SetLimit(limit int) {
    if limit > 0 {
        if limit > 50 {
            al.ExtendData["limit"] = 50
        } else {
            al.ExtendData["limit"] = limit
        }
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (al *attendanceList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(al.userIdList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    }
    if len(al.fromDate) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不能为空", nil))
    }
    al.ExtendData["userIdList"] = al.userIdList
    al.ExtendData["workDateFrom"] = al.fromDate
    al.ExtendData["workDateTo"] = al.toDate

    al.ReqUrl = dingtalk.UrlService + "/attendance/list?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(al.corpId, al.agentTag)

    reqBody := mpf.JsonMarshal(al.ExtendData)
    client, req := al.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewAttendanceList(corpId, agentTag string) *attendanceList {
    al := &attendanceList{dingtalk.NewCorp(), "", "", make([]string, 0), "", ""}
    al.corpId = corpId
    al.agentTag = agentTag
    al.ExtendData["offset"] = 0
    al.ExtendData["limit"] = 10
    al.ReqContentType = project.HTTPContentTypeJSON
    al.ReqMethod = fasthttp.MethodPost
    return al
}
