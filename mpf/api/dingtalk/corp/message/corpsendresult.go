package message

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询工作通知消息的发送结果
type corpSendResult struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    taskId   int // 任务ID
}

func (csr *corpSendResult) SetTaskId(taskId int) {
    if taskId > 0 {
        csr.taskId = taskId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "任务ID不合法", nil))
    }
}

func (csr *corpSendResult) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if csr.taskId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "任务ID不能为空", nil))
    }
    csr.ExtendData["task_id"] = csr.taskId

    csr.ReqURI = dingtalk.UrlService + "/topapi/message/corpconversation/getsendresult?access_token=" + dingtalk.NewUtil().GetAccessToken(csr.corpId, csr.agentTag, csr.atType)

    reqBody := mpf.JSONMarshal(csr.ExtendData)
    client, req := csr.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewCorpSendResult(corpId, agentTag, atType string) *corpSendResult {
    agentInfo := dingtalk.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    csr := &corpSendResult{dingtalk.NewCorp(), "", "", "", 0}
    csr.corpId = corpId
    csr.agentTag = agentTag
    csr.atType = atType
    csr.ExtendData["agent_id"] = agentInfo["id"]
    csr.ReqContentType = project.HTTPContentTypeJSON
    csr.ReqMethod = fasthttp.MethodPost
    return csr
}
