/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/2 0002
 * Time: 13:10
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

// 刷新开放平台授权者信息
func (util *utilWx) RefreshOpenAuthorizeInfo(appId string, operateType int, data map[string]interface{}) {
    util.outer.RefreshOpenAuthorizeInfo(appId, operateType, data)
    redisKey := project.RedisPrefix(project.RedisPrefixWxOpenAuthorize) + appId
    cache.NewRedis().GetConn().Del(redisKey).Result()
}

// 刷新开放平台访问令牌
func (util *utilWx) RefreshOpenAccessToken(verifyTicket string) {
    conf := NewConfig().GetOpen()
    atMap := make(map[string]string)
    atMap["component_appid"] = conf.GetAppId()
    atMap["component_appsecret"] = conf.GetSecret()
    atMap["component_verify_ticket"] = verifyTicket
    atUrl := "https://api.weixin.qq.com/cgi-bin/component/api_component_token"

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
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "获取第三方开放平台access token失败", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["component_access_token"]
    if !ok {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "获取第三方开放平台access token失败", nil))
    }

    redisKey := project.RedisPrefix(project.RedisPrefixWxOpenAccount) + conf.GetAppId()
    expireTime := respData["expires_in"].(int)
    atData := make([]string, 0)
    atData = append(atData, redisKey, "unique_key", redisKey, "access_token", respData["component_access_token"].(string))
    cache.NewRedis().DoHmSet(atData)
    cache.NewRedis().GetConn().Expire(redisKey, time.Duration(expireTime)*time.Second)
}

// 获取开放平台访问令牌
func (util *utilWx) GetOpenAccessToken() string {
    conf := NewConfig().GetOpen()
    redisKey := project.RedisPrefix(project.RedisPrefixWxOpenAccount) + conf.GetAppId()
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    uniqueKey, ok := redisData["unique_key"]
    if ok && (uniqueKey == redisKey) {
        return redisData["access_token"]
    } else {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "获取第三方开放平台access token失败", nil))
    }
}

// 刷新开放平台授权者访问令牌
func (util *utilWx) refreshOpenAuthorizeAccessToken(appId string) map[string]interface{} {
    redisKey := project.RedisPrefix(project.RedisPrefixWxOpenAuthorize) + appId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    refreshToken, ok := redisData["refresh_token"]
    authorizeStatus := 0
    if ok {
        authorizeStatus, _ = strconv.Atoi(redisData["authorize_status"])
    } else {
        getRes := util.outer.GetOpenAuthorizeInfo(appId)
        refreshToken = getRes.RefreshToken
        authorizeStatus = getRes.AuthorizeStatus
        authorizeData := make([]string, 0)
        authorizeData = append(authorizeData, "authorize_status", strconv.Itoa(getRes.AuthorizeStatus), "auth_code", getRes.AuthCode, "refresh_token", getRes.RefreshToken)
        cache.NewRedis().DoHmSet(authorizeData)
        cache.NewRedis().GetConn().Expire(redisKey, 86400*time.Second)
    }

    if authorizeStatus == project.WxConfigAuthorizeStatusEmpty {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "授权公众号不存在", nil))
    } else if authorizeStatus == project.WxConfigAuthorizeStatusNo {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "公众号已取消授权", nil))
    }

    conf := NewConfig().GetOpen()
    atMap := make(map[string]string)
    atMap["component_appid"] = conf.GetAppId()
    atMap["authorizer_appid"] = appId
    atMap["authorizer_refresh_token"] = refreshToken
    atUrl := "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=" + util.GetOpenAccessToken()

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
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "获取授权者access token失败", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok = respData["authorizer_access_token"]
    if !ok {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "获取授权者access token失败", nil))
    }

    return respData
}

// 获取开放平台授权者访问令牌
func (util *utilWx) GetOpenAuthorizeAccessToken(appId string) string {
    nowTime := time.Now().Second()
    redisKey := project.RedisPrefix(project.RedisPrefixWxOpenAuthorize) + appId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    accessTokenKey, ok := redisData["at_key"]
    if ok && (accessTokenKey == redisKey) {
        expireTime, _ := strconv.Atoi(redisData["at_expire"])
        if expireTime >= nowTime {
            return redisData["at_content"]
        }
    }

    refreshRes := util.refreshOpenAuthorizeAccessToken(appId)
    accessToken := refreshRes["authorizer_access_token"].(string)
    expireTime := refreshRes["expires_in"].(int) + nowTime - 10
    atData := make([]string, 0)
    atData = append(atData, redisKey, "at_key", redisKey, "at_expire", strconv.Itoa(expireTime), "at_content", accessToken)
    cache.NewRedis().DoHmSet(atData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return accessToken
}

// 获取开放平台授权者js ticket
func (util *utilWx) GetOpenAuthorizeJsTicket(appId string) string {
    nowTime := time.Now().Second()
    redisKey := project.RedisPrefix(project.RedisPrefixWxOpenAuthorize) + appId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    jsTicketKey, ok := redisData["jt_key"]
    if ok && (jsTicketKey == redisKey) {
        expireTime, _ := strconv.Atoi(redisData["jt_expire"])
        if expireTime >= nowTime {
            return redisData["jt_content"]
        }
    }

    accessToken := util.GetOpenAuthorizeAccessToken(appId)
    refreshRes := util.refreshSingleTicket(appId, accessToken, "jsapi")
    jsTicket := refreshRes["ticket"].(string)
    expireTime := refreshRes["expires_in"].(int) + nowTime - 10
    jtData := make([]string, 0)
    jtData = append(jtData, redisKey, "jt_key", redisKey, "jt_content", jsTicket, "jt_expire", strconv.Itoa(expireTime))
    cache.NewRedis().DoHmSet(jtData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return jsTicket
}

// 获取开放平台授权者卡券ticket
func (util *utilWx) GetOpenAuthorizeCardTicket(appId string) string {
    nowTime := time.Now().Second()
    redisKey := project.RedisPrefix(project.RedisPrefixWxOpenAuthorize) + appId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    cardTicketKey, ok := redisData["ct_key"]
    if ok && (cardTicketKey == redisKey) {
        expireTime, _ := strconv.Atoi(redisData["ct_expire"])
        if expireTime >= nowTime {
            return redisData["ct_content"]
        }
    }

    accessToken := util.GetOpenAuthorizeAccessToken(appId)
    refreshRes := util.refreshSingleTicket(appId, accessToken, "wx_card")
    cardTicket := refreshRes["ticket"].(string)
    expireTime := refreshRes["expires_in"].(int) + nowTime - 10
    ctData := make([]string, 0)
    ctData = append(ctData, redisKey, "ct_key", redisKey, "ct_content", cardTicket, "ct_expire", strconv.Itoa(expireTime))
    cache.NewRedis().DoHmSet(ctData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return cardTicket
}

// 刷新开放平台小程序云开发代码保护密钥
func (util *utilWx) refreshOpenAuthorizeCodeSecret(appId string) map[string]interface{} {
    csUrl := "https://api.weixin.qq.com/tcb/getcodesecret?access_token=?" + util.GetOpenAuthorizeAccessToken(appId)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(csUrl)
    req.Header.SetContentType(project.HttpContentTypeJson)
    req.Header.SetMethod(fasthttp.MethodPost)
    req.SetBody([]byte("{}"))

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "获取小程序代码保护密钥失败", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok && (respData["errcode"].(int) != 0) {
        panic(mperr.NewWxOpen(errorcode.WxOpenParam, "获取小程序代码保护密钥失败", nil))
    }

    return respData
}

// 获取开放平台小程序云开发代码保护密钥
func (util *utilWx) GetOpenAuthorizeCodeSecret(appId string) string {
    redisKey := project.RedisPrefix(project.RedisPrefixWxOpenAuthorize) + appId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    codeSecretKey, ok := redisData["cs_key"]
    if ok && (codeSecretKey == redisKey) {
        return redisData["cs_content"]
    }

    refreshRes := util.refreshOpenAuthorizeCodeSecret(appId)
    codeSecret := refreshRes["codesecret"].(string)
    ctData := make([]string, 0)
    ctData = append(ctData, redisKey, "cs_key", redisKey, "cs_content", codeSecret)
    cache.NewRedis().DoHmSet(ctData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return codeSecret
}
