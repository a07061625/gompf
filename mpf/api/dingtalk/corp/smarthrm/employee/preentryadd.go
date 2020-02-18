/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 0:12
 */
package employee

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

// 添加企业待入职员工
type preEntryAdd struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    name     string // 姓名
    mobile   string // 手机号
}

func (pea *preEntryAdd) SetName(name string) {
    if len(name) > 0 {
        pea.name = name
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "姓名不合法", nil))
    }
}

func (pea *preEntryAdd) SetMobile(mobile string) {
    match, _ := regexp.MatchString(project.RegexPhone, mobile)
    if match {
        pea.mobile = mobile
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "手机号不合法", nil))
    }
}

func (pea *preEntryAdd) SetPreEntryTime(preEntryTime int) {
    if preEntryTime > 946656000 {
        et := time.Unix(int64(preEntryTime), 0)
        pea.ExtendData["pre_entry_time"] = et.Format("2006-01-02 03:04:05")
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "预期入职时间不合法", nil))
    }
}

func (pea *preEntryAdd) SetOpUserId(opUserId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, opUserId)
    if match {
        pea.ExtendData["op_userid"] = opUserId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "操作人用户ID不合法", nil))
    }
}

func (pea *preEntryAdd) SetExtendInfo(extendInfo map[string]interface{}) {
    if len(extendInfo) > 0 {
        pea.ExtendData["extend_info"] = mpf.JsonMarshal(extendInfo)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "扩展信息不合法", nil))
    }
}

func (pea *preEntryAdd) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pea.name) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "姓名不能为空", nil))
    }
    if len(pea.mobile) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "手机号不能为空", nil))
    }
    pea.ExtendData["name"] = pea.name
    pea.ExtendData["mobile"] = pea.mobile

    pea.ReqUrl = dingtalk.UrlService + "/topapi/smartwork/hrm/employee/addpreentry?access_token=" + dingtalk.NewUtil().GetAccessToken(pea.corpId, pea.agentTag, pea.atType)

    reqData := make(map[string]interface{})
    reqData["param"] = pea.ExtendData
    reqBody := mpf.JsonMarshal(reqData)
    client, req := pea.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewPreEntryAdd(corpId, agentTag, atType string) *preEntryAdd {
    pea := &preEntryAdd{dingtalk.NewCorp(), "", "", "", "", ""}
    pea.corpId = corpId
    pea.agentTag = agentTag
    pea.atType = atType
    pea.ReqContentType = project.HTTPContentTypeJSON
    pea.ReqMethod = fasthttp.MethodPost
    return pea
}
