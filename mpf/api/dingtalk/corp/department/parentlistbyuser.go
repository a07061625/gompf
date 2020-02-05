package department

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询指定用户的所有上级父部门路径
type parentListByUser struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户id
}

func (plu *parentListByUser) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        plu.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "群主不合法", nil))
    }
}

func (plu *parentListByUser) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(plu.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不能为空", nil))
    }
    plu.ReqData["userId"] = plu.userId
    plu.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(plu.corpId, plu.agentTag, plu.atType)
    plu.ReqUrl = dingtalk.UrlService + "/department/list_parent_depts?" + mpf.HttpCreateParams(plu.ReqData, "none", 1)

    return plu.GetRequest()
}

func NewParentListByUser(corpId, agentTag, atType string) *parentListByUser {
    plu := &parentListByUser{dingtalk.NewCorp(), "", "", "", ""}
    plu.corpId = corpId
    plu.agentTag = agentTag
    plu.atType = atType
    return plu
}
