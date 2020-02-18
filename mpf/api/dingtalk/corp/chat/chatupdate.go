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

// 修改会话
type chatUpdate struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    chatId   string // 会话ID
}

func (cu *chatUpdate) SetChatId(chatId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, chatId)
    if match {
        cu.chatId = chatId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "会话ID不合法", nil))
    }
}

func (cu *chatUpdate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        cu.ExtendData["name"] = string(trueName[:10])
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "群名称不合法", nil))
    }
}

func (cu *chatUpdate) SetOwner(owner string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, owner)
    if match {
        cu.ExtendData["owner"] = owner
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "群主不合法", nil))
    }
}

func (cu *chatUpdate) SetAddUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "添加成员列表不能为空", nil))
    } else if len(userList) > 40 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "添加成员不能超过40个", nil))
    }

    addUserList := make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            addUserList = append(addUserList, v)
        }
    }

    if len(addUserList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "添加成员列表不能为空", nil))
    }
    cu.ExtendData["add_useridlist"] = addUserList
}

func (cu *chatUpdate) SetDelUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "删除成员列表不能为空", nil))
    } else if len(userList) > 40 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "删除成员不能超过40个", nil))
    }

    delUserList := make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            delUserList = append(delUserList, v)
        }
    }

    if len(delUserList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "删除成员列表不能为空", nil))
    }
    cu.ExtendData["del_useridlist"] = delUserList
}

func (cu *chatUpdate) SetIcon(icon string) {
    if len(icon) > 0 {
        cu.ExtendData["icon"] = icon
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "群头像不合法", nil))
    }
}

func (cu *chatUpdate) SetIsBan(isBan bool) {
    cu.ExtendData["isBan"] = isBan
}

// 搜索类型 0:默认,不可搜索 1:可搜索
func (cu *chatUpdate) SetSearchable(searchable int) {
    if (searchable == 0) || (searchable == 1) {
        cu.ExtendData["searchable"] = searchable
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "搜索类型不合法", nil))
    }
}

// 验证类型 0:默认,不验证 1:验证
func (cu *chatUpdate) SetValidationType(validationType int) {
    if (validationType == 0) || (validationType == 1) {
        cu.ExtendData["validationType"] = validationType
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "验证类型不合法", nil))
    }
}

// 通知所有人权限 0:默认,所有人 1:仅群主
func (cu *chatUpdate) SetMentionAllAuthority(mentionAllAuthority int) {
    if (mentionAllAuthority == 0) || (mentionAllAuthority == 1) {
        cu.ExtendData["mentionAllAuthority"] = mentionAllAuthority
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "通知所有人权限不合法", nil))
    }
}

// 管理类型 0:默认,所有人可管理 1:仅群主可管理
func (cu *chatUpdate) SetManagementType(managementType int) {
    if (managementType == 0) || (managementType == 1) {
        cu.ExtendData["managementType"] = managementType
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "管理类型不合法", nil))
    }
}

func (cu *chatUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cu.chatId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "会话ID不能为空", nil))
    }
    cu.ExtendData["chatid"] = cu.chatId

    cu.ReqURI = dingtalk.UrlService + "/chat/update?acuess_token=" + dingtalk.NewUtil().GetCorpAccessToken(cu.corpId, cu.agentTag)

    reqBody := mpf.JSONMarshal(cu.ExtendData)
    client, req := cu.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewChatUpdate(corpId, agentTag string) *chatUpdate {
    cu := &chatUpdate{dingtalk.NewCorp(), "", "", ""}
    cu.corpId = corpId
    cu.agentTag = agentTag
    cu.ReqContentType = project.HTTPContentTypeJSON
    cu.ReqMethod = fasthttp.MethodPost
    return cu
}
