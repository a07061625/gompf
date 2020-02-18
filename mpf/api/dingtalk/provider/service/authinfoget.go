package service

import (
    "regexp"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取授权企业基本信息
type authInfoGet struct {
    dingtalk.BaseProvider
    authCorpId string // 授权企业ID
}

func (aig *authInfoGet) SetAuthCorpId(authCorpId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, authCorpId)
    if match {
        aig.authCorpId = authCorpId
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "授权企业ID不合法", nil))
    }
}

func (aig *authInfoGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(aig.authCorpId) == 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "授权企业ID不合法", nil))
    }
    aig.ExtendData["auth_corpid"] = aig.authCorpId

    providerConf := dingtalk.NewConfig().GetProvider()
    suiteTicket := dingtalk.NewUtil().GetProviderSuiteTicket()
    nowTime := strconv.FormatInt(time.Now().Unix(), 10)

    aig.ReqData["timestamp"] = nowTime
    aig.ReqData["accessKey"] = providerConf.GetSuiteKey()
    aig.ReqData["suiteTicket"] = suiteTicket
    aig.ReqData["signature"] = dingtalk.NewUtil().CreateApiSign(nowTime+"\n"+suiteTicket, providerConf.GetSuiteSecret())
    aig.ReqURI = dingtalk.UrlService + "/service/get_auth_info?" + mpf.HTTPCreateParams(aig.ReqData, "none", 1)

    reqBody := mpf.JSONMarshal(aig.ReqData)
    client, req := aig.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewAuthInfoGet() *authInfoGet {
    aig := &authInfoGet{dingtalk.NewProvider(), ""}
    aig.ReqContentType = project.HTTPContentTypeJSON
    aig.ReqMethod = fasthttp.MethodPost
    return aig
}
