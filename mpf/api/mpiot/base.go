/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 15:14
 */
package mpiot

import (
    "crypto/tls"
    "net/url"
    "strconv"
    "strings"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type baseIot struct {
    api.ApiOuter
    ExtendData map[string]interface{}
}

func newBaseIot() baseIot {
    return baseIot{api.NewApiOuter(), make(map[string]interface{})}
}

type BaseBaiDu struct {
    baseIot
    serviceProtocol string // 服务协议
    serviceDomain   string // 服务域名
    ServiceUri      string // 服务URI
}

func (bd *BaseBaiDu) SetServiceProtocol(serviceProtocol string) {
    if (serviceProtocol == "http") || (serviceProtocol == "https") {
        bd.serviceProtocol = serviceProtocol
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "服务协议不合法", nil))
    }
}

func (bd *BaseBaiDu) SetServiceDomain(serviceDomain string) {
    _, ok := BaiDuDomains[serviceDomain]
    if ok {
        bd.serviceDomain = serviceDomain
        bd.ReqHeader["Host"] = serviceDomain
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "服务域名不合法", nil))
    }
}

func (bd *BaseBaiDu) GetServiceUrl() string {
    return bd.serviceProtocol + "://" + bd.serviceDomain + bd.ServiceUri
}

func (bd *BaseBaiDu) createSign() string {
    needStr := bd.ReqMethod + "\n" + url.QueryEscape(bd.ServiceUri)

    delete(bd.ReqData, "authorization")
    needStr += "\n" + mpf.HTTPCreateParams(bd.ReqData, "key", 1)

    reqHeader := "host"
    needStr += "\n" + reqHeader

    conf := NewConfig().GetBaiDu()
    nowTime := time.Now().UTC()
    authPrefix := "bce-auth-v1/" + conf.GetAccessKey() + "/" + nowTime.Format("2006-01-02T03:04:05Z") + "/30"
    signKey := mpf.HashSha256(authPrefix, conf.GetAccessSecret())
    return authPrefix + "/" + reqHeader + "/" + mpf.HashSha256(needStr, signKey)
}

func (bd *BaseBaiDu) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    bd.ReqHeader["Authorization"] = bd.createSign()

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(bd.ReqUrl)
    req.Header.SetContentType(bd.ReqContentType)
    req.Header.SetMethod(bd.ReqMethod)
    mpf.HTTPAddReqHeader(req, bd.ReqHeader)

    return client, req
}

func NewBaseBaiDu() BaseBaiDu {
    base := BaseBaiDu{newBaseIot(), "", "", ""}
    base.serviceProtocol = "https"
    base.serviceDomain = BaiDuDomainGZ
    base.ReqHeader["Host"] = base.serviceDomain
    base.ReqHeader["Content-Type"] = "application/json; charset=utf-8"
    base.ReqContentType = project.HTTPContentTypeJSON
    base.ReqMethod = fasthttp.MethodGet
    return base
}

type BaseTencent struct {
    baseIot
    serviceDomain string // 服务域名
    serviceName   string // 服务名称
}

func (t *BaseTencent) SetServiceDomain(serviceDomain string) {
    if len(serviceDomain) > 0 {
        t.serviceDomain = serviceDomain
        t.ReqHeader["Host"] = serviceDomain
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "服务域名不合法", nil))
    }
}

func (t *BaseTencent) createTC3Sign(reqBody string) {
    conf := NewConfig().GetTencent()
    t.ReqHeader["X-TC-Region"] = conf.GetRegionId()
    now := time.Now()
    nowTime := now.Unix()
    dateStr := now.Format("2006-01-02")
    signHeaders := "content-type;host"
    canonicalRequest := "POST\n/\n\n"
    canonicalRequest += "content-type:" + strings.ToLower(t.ReqHeader["Content-Type"]) + "\n"
    canonicalRequest += "host:" + strings.ToLower(t.ReqHeader["Host"]) + "\n\n"
    canonicalRequest += signHeaders + "\n"
    canonicalRequest += mpf.HashSha256(reqBody, "")
    credentialScope := dateStr + "/" + t.serviceName + "/tc3_request"
    signStr := "TC3-HMAC-SHA256\n"
    signStr += strconv.FormatInt(nowTime, 10) + "\n"
    signStr += credentialScope + "\n"
    signStr += mpf.HashSha256(canonicalRequest, "")

    secretDate := mpf.HashSha256(dateStr, "TC3"+conf.GetSecretKey())
    secretService := mpf.HashSha256(t.serviceName, secretDate)
    secretSign := mpf.HashSha256("tc3_request", secretService)
    signature := mpf.HashSha256(signStr, secretSign)
    t.ReqHeader["X-TC-Timestamp"] = strconv.FormatInt(nowTime, 10)
    t.ReqHeader["Authorization"] = "TC3-HMAC-SHA256 Credential=" + conf.GetSecretId() + "/" + credentialScope + ", SignedHeaders=" + signHeaders + ", Signature=" + signature
}

func (t *BaseTencent) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    reqBody := mpf.HTTPCreateParams(t.ReqData, "none", 1)
    t.createTC3Sign(reqBody)

    t.ReqUrl = "https://" + t.serviceDomain
    t.ReqContentType = project.HTTPContentTypeJSON
    t.ReqMethod = fasthttp.MethodPost

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(t.ReqUrl)
    req.Header.SetContentType(t.ReqContentType)
    req.Header.SetMethod(t.ReqMethod)
    req.SetBody([]byte(reqBody))
    mpf.HTTPAddReqHeader(req, t.ReqHeader)

    return client, req
}

func NewBaseTencent() BaseTencent {
    base := BaseTencent{newBaseIot(), "", ""}
    base.serviceName = "iot"
    base.serviceDomain = "iotexplorer.tencentcloudapi.com"
    base.ReqHeader["X-TC-Version"] = "2019-04-23"
    base.ReqHeader["Host"] = base.serviceDomain
    base.ReqHeader["Content-Type"] = "application/json"
    base.ReqHeader["Expect"] = ""
    return base
}
