/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 9:58
 */
package mpim

import (
    "crypto/tls"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/valyala/fasthttp"
)

type baseIM struct {
    api.ApiOuter
}

func newBaseIM() baseIM {
    return baseIM{api.NewApiOuter()}
}

type BaseTencent struct {
    baseIM
    ServiceUri string
    ExtendData map[string]interface{}
}

func (im BaseTencent) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    randNum := mpf.ToolCreateRandNum(10000000, 89999999)
    conf := NewConfigTencent()
    im.ReqData = make(map[string]string)
    im.ReqData["contenttype"] = "json"
    im.ReqData["random"] = strconv.Itoa(randNum)
    im.ReqData["sdkappid"] = conf.GetAppId()
    im.ReqData["identifier"] = conf.GetAccountAdmin()
    im.ReqData["usersig"] = NewUtilIM().GetTencentAccountSign(conf.GetAccountAdmin())
    im.ReqUrl = "https://console.tim.qq.com/v4" + im.ServiceUri + "?" + mpf.HttpCreateParams(im.ReqData, "none", 1)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(im.ReqUrl)
    req.Header.SetContentType(im.ReqContentType)
    req.Header.SetMethod(im.ReqMethod)
    mpf.HttpAddReqHeader(req, im.ReqHeader)

    return client, req
}

func NewBaseTencent() BaseTencent {
    return BaseTencent{newBaseIM(), "", make(map[string]interface{})}
}
