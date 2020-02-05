package message

import (
    "regexp"
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 发送工作通知消息
type corpSendAsync struct {
    dingtalk.BaseCorp
    corpId     string
    agentTag   string
    atType     string
    userList   []string               // 用户列表
    departList []string               // 部门列表
    msgContent map[string]interface{} // 消息内容
}

func (csa *corpSendAsync) SetUserList(userList []string) {
    if len(userList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表不能为空", nil))
    } else if len(userList) > 20 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户不能超过20个", nil))
    }

    csa.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            csa.userList = append(csa.userList, v)
        }
    }
}

func (csa *corpSendAsync) SetMsgContent(msgType string, msgData map[string]interface{}) {
    _, ok := dingtalk.MessageTypes[msgType]
    if !ok {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "消息类型不支持", nil))
    }
    if len(msgData) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "消息数据不能为空", nil))
    }
    csa.msgContent["msgtype"] = msgType
    csa.msgContent[msgType] = msgData
}

func (csa *corpSendAsync) SetDepartList(departList []int) {
    if len(departList) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门列表不能为空", nil))
    } else if len(departList) > 20 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门不能超过20个", nil))
    }

    csa.departList = make([]string, 0)
    for _, v := range departList {
        if v > 0 {
            csa.departList = append(csa.departList, strconv.Itoa(v))
        }
    }
}

func (csa *corpSendAsync) SetToAllUser(toAllUser bool) {
    csa.ExtendData["to_all_user"] = toAllUser
}

func (csa *corpSendAsync) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(csa.msgContent) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "消息内容不能为空", nil))
    }
    if len(csa.userList) > 0 {
        csa.ExtendData["userid_list"] = strings.Join(csa.userList, ",")
    } else if len(csa.departList) > 0 {
        csa.ExtendData["dept_id_list"] = strings.Join(csa.departList, ",")
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户列表和部门列表不能同时为空", nil))
    }
    csa.ExtendData["msg"] = csa.msgContent

    csa.ReqUrl = dingtalk.UrlService + "/topapi/message/corpconversation/asyncsend_v2?access_token=" + dingtalk.NewUtil().GetAccessToken(csa.corpId, csa.agentTag, csa.atType)

    reqBody := mpf.JsonMarshal(csa.ExtendData)
    client, req := csa.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewCorpSendAsync(corpId, agentTag, atType string) *corpSendAsync {
    agentInfo := dingtalk.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    csa := &corpSendAsync{dingtalk.NewCorp(), "", "", "", make([]string, 0), make([]string, 0), make(map[string]interface{})}
    csa.corpId = corpId
    csa.agentTag = agentTag
    csa.atType = atType
    csa.ExtendData["agent_id"] = agentInfo["id"]
    csa.ExtendData["to_all_user"] = false
    csa.ReqContentType = project.HttpContentTypeJson
    csa.ReqMethod = fasthttp.MethodPost
    return csa
}
