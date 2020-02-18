package dingtalk

import (
    "crypto/tls"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type baseDingTalk struct {
    api.ApiOuter
    ExtendData map[string]interface{}
}

func (dt *baseDingTalk) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(dt.ReqUrl)
    req.Header.SetContentType(dt.ReqContentType)
    req.Header.SetMethod(dt.ReqMethod)
    mpf.HTTPAddReqHeader(req, dt.ReqHeader)

    return client, req
}

func newDingTalk() baseDingTalk {
    dt := baseDingTalk{api.NewApiOuter(), make(map[string]interface{})}
    dt.ReqContentType = project.HTTPContentTypeForm
    dt.ReqMethod = fasthttp.MethodGet
    return dt
}

type BaseCorp struct {
    baseDingTalk
}

func NewCorp() BaseCorp {
    return BaseCorp{newDingTalk()}
}

type BaseProvider struct {
    baseDingTalk
}

func NewProvider() BaseProvider {
    return BaseProvider{newDingTalk()}
}
