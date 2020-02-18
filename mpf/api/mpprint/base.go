/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 19:46
 */
package mpprint

import (
    "crypto/tls"
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type basePrint struct {
    api.ApiInner
}

func newBase() basePrint {
    return basePrint{api.NewApiInner()}
}

type BaseFeYin struct {
    basePrint
    appId string // 应用ID
}

func (p *BaseFeYin) SetAppId(appId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appId)
    if match {
        p.appId = appId
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "应用ID不合法", nil))
    }
}

func (p *BaseFeYin) GetAppId() string {
    return p.appId
}

func (p *BaseFeYin) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(p.ReqUrl)
    req.Header.SetContentType(p.ReqContentType)
    req.Header.SetMethod(p.ReqMethod)
    mpf.HTTPAddReqHeader(req, p.ReqHeader)

    return client, req
}

func NewBaseFeYin() BaseFeYin {
    return BaseFeYin{newBase(), ""}
}
