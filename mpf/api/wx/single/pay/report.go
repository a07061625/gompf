/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 19:51
 */
package pay

import (
    "encoding/xml"
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type report struct {
    wx.BaseWxAccount
    appId        string
    deviceInfo   string                 // 设备号
    interfaceUrl string                 // 接口URL
    trades       map[string]interface{} // 上报数据包
}

func (r *report) SetDeviceInfo(deviceInfo string) {
    if len(deviceInfo) > 0 {
        r.ReqData["device_info"] = deviceInfo
    }
}

func (r *report) SetInterfaceUrl(interfaceUrl string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, interfaceUrl)
    if match {
        r.interfaceUrl = interfaceUrl
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "接口URL不合法", nil))
    }
}

func (r *report) SetTrades(trades map[string]interface{}) {
    if len(trades) > 0 {
        r.trades = trades
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "上报数据包不合法", nil))
    }
}

func (r *report) checkData() {
    if len(r.interfaceUrl) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "接口URL不能为空", nil))
    }
    if len(r.trades) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "上报数据包不能为空", nil))
    }
    r.ReqData["interface_url"] = r.interfaceUrl
    r.ReqData["trades"] = mpf.JsonMarshal(r.trades)
}

func (r *report) SendRequest() api.ApiResult {
    r.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(r.ReqData, r.appId, "md5")
    r.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(r.ReqData))
    r.ReqUrl = "https://api.mch.weixin.qq.com/payitil/report"
    client, req := r.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := r.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData := make(map[string]string)
    xml.Unmarshal(resp.Body, (*mpf.XmlMap)(&respData))
    if respData["return_code"] == "FAIL" {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["return_msg"]
    } else if respData["result_code"] == "FAIL" {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = "上报失败"
    } else {
        result.Data = respData
    }
    return result
}

func NewReport(appId, merchantType string) *report {
    conf := wx.NewConfig().GetAccount(appId)
    r := &report{wx.NewBaseWxAccount(), "", "", "", make(map[string]interface{})}
    r.appId = appId
    r.SetPayAccount(conf, merchantType)
    r.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    r.ReqData["user_ip"] = conf.GetClientIp()
    r.ReqContentType = project.HttpContentTypeXml
    r.ReqMethod = fasthttp.MethodPost
    return r
}
