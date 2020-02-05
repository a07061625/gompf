package department

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取子部门ID列表
type idList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    id       int // 部门id
}

func (il *idList) SetId(id int) {
    if id > 0 {
        il.id = id
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不合法", nil))
    }
}

func (il *idList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if il.id <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不能为空", nil))
    }
    il.ReqData["id"] = strconv.Itoa(il.id)
    il.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(il.corpId, il.agentTag, il.atType)
    il.ReqUrl = dingtalk.UrlService + "/department/list_ids?" + mpf.HttpCreateParams(il.ReqData, "none", 1)

    return il.GetRequest()
}

func NewIdList(corpId, agentTag, atType string) *idList {
    il := &idList{dingtalk.NewCorp(), "", "", "", 0}
    il.corpId = corpId
    il.agentTag = agentTag
    il.atType = atType
    return il
}
