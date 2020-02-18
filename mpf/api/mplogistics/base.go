/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 13:17
 */
package mplogistics

import (
    "crypto/tls"
    "encoding/base64"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type baseLogistics struct {
    api.ApiOuter
}

func newBaseLogistics() baseLogistics {
    return baseLogistics{api.NewApiOuter()}
}

type BaseAMAli struct {
    baseLogistics
    ServiceUri string // 服务uri
}

func (l *BaseAMAli) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    conf := NewConfigAMAli()
    l.ReqUrl = conf.serviceAddress + l.ServiceUri
    l.ReqHeader["Authorization"] = "APPCODE " + conf.GetAppCode()

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(l.ReqUrl)
    req.Header.SetContentType(l.ReqContentType)
    req.Header.SetMethod(l.ReqMethod)
    mpf.HTTPAddReqHeader(req, l.ReqHeader)

    return client, req
}

func NewBaseAMAli() BaseAMAli {
    l := BaseAMAli{newBaseLogistics(), ""}
    l.ReqContentType = project.HTTPContentTypeForm
    l.ReqMethod = fasthttp.MethodGet
    return l
}

type BaseKd100 struct {
    baseLogistics
    ExtendData map[string]interface{}
}

func (l *BaseKd100) createSign() {
    conf := NewConfigKd100()
    l.ReqData = make(map[string]string)
    l.ReqData["customer"] = conf.GetAppId()
    l.ReqData["param"] = mpf.JSONMarshal(l.ExtendData)
    signStr := l.ReqData["param"] + conf.GetAppKey() + conf.GetAppId()
    l.ReqData["sign"] = strings.ToUpper(mpf.HashMd5(signStr, ""))
    l.ReqUrl = "https://poll.kuaidi100.com/poll/query.do"
}

func (l *BaseKd100) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    l.createSign()

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(l.ReqUrl)
    req.Header.SetContentType(l.ReqContentType)
    req.Header.SetMethod(l.ReqMethod)
    mpf.HTTPAddReqHeader(req, l.ReqHeader)

    return client, req
}

func NewBaseKd100() BaseKd100 {
    l := BaseKd100{newBaseLogistics(), make(map[string]interface{})}
    l.ReqContentType = project.HTTPContentTypeForm
    l.ReqMethod = fasthttp.MethodPost
    return l
}

type BaseKdBird struct {
    baseLogistics
    ExtendData map[string]interface{}
}

func (l *BaseKdBird) createSign() {
    conf := NewConfigKdBird()
    l.ReqData["RequestData"] = mpf.JSONMarshal(l.ExtendData)
    signStr := mpf.HashMd5(l.ReqData["RequestData"]+conf.GetAppKey(), "")
    l.ReqData["DataSign"] = base64.StdEncoding.EncodeToString([]byte(signStr))
    l.ReqData["EBusinessID"] = conf.GetBusinessId()
    l.ReqData["DataType"] = "2"
    l.ReqContentType = project.HTTPContentTypeForm
    l.ReqMethod = fasthttp.MethodPost
}

func (l *BaseKdBird) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    l.createSign()

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(l.ReqUrl)
    req.Header.SetContentType(l.ReqContentType)
    req.Header.SetMethod(l.ReqMethod)
    mpf.HTTPAddReqHeader(req, l.ReqHeader)

    return client, req
}

func NewBaseKdBird() BaseKdBird {
    l := BaseKdBird{newBaseLogistics(), make(map[string]interface{})}
    return l
}
