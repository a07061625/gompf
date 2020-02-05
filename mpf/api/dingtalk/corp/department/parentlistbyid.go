package department

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询部门的所有上级父部门路径
type parentListById struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    id       int // 部门id
}

func (pli *parentListById) SetId(id int) {
    if id > 0 {
        pli.id = id
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不合法", nil))
    }
}

func (pli *parentListById) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if pli.id <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不能为空", nil))
    }
    pli.ReqData["id"] = strconv.Itoa(pli.id)
    pli.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(pli.corpId, pli.agentTag, pli.atType)
    pli.ReqUrl = dingtalk.UrlService + "/department/list_parent_depts_by_dept?" + mpf.HttpCreateParams(pli.ReqData, "none", 1)

    return pli.GetRequest()
}

func NewParentListById(corpId, agentTag, atType string) *parentListById {
    pli := &parentListById{dingtalk.NewCorp(), "", "", "", 0}
    pli.corpId = corpId
    pli.agentTag = agentTag
    pli.atType = atType
    return pli
}
