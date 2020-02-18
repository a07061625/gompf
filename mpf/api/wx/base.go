/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 19:48
 */
package wx

import (
    "crypto/tls"
    "encoding/xml"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type DataOpenAuthorize struct {
    AuthorizeStatus int
    AuthCode        string
    RefreshToken    string
}

func NewDataOpenAuthorize() *DataOpenAuthorize {
    return &DataOpenAuthorize{0, "", ""}
}

type DataProviderAuthorize struct {
    AuthorizeStatus int
    PermanentCode   string
}

func NewDataProviderAuthorize() *DataProviderAuthorize {
    return &DataProviderAuthorize{0, ""}
}

type WxNotifyOpen struct {
    XMLName xml.Name `xml:"xml"`
    AppId   string
    Encrypt string
}

type WxCDATAText struct {
    Text string `xml:",innerxml"`
}

type WxResponse struct {
    XMLName      xml.Name `xml:"xml"`
    Encrypt      WxCDATAText
    MsgSignature WxCDATAText
    TimeStamp    string
    Nonce        WxCDATAText
}

type baseWx struct {
    api.ApiInner
}

func (bw *baseWx) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(bw.ReqUrl)
    req.Header.SetContentType(bw.ReqContentType)
    req.Header.SetMethod(bw.ReqMethod)
    mpf.HttpAddReqHeader(req, bw.ReqHeader)

    return client, req
}

func newBaseWx() baseWx {
    bw := baseWx{api.NewApiInner()}
    bw.ReqContentType = project.HTTPContentTypeForm
    bw.ReqMethod = fasthttp.MethodGet
    return bw
}

type BaseWxAccount struct {
    baseWx
}

func (wa *BaseWxAccount) SetPayAccount(accountConfig *configAccount, merchantType string) {
    switch merchantType {
    case AccountMerchantTypeSelf:
        wa.ReqData["appid"] = accountConfig.GetAppId()
        wa.ReqData["mch_id"] = accountConfig.GetPayMchId()
    case AccountMerchantTypeSub:
        merchantAppId := accountConfig.GetMerchantAppId()
        if len(merchantAppId) == 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "服务商微信号不能为空", nil))
        }
        merchantConfig := NewConfig().GetAccount(merchantAppId)
        wa.ReqData["appid"] = merchantConfig.GetAppId()
        wa.ReqData["mch_id"] = merchantConfig.GetPayMchId()
        wa.ReqData["sub_appid"] = accountConfig.GetAppId()
        wa.ReqData["sub_mch_id"] = accountConfig.GetPayMchId()
    default:
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户类型不合法", nil))
    }
}

func NewBaseWxAccount() BaseWxAccount {
    return BaseWxAccount{newBaseWx()}
}

type BaseWxMini struct {
    baseWx
}

func NewBaseWxMini() BaseWxMini {
    return BaseWxMini{newBaseWx()}
}

type BaseWxCorp struct {
    baseWx
}

func NewBaseWxCorp() BaseWxCorp {
    return BaseWxCorp{newBaseWx()}
}

type BaseWxOpen struct {
    baseWx
}

func NewBaseWxOpen() BaseWxOpen {
    return BaseWxOpen{newBaseWx()}
}

type BaseWxProvider struct {
    baseWx
}

func NewBaseWxProvider() BaseWxProvider {
    return BaseWxProvider{newBaseWx()}
}

type BaseWxSingle struct {
    baseWx
}

func NewBaseWxSingle() BaseWxSingle {
    return BaseWxSingle{newBaseWx()}
}
