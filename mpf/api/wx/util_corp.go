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
    "github.com/a07061625/gompf/mpf/mpcache"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 刷新企业号访问令牌
func (util *utilWx) refreshCorpAccessToken(corpId string, agentTag string) map[string]interface{} {
    conf := NewConfig().GetCorp(corpId)
    agentInfo := conf.GetAgentInfo(agentTag)
    atMap := make(map[string]string)
    atMap["corpid"] = conf.GetCorpId()
    atMap["corpsecret"] = agentInfo["secret"]
    atUrl := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?" + mpf.HttpCreateParams(atMap, "none", 1)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(atUrl)
    req.Header.SetContentType(project.HTTPContentTypeForm)
    req.Header.SetMethod(fasthttp.MethodGet)

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "获取企业号access token失败", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["access_token"]
    if !ok {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, respData["errmsg"].(string), nil))
    }

    return respData
}

// 获取企业号访问令牌
func (util *utilWx) GetCorpAccessToken(corpId string, agentTag string) string {
    nowTime := time.Now().Unix()
    agentInfo := NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    redisKey := project.RedisPrefix(project.RedisPrefixWxCorp) + corpId + "_" + agentInfo["id"]
    redisData := mpcache.NewRedis().GetConn().HGetAll(redisKey).Val()
    accessTokenKey, ok := redisData["at_key"]
    if ok && (accessTokenKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["at_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["at_content"]
        }
    }

    refreshRes := util.refreshCorpAccessToken(corpId, agentTag)
    accessToken := refreshRes["access_token"].(string)
    expireTime := refreshRes["expires_in"].(int64) + nowTime - 10
    atData := make([]string, 0)
    atData = append(atData, redisKey, "at_key", redisKey, "at_content", accessToken, "at_expire", strconv.FormatInt(expireTime, 10))
    mpcache.NewRedis().DoHmSet(atData)
    mpcache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return accessToken
}

// 刷新企业号js ticket
func (util *utilWx) refreshCorpJsTicket(accessToken string) map[string]interface{} {
    jtUrl := "https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket?access_token=" + accessToken

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(jtUrl)
    req.Header.SetContentType(project.HTTPContentTypeForm)
    req.Header.SetMethod(fasthttp.MethodGet)

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "获取企业号js ticket失败", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok && (respData["errcode"].(int) != 0) {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, respData["errmsg"].(string), nil))
    }

    return respData
}

// 获取企业号js ticket
func (util *utilWx) GetCorpJsTicket(corpId string, agentTag string) string {
    nowTime := time.Now().Unix()
    agentInfo := NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    redisKey := project.RedisPrefix(project.RedisPrefixWxCorp) + corpId + "_" + agentInfo["id"]
    redisData := mpcache.NewRedis().GetConn().HGetAll(redisKey).Val()
    jsTicketKey, ok := redisData["jt_key"]
    if ok && (jsTicketKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["jt_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["jt_content"]
        }
    }

    accessToken := util.GetCorpAccessToken(corpId, agentTag)
    refreshRes := util.refreshCorpJsTicket(accessToken)
    jsTicket := refreshRes["ticket"].(string)
    expireTime := refreshRes["expires_in"].(int64) + nowTime - 10
    jtData := make([]string, 0)
    jtData = append(jtData, redisKey, "jt_key", redisKey, "jt_content", jsTicket, "jt_expire", strconv.FormatInt(expireTime, 10))
    mpcache.NewRedis().DoHmSet(jtData)
    mpcache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return jsTicket
}
