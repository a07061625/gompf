/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 19:48
 */
package alipay

import (
    "crypto/tls"
    "net/url"
    "regexp"
    "strings"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type BaseAliPay struct {
    api.ApiOuter
    respTag       string                 // 响应标识
    BizContent    map[string]interface{} // 业务请求参数集合
    urlNotify     string                 // 主动通知地址
    urlReturnBase string                 // 跳转基础url地址
}

func (ap *BaseAliPay) SetMethod(method string) {
    ap.ReqData["method"] = method
    ap.respTag = strings.Replace(method, ".", "_", -1) + "_response"
}

func (ap *BaseAliPay) GetRespTag() string {
    return ap.respTag
}

func (ap *BaseAliPay) SetUrlNotify(notifyFlag bool) {
    if notifyFlag {
        ap.ReqData["notify_url"] = ap.urlNotify
    } else {
        delete(ap.ReqData, "notify_url")
    }
}

func (ap *BaseAliPay) SetUrlReturn(urlReturn string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, urlReturn)
    if match {
        ap.ReqData["return_url"] = ap.urlReturnBase + url.QueryEscape(urlReturn)
    } else {
        panic(mperr.NewAliPay(errorcode.AliPayParam, "同步通知地址不合法", nil))
    }
}

func (ap *BaseAliPay) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    ap.ReqData["biz_content"] = mpf.JsonMarshal(ap.BizContent)
    delete(ap.ReqData, "sign")
    sign := NewUtil().CreateSign(ap.ReqData, ap.ReqData["sign_type"])
    ap.ReqData["sign"] = sign
    reqBody := mpf.HttpCreateParams(ap.ReqData, "none", 1)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.SetBody([]byte(reqBody))
    req.Header.SetRequestURI(UrlGateWay)
    req.Header.SetContentType(project.HttpContentTypeForm)
    req.Header.SetMethod(fasthttp.MethodPost)
    mpf.HttpAddReqHeader(req, ap.ReqHeader)

    return client, req
}

func NewBase(appId string) BaseAliPay {
    now := time.Now()
    conf := NewConfig().GetAccount(appId)
    ap := BaseAliPay{api.NewApiOuter(), "", make(map[string]interface{}), "", ""}
    ap.urlNotify = conf.GetUrlNotify()
    ap.urlReturnBase = conf.GetUrlReturn()
    ap.ReqData["app_id"] = appId
    ap.ReqData["format"] = "json"
    ap.ReqData["charset"] = "utf-8"
    ap.ReqData["sign_type"] = "RSA2"
    ap.ReqData["timestamp"] = now.Format("2006-01-02 03:04:05")
    ap.ReqData["version"] = "1.0"
    ap.ReqHeader["Expect"] = ""
    return ap
}
