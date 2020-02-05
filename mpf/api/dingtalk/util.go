package dingtalk

import (
    "crypto/tls"
    "encoding/base64"
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

type utilDingTalk struct {
    api.UtilApi
}

func (util *utilDingTalk) CreateApiSign(data, secret string) string {
    signStr := mpf.HashSha256(data, secret)
    return base64.StdEncoding.EncodeToString([]byte(signStr))
}

func (util *utilDingTalk) refreshSsoToken(corpId string) map[string]interface{} {
    atMap := make(map[string]string)
    if len(corpId) > 0 {
        conf := NewConfig().GetCorp(corpId)
        atMap["corpid"] = conf.GetCorpId()
        atMap["corpsecret"] = conf.GetSsoSecret()
    } else {
        conf := NewConfig().GetProvider()
        atMap["corpid"] = conf.GetCorpId()
        atMap["corpsecret"] = conf.GetSsoSecret()
    }
    atUrl := UrlService + "/sso/gettoken?" + mpf.HttpCreateParams(atMap, "none", 1)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(atUrl)
    req.Header.SetContentType(project.HttpContentTypeForm)
    req.Header.SetMethod(fasthttp.MethodGet)

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewDingTalk(errorcode.DingTalkParam, "获取sso token出错", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["access_token"]
    if !ok {
        panic(mperr.NewDingTalk(errorcode.DingTalkParam, "获取sso token出错", nil))
    }

    return respData
}

func (util *utilDingTalk) refreshSnsToken(corpId string) map[string]interface{} {
    stMap := make(map[string]string)
    if len(corpId) > 0 {
        conf := NewConfig().GetCorp(corpId)
        stMap["appid"] = conf.GetLoginAppId()
        stMap["appsecret"] = conf.GetLoginAppSecret()
    } else {
        conf := NewConfig().GetProvider()
        stMap["appid"] = conf.GetLoginAppId()
        stMap["appsecret"] = conf.GetLoginAppSecret()
    }
    stUrl := UrlService + "/sns/gettoken?" + mpf.HttpCreateParams(stMap, "none", 1)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(stUrl)
    req.Header.SetContentType(project.HttpContentTypeForm)
    req.Header.SetMethod(fasthttp.MethodGet)

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewDingTalk(errorcode.DingTalkParam, "获取开放应用令牌出错", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["access_token"]
    if !ok {
        panic(mperr.NewDingTalk(errorcode.DingTalkParam, "获取开放应用令牌出错", nil))
    }

    return respData
}

func (util *utilDingTalk) refreshUserSnsToken(corpId, openid, persistentCode string) map[string]interface{} {
    accessToken := ""
    if len(corpId) > 0 {
        accessToken = util.GetCorpSnsToken(corpId)
    } else {
        accessToken = util.GetProviderSnsToken()
    }
    stUrl := UrlService + "/sns/get_sns_token?access_token=" + url.QueryEscape(accessToken)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(stUrl)
    req.Header.SetContentType(project.HttpContentTypeJson)
    req.Header.SetMethod(fasthttp.MethodPost)

    stData := make(map[string]string)
    stData["openid"] = openid
    stData["persistent_code"] = persistentCode
    reqBody := mpf.JsonMarshal(stData)
    req.SetBody([]byte(reqBody))

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewDingTalk(errorcode.DingTalkParam, "获取用户授权令牌出错", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["sns_token"]
    if !ok {
        panic(mperr.NewDingTalk(errorcode.DingTalkParam, "获取用户授权令牌出错", nil))
    }

    return respData
}

func (util *utilDingTalk) DecryptMsg(encryptData string, encodeAesKey string, corpId string) string {
    return ""
}

func (util *utilDingTalk) EncryptMsg(replyMsg string, corpId string, encodeAesKey string) map[string]interface{} {
    result := make(map[string]interface{})
    return result
}

func (util *utilDingTalk) CreateCallbackSign(token, encryptData, nonceStr string, timestamp int) string {
    signMap := make(map[string]string)
    signMap["token"] = token
    signMap["encrypt"] = encryptData
    signMap["nonce"] = nonceStr
    signMap["timestamp"] = strconv.Itoa(timestamp)
    signData := mpf.NewHttpParamVal(signMap)
    sort.Sort(signData)

    signStr := ""
    for i := 0; i <= 3; i++ {
        signStr += signData.Params[i].Val
    }

    return mpf.HashSha1(signStr, "")
}

func (util *utilDingTalk) CheckCallbackSign(data map[string]string, sign string) bool {
    signTime, _ := strconv.Atoi(data["timestamp"])
    nowSign := util.CreateCallbackSign(data["token"], data["encrypt"], data["nonce"], signTime)

    return nowSign == sign
}

func (util *utilDingTalk) GetAccessToken(corpId, agentTag, atType string) string {
    if atType == AccessTokenTypeCorp {
        return util.GetCorpAccessToken(corpId, agentTag)
    } else {
        return util.GetProviderAuthorizeAccessToken(corpId)
    }
}

func (util *utilDingTalk) SendRequest(service api.IApiOuter, errorCode uint) api.ApiResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorCode
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

var (
    insUtil *utilDingTalk
)

func init() {
    insUtil = &utilDingTalk{api.NewUtilApi()}
}

func NewUtil() *utilDingTalk {
    return insUtil
}
