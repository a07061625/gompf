/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/25 0025
 * Time: 22:52
 */
package mpmap

import (
    "crypto/tls"
    "net/url"
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type baseMap struct {
    api.ApiOuter
    respTag string // 响应标识
}

func (m *baseMap) GetRespTag() string {
    return m.respTag
}

func (m *baseMap) SetRespTag(respTag string) {
    m.respTag = respTag
}

func newBase() baseMap {
    return baseMap{api.NewApiOuter(), ""}
}

type BaseBaiDu struct {
    baseMap
    serviceUri    string // 服务uri
    serviceDomain string // 服务域名
    ak            string // 应用密钥
    serverIp      string // 服务端IP
    output        string // 输出格式
    checkType     string // 校验类型
    sk            string // 用户签名
    reqRefer      string // 请求引用地址
}

func (m *BaseBaiDu) SetServiceUri(uri string) {
    m.serviceUri = uri
    m.ReqUrl = m.serviceDomain + uri
}

func (m *BaseBaiDu) SetCheckType(checkType string) {
    _, ok := BdCheckTypes[checkType]
    if ok {
        m.checkType = checkType
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "校验类型不支持", nil))
    }
}

func (m *BaseBaiDu) SetSk(sk string) {
    match, _ := regexp.MatchString(`^[0-9a-z]{32}$`, sk)
    if match {
        m.sk = sk
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "用户签名不合法", nil))
    }
}

func (m *BaseBaiDu) SetReqRefer(reqRefer string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, reqRefer)
    if match {
        m.reqRefer = reqRefer
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "请求引用地址不合法", nil))
    }
}

func (m *BaseBaiDu) verifyData() {
    m.ReqData["ak"] = m.ak
    m.ReqData["output"] = m.output

    switch m.checkType {
    case BaiDuCheckTypeServerIp:
        m.ReqHeader["X-Forwarded-For"] = m.serverIp
        m.ReqHeader["Client-Ip"] = m.serverIp
    case BaiDuCheckTypeServerSn:
        if len(m.sk) == 0 {
            panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "签名校验码不能为空", nil))
        }
        snStr := m.serviceUri + "?" + mpf.HttpCreateParams(m.ReqData, "key", 1) + m.sk
        m.ReqData["sn"] = mpf.HashMd5(url.QueryEscape(snStr), "")
    default:
        if len(m.reqRefer) == 0 {
            panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "请求引用地址不能为空", nil))
        }
        m.ReqHeader["Referer"] = m.reqRefer
        m.ReqHeader["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.57 Safari/536.11"
    }
    m.ReqUrl += "?" + mpf.HttpCreateParams(m.ReqData, "key", 1)
}

func (m *BaseBaiDu) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    m.verifyData()

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(m.ReqUrl)
    req.Header.SetContentType(m.ReqContentType)
    req.Header.SetMethod(m.ReqMethod)

    referer, ok := m.ReqHeader["Referer"]
    if ok {
        req.Header.SetReferer(referer)
        delete(m.ReqHeader, "Referer")
    }
    userAgent, ok := m.ReqHeader["User-Agent"]
    if ok {
        req.Header.SetUserAgent(userAgent)
        delete(m.ReqHeader, "User-Agent")
    }
    mpf.HttpAddReqHeader(req, m.ReqHeader)

    return client, req
}

func NewBaseBaiDu() BaseBaiDu {
    conf := NewConfigBaiDu()
    base := BaseBaiDu{newBase(), "", "", "", "", "", "", "", ""}
    base.serviceDomain = "http://api.map.baidu.com"
    base.ak = conf.GetAk()
    base.serverIp = conf.GetServerIp()
    base.output = "json"
    base.checkType = BaiDuCheckTypeServerIp
    base.ReqMethod = fasthttp.MethodGet
    return base
}

type BaseGaoDe struct {
    baseMap
    serviceUri    string // 服务uri
    serviceDomain string // 服务域名
    key           string // 应用key
    secret        string // 应用密钥
}

func (m *BaseGaoDe) SetServiceUri(uri string) {
    m.serviceUri = uri
    m.ReqUrl = m.serviceDomain + uri
}

func (m *BaseGaoDe) createSign() {
    m.ReqData["key"] = m.key
    m.ReqData["output"] = "JSON"

    delete(m.ReqData, "sig")
    signStr := mpf.HttpCreateParams(m.ReqData, "key", 4)
    m.ReqData["sig"] = mpf.HashMd5(signStr+m.secret, "")
    m.ReqUrl += "?" + mpf.HttpCreateParams(m.ReqData, "key", 1)
}

func (m *BaseGaoDe) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    m.createSign()

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(m.ReqUrl)
    req.Header.SetContentType(m.ReqContentType)
    req.Header.SetMethod(m.ReqMethod)
    mpf.HttpAddReqHeader(req, m.ReqHeader)

    return client, req
}

func NewBaseGaoDe() BaseGaoDe {
    conf := NewConfigGaoDe()
    base := BaseGaoDe{newBase(), "", "", "", ""}
    base.serviceDomain = "https://restapi.amap.com/v3"
    base.key = conf.GetKey()
    base.secret = conf.GetSecret()
    base.ReqMethod = fasthttp.MethodGet
    return base
}

type BaseTencent struct {
    baseMap
    serviceUrl    string // 服务请求地址
    key           string // 应用key
    serverIp      string // 服务端IP
    webUrl        string // 页面URL
    appIdentifier string // 手机应用标识符
    output        string // 返回格式,默认JSON
    getType       string // 获取类型
}

func (t *BaseTencent) SetServiceUrl(serviceUrl string) {
    t.serviceUrl = serviceUrl
}

func (t *BaseTencent) SetAppIdentifier(appIdentifier string) {
    if len(appIdentifier) > 0 {
        t.appIdentifier = appIdentifier
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "应用标识符不能为空", nil))
    }
}

func (t *BaseTencent) SetGetType(getType string) {
    _, ok := TencentGetTypes[getType]
    if ok {
        t.getType = getType
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "获取类型不合法", nil))
    }
}

func (t *BaseTencent) verifyData() {
    t.ReqData["key"] = t.key

    switch t.getType {
    case TencentGetTypeBrowse:
        t.ReqHeader["Referer"] = t.webUrl
        t.ReqHeader["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.57 Safari/536.11"
    case TencentGetTypeMobile:
        if len(t.appIdentifier) == 0 {
            panic(mperr.NewMapTencent(errorcode.MapTencentParam, "应用标识符不能为空", nil))
        }
        t.ReqHeader["Referer"] = t.appIdentifier
    default:
        t.ReqHeader["X-Forwarded-For"] = t.serverIp
        t.ReqHeader["Client-Ip"] = t.serverIp
    }

    t.ReqUrl = t.serviceUrl + "?" + mpf.HttpCreateParams(t.ReqData, "key", 1)
}

func (t *BaseTencent) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    t.verifyData()

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(t.ReqUrl)
    req.Header.SetContentType(t.ReqContentType)
    req.Header.SetMethod(t.ReqMethod)

    referer, ok := t.ReqHeader["Referer"]
    if ok {
        req.Header.SetReferer(referer)
        delete(t.ReqHeader, "Referer")
    }
    userAgent, ok := t.ReqHeader["User-Agent"]
    if ok {
        req.Header.SetUserAgent(userAgent)
        delete(t.ReqHeader, "User-Agent")
    }
    mpf.HttpAddReqHeader(req, t.ReqHeader)

    return client, req
}

func NewBaseTencent() BaseTencent {
    conf := NewConfigTencent()
    base := BaseTencent{newBase(), "", "", "", "", "", "", ""}
    base.serverIp = conf.GetServerIp()
    base.webUrl = conf.GetDomain()
    base.key = conf.GetKey()
    base.getType = TencentGetTypeServer
    base.ReqData["output"] = "json"
    base.ReqMethod = fasthttp.MethodGet
    return base
}
