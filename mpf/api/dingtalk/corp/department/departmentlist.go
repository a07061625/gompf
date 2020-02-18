package department

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取部门列表
type departmentList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    id       int // 部门id
}

func (dl *departmentList) SetId(id int) {
    if id > 0 {
        dl.id = id
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不合法", nil))
    }
}

func (dl *departmentList) SetFtchChild(fetchChild int) {
    if (fetchChild == 0) || (fetchChild == 1) {
        dl.ReqData["fetch_child"] = strconv.Itoa(fetchChild)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "递归部门标识不合法", nil))
    }
}

func (dl *departmentList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if dl.id <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不能为空", nil))
    }
    dl.ReqData["id"] = strconv.Itoa(dl.id)
    dl.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(dl.corpId, dl.agentTag, dl.atType)
    dl.ReqURI = dingtalk.UrlService + "/department/list?" + mpf.HTTPCreateParams(dl.ReqData, "none", 1)

    return dl.GetRequest()
}

func NewDepartmentList(corpId, agentTag, atType string) *departmentList {
    dl := &departmentList{dingtalk.NewCorp(), "", "", "", 0}
    dl.corpId = corpId
    dl.agentTag = agentTag
    dl.atType = atType
    dl.ReqData["lang"] = "zh_CN"
    dl.ReqData["fetch_child"] = "0"
    return dl
}
