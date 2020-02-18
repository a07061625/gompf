/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 23:03
 */
package authorize

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

type code2Openid struct {
    wx.BaseWxAccount
    appId    string
    authCode string // 授权码
}

func (co *code2Openid) SetAuthCode(authCode string) {
    match, _ := regexp.MatchString(`^1[0-9]{17}$`, authCode)
    if match {
        co.authCode = authCode
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "授权码不合法", nil))
    }
}

func (co *code2Openid) checkData() {
    if len(co.authCode) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "授权码不能为空", nil))
    }
    co.ReqData["auth_code"] = co.authCode
}

func (co *code2Openid) SendRequest() api.APIResult {
    co.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(co.ReqData, co.appId, "md5")
    co.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(co.ReqData))
    co.ReqURI = "https://api.mch.weixin.qq.com/tools/authcodetoopenid"
    client, req := co.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := co.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData := make(map[string]string)
    xml.Unmarshal(resp.Body, (*mpf.XMLMap)(&respData))
    if respData["return_code"] == "FAIL" {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["return_msg"]
    } else if respData["result_code"] == "FAIL" {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["err_code"]
    } else {
        result.Data = respData
    }
    return result
}

func NewCode2Openid(appId string) *code2Openid {
    conf := wx.NewConfig().GetAccount(appId)
    co := &code2Openid{wx.NewBaseWxAccount(), "", ""}
    co.appId = appId
    co.ReqData["appid"] = conf.GetAppId()
    co.ReqData["mch_id"] = conf.GetPayMchId()
    co.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    co.ReqContentType = project.HTTPContentTypeXML
    co.ReqMethod = fasthttp.MethodPost
    return co
}
