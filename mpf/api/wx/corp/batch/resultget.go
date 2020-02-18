/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 23:17
 */
package batch

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取异步任务结果
type resultGet struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    jobId    string // 任务id
}

func (rg *resultGet) SetJobId(jobId string) {
    if len(jobId) > 0 {
        rg.jobId = jobId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "任务id不合法", nil))
    }
}

func (rg *resultGet) checkData() {
    if len(rg.jobId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "任务id不能为空", nil))
    }
    rg.ReqData["jobid"] = rg.jobId
}

func (rg *resultGet) SendRequest() api.APIResult {
    rg.checkData()

    rg.ReqData["access_token"] = wx.NewUtilWx().GetCorpAccessToken(rg.corpId, rg.agentTag)
    rg.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/batch/getresult?" + mpf.HTTPCreateParams(rg.ReqData, "none", 1)
    client, req := rg.GetRequest()

    resp, result := rg.SendInner(client, req, errorcode.WxCorpRequestGet)
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

func NewResultGet(corpId, agentTag string) *resultGet {
    rg := &resultGet{wx.NewBaseWxCorp(), "", "", ""}
    rg.corpId = corpId
    rg.agentTag = agentTag
    rg.ReqContentType = project.HTTPContentTypeJSON
    rg.ReqMethod = fasthttp.MethodPost
    return rg
}
