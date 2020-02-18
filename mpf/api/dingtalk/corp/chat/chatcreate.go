package chat

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建会话
type chatCreate struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    name     string   // 群名称
    owner    string   // 群主
    userList []string // 成员列表
}

func (cc *chatCreate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        cc.name = string(trueName[:10])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "群名称不合法", nil))
    }
}

func (cc *chatCreate) SetOwner(owner string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, owner)
    if match {
        cc.owner = owner
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "群主不合法", nil))
    }
}

func (cc *chatCreate) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "成员列表不能为空", nil))
    } else if len(userList) > 40 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "成员不能超过40个", nil))
    }

    cc.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            cc.userList = append(cc.userList, v)
        }
    }

    if len(cc.userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "成员列表不能为空", nil))
    }
}

// 新成员查看聊天历史消息标识 0:否 1:是
func (cc *chatCreate) SetShowHistoryType(showHistoryType int) {
    if (showHistoryType == 0) || (showHistoryType == 1) {
        cc.ExtendData["showHistoryType"] = showHistoryType
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "新成员查看聊天历史消息标识不合法", nil))
    }
}

// 搜索类型 0:默认,不可搜索 1:可搜索
func (cc *chatCreate) SetSearchable(searchable int) {
    if (searchable == 0) || (searchable == 1) {
        cc.ExtendData["searchable"] = searchable
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "搜索类型不合法", nil))
    }
}

// 验证类型 0:默认,不验证 1:验证
func (cc *chatCreate) SetValidationType(validationType int) {
    if (validationType == 0) || (validationType == 1) {
        cc.ExtendData["validationType"] = validationType
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "验证类型不合法", nil))
    }
}

// 通知所有人权限 0:默认,所有人 1:仅群主
func (cc *chatCreate) SetMentionAllAuthority(mentionAllAuthority int) {
    if (mentionAllAuthority == 0) || (mentionAllAuthority == 1) {
        cc.ExtendData["mentionAllAuthority"] = mentionAllAuthority
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "通知所有人权限不合法", nil))
    }
}

// 管理类型 0:默认,所有人可管理 1:仅群主可管理
func (cc *chatCreate) SetManagementType(managementType int) {
    if (managementType == 0) || (managementType == 1) {
        cc.ExtendData["managementType"] = managementType
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "管理类型不合法", nil))
    }
}

func (cc *chatCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cc.name) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "群名称不能为空", nil))
    }
    if len(cc.owner) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "群主不能为空", nil))
    }
    if len(cc.userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "成员列表不能为空", nil))
    }
    cc.ExtendData["name"] = cc.name
    cc.ExtendData["owner"] = cc.owner
    cc.ExtendData["useridlist"] = cc.userList

    cc.ReqUrl = dingtalk.UrlService + "/chat/create?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(cc.corpId, cc.agentTag)

    reqBody := mpf.JsonMarshal(cc.ExtendData)
    client, req := cc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewChatCreate(corpId, agentTag string) *chatCreate {
    cc := &chatCreate{dingtalk.NewCorp(), "", "", "", "", make([]string, 0)}
    cc.corpId = corpId
    cc.agentTag = agentTag
    cc.ExtendData["showHistoryType"] = 0
    cc.ExtendData["searchable"] = 1
    cc.ExtendData["validationType"] = 0
    cc.ExtendData["mentionAllAuthority"] = 0
    cc.ExtendData["managementType"] = 1
    cc.ReqContentType = project.HTTPContentTypeJSON
    cc.ReqMethod = fasthttp.MethodPost
    return cc
}
