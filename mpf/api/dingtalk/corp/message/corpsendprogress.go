package message

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询工作通知消息的发送进度
type corpSendProgress struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    taskId   int // 任务ID
}

func (csp *corpSendProgress) SetTaskId(taskId int) {
    if taskId > 0 {
        csp.taskId = taskId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "任务ID不合法", nil))
    }
}

func (csp *corpSendProgress) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if csp.taskId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "任务ID不能为空", nil))
    }
    csp.ExtendData["task_id"] = csp.taskId

    csp.ReqUrl = dingtalk.UrlService + "/topapi/message/corpconversation/getsendprogress?access_token=" + dingtalk.NewUtil().GetAccessToken(csp.corpId, csp.agentTag, csp.atType)

    reqBody := mpf.JsonMarshal(csp.ExtendData)
    client, req := csp.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewCorpSendProgress(corpId, agentTag, atType string) *corpSendProgress {
    agentInfo := dingtalk.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    csp := &corpSendProgress{dingtalk.NewCorp(), "", "", "", 0}
    csp.corpId = corpId
    csp.agentTag = agentTag
    csp.atType = atType
    csp.ExtendData["agent_id"] = agentInfo["id"]
    csp.ReqContentType = project.HttpContentTypeJson
    csp.ReqMethod = fasthttp.MethodPost
    return csp
}
