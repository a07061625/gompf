/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 23:32
 */
package oa

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取审批数据
type approvalDataGet struct {
    wx.BaseWxCorp
    corpId    string
    agentTag  string
    startTime int    // 开始时间
    endTime   int    // 结束时间
    nextSpNum string // 起始审批单号
}

func (adg *approvalDataGet) SetStartAndEndTime(startTime, endTime int) {
    if startTime <= 1000000000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "开始时间不合法", nil))
    } else if endTime <= 1000000000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "结束时间不合法", nil))
    } else if startTime > endTime {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "开始时间不能大于结束时间", nil))
    }
    adg.startTime = startTime
    adg.endTime = endTime
}

func (adg *approvalDataGet) SetNextSpNum(nextSpNum string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, nextSpNum)
    if match {
        adg.nextSpNum = nextSpNum
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "起始审批单号不合法", nil))
    }
}

func (adg *approvalDataGet) checkData() {
    if adg.startTime <= 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "开始时间不能为空", nil))
    }
}

func (adg *approvalDataGet) SendRequest() api.APIResult {
    adg.checkData()

    reqData := make(map[string]interface{})
    reqData["starttime"] = adg.startTime
    reqData["endtime"] = adg.endTime
    if len(adg.nextSpNum) > 0 {
        reqData["next_spnum"] = adg.nextSpNum
    }
    reqBody := mpf.JSONMarshal(reqData)

    adg.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/corp/getapprovaldata?access_token=" + wx.NewUtilWx().GetCorpAccessToken(adg.corpId, adg.agentTag)
    client, req := adg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := adg.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewApprovalDataGet(corpId, agentTag string) *approvalDataGet {
    adg := &approvalDataGet{wx.NewBaseWxCorp(), "", "", 0, 0, ""}
    adg.corpId = corpId
    adg.agentTag = agentTag
    adg.ReqContentType = project.HTTPContentTypeJSON
    adg.ReqMethod = fasthttp.MethodPost
    return adg
}
