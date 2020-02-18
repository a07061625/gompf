/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 16:19
 */
package mpcurrency

import (
    "crypto/tls"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type baseCurrency struct {
    api.ApiOuter
}

func newBaseCurrency() baseCurrency {
    return baseCurrency{api.NewApiOuter()}
}

type BaseAMJiSu struct {
    baseCurrency
    ServiceUri string // 服务uri
}

func (c *BaseAMJiSu) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    conf := NewConfigAMJiSu()
    c.ReqUrl = conf.serviceAddress + c.ServiceUri
    c.ReqHeader["Authorization"] = "APPCODE " + conf.GetAppCode()

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(c.ReqUrl)
    req.Header.SetContentType(c.ReqContentType)
    req.Header.SetMethod(c.ReqMethod)
    mpf.HTTPAddReqHeader(req, c.ReqHeader)

    return client, req
}

func NewBaseAMJiSu() BaseAMJiSu {
    c := BaseAMJiSu{newBaseCurrency(), ""}
    c.ReqContentType = project.HTTPContentTypeForm
    c.ReqMethod = fasthttp.MethodGet
    return c
}

type BaseAMYiYuan struct {
    baseCurrency
    ServiceUri string // 服务uri
}

func (c *BaseAMYiYuan) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    conf := NewConfigAMYiYuan()
    c.ReqUrl = conf.serviceAddress + c.ServiceUri
    c.ReqHeader["Authorization"] = "APPCODE " + conf.GetAppCode()

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(c.ReqUrl)
    req.Header.SetContentType(c.ReqContentType)
    req.Header.SetMethod(c.ReqMethod)
    mpf.HTTPAddReqHeader(req, c.ReqHeader)

    return client, req
}

func NewBaseAMYiYuan() BaseAMYiYuan {
    c := BaseAMYiYuan{newBaseCurrency(), ""}
    c.ReqContentType = project.HTTPContentTypeForm
    c.ReqMethod = fasthttp.MethodGet
    return c
}
