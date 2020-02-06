/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/2 0002
 * Time: 13:09
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

// 刷新公众号或小程序访问令牌
func (util *utilWx) refreshSingleAccessToken(appId string) map[string]interface{} {
    conf := NewConfig().GetAccount(appId)
    atMap := make(map[string]string)
    atMap["appid"] = conf.GetAppId()
    atMap["secret"] = conf.GetSecret()
    atMap["grant_type"] = "client_credential"
    atUrl := "https://api.weixin.qq.com/cgi-bin/token?" + mpf.HttpCreateParams(atMap, "none", 1)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(atUrl)
    req.Header.SetContentType(project.HttpContentTypeForm)
    req.Header.SetMethod(fasthttp.MethodGet)

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewWx(errorcode.WxParam, "获取公众号access token失败", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["access_token"]
    if !ok {
        panic(mperr.NewWx(errorcode.WxParam, respData["errmsg"].(string), nil))
    }

    return respData
}

// 获取公众号或小程序访问令牌
func (util *utilWx) GetSingleAccessToken(appId string) string {
    nowTime := time.Now().Unix()
    redisKey := project.RedisPrefix(project.RedisPrefixWxAccount) + appId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    accessTokenKey, ok := redisData["at_key"]
    if ok && (accessTokenKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["at_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["at_content"]
        }
    }

    refreshRes := util.refreshSingleAccessToken(appId)
    accessToken := refreshRes["access_token"].(string)
    expireTime := refreshRes["expires_in"].(int64) + nowTime - 10
    atData := make([]string, 0)
    atData = append(atData, redisKey, "at_key", redisKey, "at_content", accessToken, "at_expire", strconv.FormatInt(expireTime, 10))
    cache.NewRedis().DoHmSet(atData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return accessToken
}

// 刷新公众号或小程序ticket
func (util *utilWx) refreshSingleTicket(appId string, accessToken string, refreshType string) map[string]interface{} {
    ticketMap := make(map[string]string)
    ticketMap["access_token"] = accessToken
    ticketMap["type"] = refreshType
    ticketUrl := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?" + mpf.HttpCreateParams(ticketMap, "none", 1)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(ticketUrl)
    req.Header.SetContentType(project.HttpContentTypeForm)
    req.Header.SetMethod(fasthttp.MethodGet)

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewWx(errorcode.WxParam, "获取公众号ticket失败", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok && (respData["errcode"].(int) != 0) {
        panic(mperr.NewWx(errorcode.WxParam, respData["errmsg"].(string), nil))
    }

    return respData
}

// 获取公众号或小程序js ticket
func (util *utilWx) GetSingleJsTicket(appId string) string {
    nowTime := time.Now().Unix()
    redisKey := project.RedisPrefix(project.RedisPrefixWxAccount) + appId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    jsTicketKey, ok := redisData["jt_key"]
    if ok && (jsTicketKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["jt_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["jt_content"]
        }
    }

    accessToken := util.GetSingleAccessToken(appId)
    refreshRes := util.refreshSingleTicket(appId, accessToken, "jsapi")
    jsTicket := refreshRes["ticket"].(string)
    expireTime := refreshRes["expires_in"].(int64) + nowTime - 10
    jtData := make([]string, 0)
    jtData = append(jtData, redisKey, "jt_key", redisKey, "jt_content", jsTicket, "jt_expire", strconv.FormatInt(expireTime, 10))
    cache.NewRedis().DoHmSet(jtData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return jsTicket
}

// 获取公众号或小程序卡券 ticket
func (util *utilWx) GetSingleCardTicket(appId string) string {
    nowTime := time.Now().Unix()
    redisKey := project.RedisPrefix(project.RedisPrefixWxAccount) + appId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    cardTicketKey, ok := redisData["ct_key"]
    if ok && (cardTicketKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["ct_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["ct_content"]
        }
    }

    accessToken := util.GetSingleAccessToken(appId)
    refreshRes := util.refreshSingleTicket(appId, accessToken, "wx_card")
    cardTicket := refreshRes["ticket"].(string)
    expireTime := refreshRes["expires_in"].(int64) + nowTime - 10
    ctData := make([]string, 0)
    ctData = append(ctData, redisKey, "ct_key", redisKey, "ct_content", cardTicket, "ct_expire", strconv.FormatInt(expireTime, 10))
    cache.NewRedis().DoHmSet(ctData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return cardTicket
}
