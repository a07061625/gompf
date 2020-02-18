/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 22:27
 */
package mppush

import (
    "crypto/tls"
    "net/url"
    "sort"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type basePush struct {
    api.APIOuter
    ExtendData map[string]interface{}
}

func newBasePush() basePush {
    return basePush{api.NewAPIOuter(), make(map[string]interface{})}
}

type BaseXinGe struct {
    basePush
    platform   string // 平台类型
    ServiceUri string // 服务URI
}

func (xg *BaseXinGe) setPlatform(platform string) {
    if (platform == XinGePlatformTypeAndroid) || (platform == XinGePlatformTypeIOS) {
        xg.platform = platform
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "平台类型不支持", nil))
    }
}

func (xg *BaseXinGe) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    if xg.platform == XinGePlatformTypeAndroid {
        conf := NewConfig().GetXinGeAndroid()
        xg.ReqHeader["Authorization"] = conf.GetAppAuth()
    } else {
        conf := NewConfig().GetXinGeIos()
        xg.ReqHeader["Authorization"] = conf.GetAppAuth()
    }
    xg.ReqURI = XinGeServiceDomain + xg.ServiceUri
    xg.ReqContentType = project.HTTPContentTypeJSON
    xg.ReqMethod = fasthttp.MethodPost

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    reqBody := mpf.JSONMarshal(xg.ExtendData)
    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(xg.ReqURI)
    req.Header.SetContentType(xg.ReqContentType)
    req.Header.SetMethod(xg.ReqMethod)
    req.SetBody([]byte(reqBody))
    mpf.HTTPAddReqHeader(req, xg.ReqHeader)

    return client, req
}

func NewBaseXinGe(platform string) BaseXinGe {
    xg := BaseXinGe{newBasePush(), "", ""}
    xg.setPlatform(platform)
    return xg
}

type BaseBaiDu struct {
    basePush
    ServiceUri string // 服务URI
}

func (bd *BaseBaiDu) createSign(secret string) {
    signStr := bd.ReqMethod + bd.ReqURI
    delete(bd.ReqData, "sign")
    pkParams := mpf.NewHTTPParamKey(bd.ReqData)
    sort.Sort(pkParams)
    paramNum := len(pkParams.Params)
    for i := 0; i < paramNum; i++ {
        signStr += pkParams.Params[i].Key + "=" + pkParams.Params[i].Val
    }
    signStr += secret
    bd.ReqData["sign"] = mpf.HashMd5(url.QueryEscape(signStr), "")
}

func (bd *BaseBaiDu) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    nowTime := time.Now().Unix()
    conf := NewConfig().GetBaiDu()
    bd.ReqData["apikey"] = conf.GetAppKey()
    bd.ReqData["timestamp"] = strconv.FormatInt(nowTime, 10)
    bd.ReqData["expires"] = strconv.FormatInt(nowTime+60, 10)
    if conf.GetDeviceType() == BaiDuDeviceTypeAll {
        bd.ReqData["device_type"] = BaiDuDeviceTypeAndroid
    } else {
        bd.ReqData["device_type"] = conf.GetDeviceType()
    }
    bd.ReqURI = BaiDuServiceDomain + BaiDuServiceUriPrefix + bd.ServiceUri
    bd.ReqContentType = project.HTTPContentTypeForm
    bd.ReqMethod = fasthttp.MethodPost
    bd.ReqHeader["Content-Type"] = project.HTTPContentTypeForm
    bd.ReqHeader["User-Agent"] = conf.GetUserAgent()
    bd.createSign(conf.GetAppSecret())

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    reqBody := mpf.HTTPCreateParams(bd.ReqData, "none", 1)
    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(bd.ReqURI)
    req.Header.SetContentType(bd.ReqContentType)
    req.Header.SetMethod(bd.ReqMethod)
    req.SetBody([]byte(reqBody))
    mpf.HTTPAddReqHeader(req, bd.ReqHeader)

    return client, req
}

func NewBaseBaiDu() BaseBaiDu {
    return BaseBaiDu{newBasePush(), ""}
}

type BaseJPush struct {
    basePush
    serviceDomain string // 服务域名
    ServiceUri    string // 服务URI
}

func (jp *BaseJPush) GetServiceUrl() string {
    return jp.serviceDomain + jp.ServiceUri
}

// 获取授权字符串
//   key: 应用标识
//   authType: 授权类型 app:应用 group:分组 dev:开发
func (jp *BaseJPush) getServiceAuth(key, authType string) string {
    switch authType {
    case "app":
        conf := NewConfig().GetJPushApp(key)
        return conf.GetAuth()
    case "group":
        conf := NewConfig().GetJPushGroup(key)
        return conf.GetAuth()
    default:
        conf := NewConfig().GetJPushDev()
        return conf.GetAuth()
    }
}

func (jp *BaseJPush) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(jp.ReqURI)
    req.Header.SetContentType(jp.ReqContentType)
    req.Header.SetMethod(jp.ReqMethod)
    mpf.HTTPAddReqHeader(req, jp.ReqHeader)

    return client, req
}

func NewBaseJPush(domain, key, authType string) BaseJPush {
    jp := BaseJPush{newBasePush(), "", ""}
    jp.serviceDomain = domain
    jp.ReqContentType = project.HTTPContentTypeJSON
    jp.ReqMethod = fasthttp.MethodGet
    jp.ReqHeader["Content-Type"] = "application/json"
    jp.ReqHeader["Authorization"] = jp.getServiceAuth(key, authType)
    return jp
}
