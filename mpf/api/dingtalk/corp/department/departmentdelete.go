package department

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除部门
type departmentDelete struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    id       int // 部门id
}

func (dd *departmentDelete) SetId(id int) {
    if id > 0 {
        dd.id = id
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不合法", nil))
    }
}

func (dd *departmentDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if dd.id <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不能为空", nil))
    }
    dd.ReqData["id"] = strconv.Itoa(dd.id)
    dd.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(dd.corpId, dd.agentTag, dd.atType)
    dd.ReqUrl = dingtalk.UrlService + "/department/delete?" + mpf.HttpCreateParams(dd.ReqData, "none", 1)

    return dd.GetRequest()
}

func NewDepartmentDelete(corpId, agentTag, atType string) *departmentDelete {
    dd := &departmentDelete{dingtalk.NewCorp(), "", "", "", 0}
    dd.corpId = corpId
    dd.agentTag = agentTag
    dd.atType = atType
    return dd
}
