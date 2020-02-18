/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/8 0008
 * Time: 0:14
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

// 查询自建应用审批单当前状态
type openApprovalDataGet struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    thirdNo  string // 审批单号
}

func (adg *openApprovalDataGet) SetThirdNo(thirdNo string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, thirdNo)
    if match {
        adg.thirdNo = thirdNo
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "审批单号不合法", nil))
    }
}

func (adg *openApprovalDataGet) checkData() {
    if len(adg.thirdNo) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "审批单号不能为空", nil))
    }
    adg.ReqData["thirdNo"] = adg.thirdNo
}

func (adg *openApprovalDataGet) SendRequest(getType string) api.ApiResult {
    adg.checkData()

    reqBody := mpf.JSONMarshal(adg.ReqData)
    adg.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/corp/getopenapprovaldata?access_token=" + wx.NewUtilWx().GetCorpCache(adg.corpId, adg.agentTag, getType)
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

func NewOpenApprovalDataGet(corpId, agentTag string) *openApprovalDataGet {
    adg := &openApprovalDataGet{wx.NewBaseWxCorp(), "", "", ""}
    adg.corpId = corpId
    adg.agentTag = agentTag
    adg.ReqContentType = project.HTTPContentTypeJSON
    adg.ReqMethod = fasthttp.MethodPost
    return adg
}
