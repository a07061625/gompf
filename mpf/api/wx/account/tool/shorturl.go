/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 23:25
 */
package tool

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

type shortUrl struct {
    wx.BaseWxAccount
    appId   string
    longUrl string // URL链接
}

func (su *shortUrl) SetLongUrl(longUrl string) {
    match, _ := regexp.MatchString(`^weixin`, longUrl)
    if match {
        su.longUrl = longUrl
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "长链接不合法", nil))
    }
}

func (su *shortUrl) checkData() {
    if len(su.longUrl) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "长链接不能为空", nil))
    }
    su.ReqData["long_url"] = su.longUrl
}

func (su *shortUrl) SendRequest() api.APIResult {
    su.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(su.ReqData, su.appId, "md5")
    su.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(su.ReqData))
    su.ReqURI = "https://api.mch.weixin.qq.com/tools/shorturl"
    client, req := su.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := su.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewShortUrl(appId string) *shortUrl {
    conf := wx.NewConfig().GetAccount(appId)
    su := &shortUrl{wx.NewBaseWxAccount(), "", ""}
    su.appId = appId
    su.ReqData["appid"] = conf.GetAppId()
    su.ReqData["mch_id"] = conf.GetPayMchId()
    su.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    su.ReqData["sign_type"] = "MD5"
    su.ReqContentType = project.HTTPContentTypeXML
    su.ReqMethod = fasthttp.MethodPost
    return su
}
