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

// 获取打卡详情
type recordList struct {
    dingtalk.BaseCorp
    corpId     string
    agentTag   string
    userIdList []string // 员工列表
    fromDate   string   // 开始时间
    toDate     string   // 结束时间
}

func (rl *recordList) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    } else if len(userList) > 50 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工不能超过50个", nil))
    }

    rl.userIdList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            rl.userIdList = append(rl.userIdList, v)
        }
    }

    if len(rl.userIdList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    }
}

func (rl *recordList) SetStartAndEndTime(startTime, endTime int) {
    if startTime < 946656000 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不合法", nil))
    } else if endTime < startTime {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能小于开始时间", nil))
    } else if (endTime - startTime) > 604800 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "结束时间不能超过开始时间7天", nil))
    }

    st := time.Unix(int64(startTime), 0)
    et := time.Unix(int64(endTime), 0)
    rl.fromDate = st.Format("2006-01-02 03:04:05")
    rl.toDate = et.Format("2006-01-02 03:04:05")
}

func (rl *recordList) SetIsI18n(isI18n string) {
    if (isI18n == "true") || (isI18n == "false") {
        rl.ExtendData["isI18n"] = isI18n
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "海外使用标识不合法", nil))
    }
}

func (rl *recordList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rl.userIdList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "员工列表不能为空", nil))
    }
    if len(rl.fromDate) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "开始时间不能为空", nil))
    }
    rl.ExtendData["userIds"] = rl.userIdList
    rl.ExtendData["checkDateFrom"] = rl.fromDate
    rl.ExtendData["checkDateTo"] = rl.toDate

    rl.ReqUrl = dingtalk.UrlService + "/attendance/listRecord?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(rl.corpId, rl.agentTag)

    reqBody := mpf.JsonMarshal(rl.ExtendData)
    client, req := rl.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewRecordList(corpId, agentTag string) *recordList {
    rl := &recordList{dingtalk.NewCorp(), "", "", make([]string, 0), "", ""}
    rl.corpId = corpId
    rl.agentTag = agentTag
    rl.ExtendData["isI18n"] = "false"
    rl.ReqContentType = project.HTTPContentTypeJSON
    rl.ReqMethod = fasthttp.MethodPost
    return rl
}
