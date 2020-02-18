/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 12:50
 */
package agent

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/valyala/fasthttp"
)

// 获取指定的应用详情
type agentGet struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    agentId  string // 应用ID
}

func (ag *agentGet) checkData() (*fasthttp.Client, *fasthttp.Request) {
    agentInfo := wx.NewConfig().GetCorp(ag.corpId).GetAgentInfo(ag.agentTag)
    ag.ReqData["agentid"] = agentInfo["id"]
    ag.ReqData["access_token"] = wx.NewUtilWx().GetCorpAccessToken(ag.corpId, ag.agentTag)
    ag.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/agent/get?" + mpf.HTTPCreateParams(ag.ReqData, "none", 1)
    return ag.GetRequest()
}

func (ag *agentGet) SendRequest() api.APIResult {
    client, req := ag.checkData()
    resp, result := ag.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewAgentGet(corpId, agentTag string) *agentGet {
    ag := &agentGet{wx.NewBaseWxCorp(), "", "", ""}
    ag.corpId = corpId
    ag.agentTag = agentTag
    return ag
}
