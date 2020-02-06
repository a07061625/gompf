/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/2 0002
 * Time: 13:11
 */
package wx

import (
    "crypto/tls"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/cache"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 刷新服务商凭证
func (util *utilWx) refreshProviderToken() map[string]interface{} {
    conf := NewConfig().GetProvider()
    ptMap := make(map[string]string)
    ptMap["corpid"] = conf.GetCorpId()
    ptMap["provider_secret"] = conf.GetCorpSecret()
    ptUrl := "https://qyapi.weixin.qq.com/cgi-bin/service/get_provider_token"

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(ptUrl)
    req.Header.SetContentType(project.HttpContentTypeJson)
    req.Header.SetMethod(fasthttp.MethodPost)

    reqBody := mpf.JsonMarshal(ptMap)
    req.SetBody([]byte(reqBody))

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "获取服务商凭证失败", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["provider_access_token"]
    if !ok {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "获取服务商凭证失败", nil))
    }

    return respData
}

// 获取服务商凭证
func (util *utilWx) GetProviderToken() string {
    nowTime := time.Now().Unix()
    conf := NewConfig().GetProvider()
    redisKey := project.RedisPrefix(project.RedisPrefixWxProviderAccount) + conf.GetCorpId()
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    accessTokenKey, ok := redisData["at_key"]
    if ok && (accessTokenKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["at_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["at_content"]
        }
    }

    refreshRes := util.refreshProviderToken()
    accessToken := refreshRes["provider_access_token"].(string)
    expireTime := refreshRes["expires_in"].(int64) + nowTime - 10
    ptData := make([]string, 0)
    ptData = append(ptData, redisKey, "at_key", redisKey, "at_expire", strconv.FormatInt(expireTime, 10), "at_content", accessToken)
    cache.NewRedis().DoHmSet(ptData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return accessToken
}

// 刷新服务商套件ticket
func (util *utilWx) RefreshProviderSuiteTicket(suiteId string, suiteTicket string) {
    redisKey := project.RedisPrefix(project.RedisPrefixWxProviderSuite) + suiteId
    stData := make([]string, 0)
    stData = append(stData, redisKey, "unique_key", redisKey, "ticket", suiteTicket)
    cache.NewRedis().DoHmSet(stData)
    cache.NewRedis().GetConn().Expire(redisKey, 1800*time.Second)
}

// 获取服务商套件ticket
func (util *utilWx) GetProviderSuiteTicket() string {
    conf := NewConfig().GetProvider()
    redisKey := project.RedisPrefix(project.RedisPrefixWxProviderSuite) + conf.GetSuiteId()
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    uniqueKey, ok := redisData["unique_key"]
    if ok && (uniqueKey == redisKey) {
        return redisData["ticket"]
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "获取微信服务商套件ticket失败", nil))
    }
}

// 刷新服务商第三方应用凭证
func (util *utilWx) refreshProviderSuiteToken() map[string]interface{} {
    conf := NewConfig().GetProvider()
    atMap := make(map[string]string)
    atMap["suite_id"] = conf.GetSuiteId()
    atMap["suite_secret"] = conf.GetSuiteSecret()
    atMap["suite_ticket"] = util.GetProviderSuiteTicket()
    atUrl := "https://qyapi.weixin.qq.com/cgi-bin/service/get_suite_token"

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(atUrl)
    req.Header.SetContentType(project.HttpContentTypeJson)
    req.Header.SetMethod(fasthttp.MethodPost)

    reqBody := mpf.JsonMarshal(atMap)
    req.SetBody([]byte(reqBody))

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "获取第三方应用凭证出错", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["suite_access_token"]
    if !ok {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "获取第三方应用凭证出错", nil))
    }

    return respData
}

// 获取服务商第三方应用凭证
func (util *utilWx) GetProviderSuiteToken() string {
    nowTime := time.Now().Unix()
    conf := NewConfig().GetProvider()
    redisKey := project.RedisPrefix(project.RedisPrefixWxProviderSuite) + conf.GetSuiteId()
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    accessTokenKey, ok := redisData["at_key"]
    if ok && (accessTokenKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["at_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["at_content"]
        }
    }

    refreshRes := util.refreshProviderSuiteToken()
    accessToken := refreshRes["suite_access_token"].(string)
    expireTime := refreshRes["expires_in"].(int64) + nowTime - 10
    atData := make([]string, 0)
    atData = append(atData, redisKey, "at_key", redisKey, "at_expire", strconv.FormatInt(expireTime, 10), "at_content", accessToken)
    cache.NewRedis().DoHmSet(atData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return accessToken
}

// 刷新服务商授权者信息
func (util *utilWx) RefreshProviderAuthorizeInfo(corpId string, operateType int, data map[string]interface{}) {
    util.outer.RefreshProviderAuthorizeInfo(corpId, operateType, data)
    redisKey := project.RedisPrefix(project.RedisPrefixWxProviderAuthorize) + corpId
    cache.NewRedis().GetConn().Del(redisKey).Result()
}

// 获取服务商授权者信息
func (util *utilWx) getProviderAuthorizeInfo(corpId string) *DataProviderAuthorize {
    refreshRes := util.outer.GetProviderAuthorizeInfo(corpId)
    if refreshRes.AuthorizeStatus == project.WxConfigAuthorizeStatusEmpty {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "授权企业微信不存在", nil))
    } else if refreshRes.AuthorizeStatus == project.WxConfigAuthorizeStatusNo {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "企业微信已取消授权", nil))
    }

    return refreshRes
}

// 获取服务商授权者永久授权码
func (util *utilWx) GetProviderAuthorizePermanentCode(corpId string) string {
    redisKey := project.RedisPrefix(project.RedisPrefixWxProviderAuthorize) + corpId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    uniqueKey, ok := redisData["unique_key"]
    if ok && (uniqueKey == redisKey) {
        return redisData["permanent_code"]
    }

    refreshRes := util.getProviderAuthorizeInfo(corpId)
    pcData := make([]string, 0)
    pcData = append(pcData, redisKey, "unique_key", redisKey, "permanent_code", refreshRes.PermanentCode)
    cache.NewRedis().DoHmSet(pcData)
    cache.NewRedis().GetConn().Expire(redisKey, 86400*time.Second)
    return refreshRes.PermanentCode
}

// 刷新服务商授权者访问令牌
func (util *utilWx) refreshProviderAuthorizeAccessToken(corpId string) map[string]interface{} {
    atMap := make(map[string]string)
    atMap["auth_corpid"] = corpId
    atMap["permanent_code"] = util.GetProviderAuthorizePermanentCode(corpId)
    atUrl := "https://qyapi.weixin.qq.com/cgi-bin/service/get_corp_token?suite_access_token=" + util.GetProviderSuiteToken()

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(atUrl)
    req.Header.SetContentType(project.HttpContentTypeJson)
    req.Header.SetMethod(fasthttp.MethodPost)

    reqBody := mpf.JsonMarshal(atMap)
    req.SetBody([]byte(reqBody))

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "获取服务商授权者访问令牌出错", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["access_token"]
    if !ok {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "获取服务商授权者访问令牌出错", nil))
    }

    return respData
}

// 获取服务商授权者访问令牌
func (util *utilWx) GetProviderAuthorizeAccessToken(corpId string) string {
    nowTime := time.Now().Unix()
    redisKey := project.RedisPrefix(project.RedisPrefixWxProviderAuthorize) + corpId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    accessTokenKey, ok := redisData["at_key"]
    if ok && (accessTokenKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["at_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["at_content"]
        }
    }

    refreshRes := util.refreshProviderAuthorizeAccessToken(corpId)
    accessToken := refreshRes["access_token"].(string)
    expireTime := refreshRes["expires_in"].(int64) + nowTime - 10
    atData := make([]string, 0)
    atData = append(atData, redisKey, "at_key", redisKey, "at_expire", strconv.FormatInt(expireTime, 10), "at_content", accessToken)
    cache.NewRedis().DoHmSet(atData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return accessToken
}

// 获取服务商授权者js ticket
func (util *utilWx) GetProviderAuthorizeJsTicket(corpId string) string {
    nowTime := time.Now().Unix()
    redisKey := project.RedisPrefix(project.RedisPrefixWxProviderAuthorize) + corpId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    jsTicketKey, ok := redisData["jt_key"]
    if ok && (jsTicketKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["jt_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["jt_content"]
        }
    }

    accessToken := util.GetProviderAuthorizeAccessToken(corpId)
    refreshRes := util.refreshCorpJsTicket(accessToken)
    jsTicket := refreshRes["ticket"].(string)
    expireTime := refreshRes["expires_in"].(int64) + nowTime - 10
    jtData := make([]string, 0)
    jtData = append(jtData, redisKey, "jt_key", redisKey, "jt_content", jsTicket, "jt_expire", strconv.FormatInt(expireTime, 10))
    cache.NewRedis().DoHmSet(jtData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return jsTicket
}
