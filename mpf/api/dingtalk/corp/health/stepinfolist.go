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

// 获取个人或部门的钉钉运动数据
type stepInfoList struct {
    dingtalk.BaseCorp
    corpId       string
    agentTag     string
    atType       string
    dataType     int      // 数据类型 0:取用户步数 1:取部门步数
    userId       string   // 用户ID
    departmentId int      // 部门ID
    statDateList []string // 时间列表
}

func (sil *stepInfoList) SetDataType(dataType int) {
    if (dataType == 0) || (dataType == 1) {
        sil.dataType = dataType
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "数据类型不合法", nil))
    }
}

func (sil *stepInfoList) SetDepartmentId(departmentId int) {
    if departmentId > 0 {
        sil.departmentId = departmentId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门ID不合法", nil))
    }
}

func (sil *stepInfoList) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        sil.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (sil *stepInfoList) SetStatDateList(statDateList []string) {
    if len(statDateList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "时间列表不能为空", nil))
    } else if len(statDateList) > 31 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "时间总数不能超过31天", nil))
    }

    sil.statDateList = make([]string, 0)
    for _, v := range statDateList {
        match, _ := regexp.MatchString(`^2[0-9]{7}$`, v)
        if match {
            sil.statDateList = append(sil.statDateList, v)
        }
    }
}

func (sil *stepInfoList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sil.statDateList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "时间列表不能为空", nil))
    }
    if sil.dataType == 0 {
        if len(sil.userId) == 0 {
            panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
        }
        sil.ExtendData["object_id"] = sil.userId
    } else if sil.dataType == 1 {
        if sil.departmentId <= 0 {
            panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门ID不能为空", nil))
        }
        sil.ExtendData["object_id"] = sil.departmentId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "数据类型不能为空", nil))
    }
    sil.ExtendData["type"] = sil.dataType
    sil.ExtendData["stat_dates"] = strings.Join(sil.statDateList, ",")

    sil.ReqURI = dingtalk.UrlService + "/topapi/health/stepinfo/list?access_token=" + dingtalk.NewUtil().GetAccessToken(sil.corpId, sil.agentTag, sil.atType)

    reqBody := mpf.JSONMarshal(sil.ExtendData)
    client, req := sil.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewStepInfoList(corpId, agentTag, atType string) *stepInfoList {
    sil := &stepInfoList{dingtalk.NewCorp(), "", "", "", -1, "", 0, make([]string, 0)}
    sil.corpId = corpId
    sil.agentTag = agentTag
    sil.atType = atType
    sil.dataType = -1
    sil.ReqContentType = project.HTTPContentTypeJSON
    sil.ReqMethod = fasthttp.MethodPost
    return sil
}
