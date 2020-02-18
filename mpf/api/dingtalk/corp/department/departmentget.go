package department

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取部门详情
type departmentGet struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    id       int // 部门id
}

func (dg *departmentGet) SetId(id int) {
    if id > 0 {
        dg.id = id
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不合法", nil))
    }
}

func (dg *departmentGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if dg.id <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不能为空", nil))
    }
    dg.ReqData["id"] = strconv.Itoa(dg.id)
    dg.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(dg.corpId, dg.agentTag, dg.atType)
    dg.ReqUrl = dingtalk.UrlService + "/department/get?" + mpf.HTTPCreateParams(dg.ReqData, "none", 1)

    return dg.GetRequest()
}

func NewDepartmentGet(corpId, agentTag, atType string) *departmentGet {
    dg := &departmentGet{dingtalk.NewCorp(), "", "", "", 0}
    dg.corpId = corpId
    dg.agentTag = agentTag
    dg.atType = atType
    dg.ReqData["lang"] = "zh_CN"
    return dg
}
