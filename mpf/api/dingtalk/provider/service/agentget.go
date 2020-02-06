package service

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

// 获取授权企业的应用信息
type agentGet struct {
    dingtalk.BaseProvider
    corpId   string
    agentTag string
}

func (ag *agentGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    providerConf := dingtalk.NewConfig().GetProvider()
    suiteTicket := dingtalk.NewUtil().GetProviderSuiteTicket()
    nowTime := strconv.FormatInt(time.Now().Unix(), 10)
    corpConf := dingtalk.NewConfig().GetCorp(ag.corpId)
    agentInfo := corpConf.GetAgentInfo(ag.agentTag)
    ag.ExtendData["agentid"] = agentInfo["id"]
    ag.ExtendData["suite_key"] = providerConf.GetSuiteKey()

    ag.ReqData["timestamp"] = nowTime
    ag.ReqData["accessKey"] = providerConf.GetSuiteKey()
    ag.ReqData["suiteTicket"] = suiteTicket
    ag.ReqData["signature"] = dingtalk.NewUtil().CreateApiSign(nowTime+"\n"+suiteTicket, providerConf.GetSuiteSecret())
    ag.ReqUrl = dingtalk.UrlService + "/service/get_agent?" + mpf.HttpCreateParams(ag.ReqData, "none", 1)

    reqBody := mpf.JsonMarshal(ag.ExtendData)
    client, req := ag.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewAgentGet(corpId, agentTag string) *agentGet {
    ag := &agentGet{dingtalk.NewProvider(), "", ""}
    ag.corpId = corpId
    ag.agentTag = agentTag
    ag.ExtendData["auth_corpid"] = corpId
    ag.ReqContentType = project.HttpContentTypeJson
    ag.ReqMethod = fasthttp.MethodPost
    return ag
}
