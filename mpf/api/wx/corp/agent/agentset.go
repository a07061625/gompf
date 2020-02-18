/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 13:06
 */
package agent

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 设置应用
type agentSet struct {
    wx.BaseWxCorp
    corpId             string
    agentTag           string
    reportLocationFlag int    // 地理位置上报标识 0:不上报 1:上报
    logoMediaId        string // 应用头像
    name               string // 应用名称
    description        string // 应用详情
    redirectDomain     string // 应用可信域名
    isReportEnter      int    // 用户进入上报标识 0:不接收 1:接收
    homeUrl            string // 应用主页url
}

func (as *agentSet) SetReportLocationFlag(reportLocationFlag int) {
    if (reportLocationFlag == 0) || (reportLocationFlag == 1) {
        as.ReqData["report_location_flag"] = strconv.Itoa(reportLocationFlag)
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "地理位置上报标识不合法", nil))
    }
}

func (as *agentSet) SetLogoMediaId(logoMediaId string) {
    if len(logoMediaId) > 0 {
        as.ReqData["logo_mediaid"] = logoMediaId
    } else {
        delete(as.ReqData, "logo_mediaid")
    }
}

func (as *agentSet) SetName(name string) {
    if len(name) > 0 {
        nameRune := []rune(name)
        as.ReqData["name"] = string(nameRune[:16])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "应用名称不合法", nil))
    }
}

func (as *agentSet) SetDescription(description string) {
    length := len(description)
    if length < 4 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "应用详情不能少于4个字节", nil))
    } else if length > 120 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "应用详情不能大于120个字节", nil))
    }
    as.ReqData["description"] = description
}

func (as *agentSet) SetRedirectDomain(redirectDomain string) {
    if len(redirectDomain) > 0 {
        as.ReqData["redirect_domain"] = redirectDomain
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "应用可信域名不合法", nil))
    }
}

func (as *agentSet) SetReportEnterFlag(reportEnterFlag int) {
    if (reportEnterFlag == 0) || (reportEnterFlag == 1) {
        as.ReqData["isreportenter"] = strconv.Itoa(reportEnterFlag)
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户进入上报标识不合法", nil))
    }
}

func (as *agentSet) SetHomeUrl(homeUrl string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, homeUrl)
    if match {
        as.ReqData["home_url"] = homeUrl
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "应用主页url不合法", nil))
    }
}

func (as *agentSet) checkData() (*fasthttp.Client, *fasthttp.Request) {
    as.ReqData["access_token"] = wx.NewUtilWx().GetCorpAccessToken(as.corpId, as.agentTag)
    as.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/agent/set?" + mpf.HttpCreateParams(as.ReqData, "none", 1)
    return as.GetRequest()
}

func (as *agentSet) SendRequest() api.ApiResult {
    client, req := as.checkData()
    resp, result := as.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewAgentSet(corpId, agentTag string) *agentSet {
    agentInfo := wx.NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    as := &agentSet{wx.NewBaseWxCorp(), "", "", 0, "", "", "", "", 0, ""}
    as.corpId = corpId
    as.agentTag = agentTag
    as.ReqData["agentid"] = agentInfo["id"]
    as.ReqData["report_location_flag"] = "0"
    as.ReqData["isreportenter"] = "0"
    return as
}
