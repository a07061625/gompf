/**
 * 短信基础结构体
 * User: 姜伟
 * Date: 2019/12/22 0022
 * Time: 11:18
 */
package mpsms

import (
    "crypto/tls"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type BaseYun253 struct {
    api.ApiOuter
    Account  string // API账号
    Password string // API密码
}

func (b *BaseYun253) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    reqBody := mpf.HttpCreateParams(b.ReqData, "key", 2)
    req := fasthttp.AcquireRequest()
    req.SetBody([]byte(reqBody))
    req.Header.SetRequestURI(b.ReqUrl)
    req.Header.SetContentType(b.ReqContentType)
    req.Header.SetMethod(b.ReqMethod)
    mpf.HttpAddReqHeader(req, b.ReqHeader)

    return client, req
}

func NewBaseYun253() BaseYun253 {
    base := BaseYun253{api.NewApiOuter(), "", ""}
    base.ReqMethod = fasthttp.MethodPost
    base.ReqContentType = project.HTTPContentTypeJSON
    base.ReqHeader["Expect"] = ""
    return base
}
