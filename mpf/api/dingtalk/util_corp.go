package dingtalk

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

func (util *utilDingTalk) refreshCorpAccessToken(corpId, agentTag string) map[string]interface{} {
    agentInfo := NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    atMap := make(map[string]string)
    atMap["appkey"] = agentInfo["key"]
    atMap["appsecret"] = agentInfo["secret"]
    atUrl := UrlService + "/gettoken?" + mpf.HTTPCreateParams(atMap, "none", 1)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(atUrl)
    req.Header.SetContentType(project.HTTPContentTypeForm)
    req.Header.SetMethod(fasthttp.MethodGet)

    resp := mpf.HTTPSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "获取企业号访问令牌出错", nil))
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["access_token"]
    if !ok {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "获取企业号访问令牌出错", nil))
    }

    return respData
}

// 获取企业号访问令牌
func (util *utilDingTalk) GetCorpAccessToken(corpId, agentTag string) string {
    nowTime := time.Now().Unix()
    agentInfo := NewConfig().GetCorp(corpId).GetAgentInfo(agentTag)
    redisKey := project.RedisPrefix(project.RedisPrefixDingTalkCorp) + corpId + "_" + agentInfo["id"]
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
    expireTime := nowTime + 7000
    atData := make([]string, 0)
    atData = append(atData, redisKey, "at_key", redisKey, "at_content", accessToken, "at_expire", strconv.FormatInt(expireTime, 10))
    mpcache.NewRedis().DoHmSet(atData)
    mpcache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return accessToken
}

func (util *utilDingTalk) GetCorpSsoToken(corpId string) string {
    nowTime := time.Now().Unix()
    redisKey := project.RedisPrefix(project.RedisPrefixDingTalkCorp) + corpId
    redisData := mpcache.NewRedis().GetConn().HGetAll(redisKey).Val()
    ssoTokenKey, ok := redisData["sso_key"]
    if ok && (ssoTokenKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["sso_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["sso_content"]
        }
    }

    refreshRes := util.refreshSsoToken(corpId)
    ssoToken := refreshRes["access_token"].(string)
    expireTime := nowTime + 7000
    stData := make([]string, 0)
    stData = append(stData, redisKey, "sso_key", redisKey, "sso_content", ssoToken, "sso_expire", strconv.FormatInt(expireTime, 10))
    mpcache.NewRedis().DoHmSet(stData)
    mpcache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return ssoToken
}

func (util *utilDingTalk) GetCorpSnsToken(corpId string) string {
    nowTime := time.Now().Unix()
    redisKey := project.RedisPrefix(project.RedisPrefixDingTalkCorp) + corpId
    redisData := mpcache.NewRedis().GetConn().HGetAll(redisKey).Val()
    snsTokenKey, ok := redisData["sns_key"]
    if ok && (snsTokenKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["sns_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["sns_content"]
        }
    }

    refreshRes := util.refreshSnsToken(corpId)
    snsToken := refreshRes["access_token"].(string)
    expireTime := nowTime + 7000
    stData := make([]string, 0)
    stData = append(stData, redisKey, "sns_key", redisKey, "sns_content", snsToken, "sns_expire", strconv.FormatInt(expireTime, 10))
    mpcache.NewRedis().DoHmSet(stData)
    mpcache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return snsToken
}

func (util *utilDingTalk) GetCorpUserSnsToken(corpId, openid, persistentCode string) string {
    nowTime := time.Now().Unix()
    redisKey := project.RedisPrefix(project.RedisPrefixDingTalkCorp) + corpId + "_" + mpf.HashCrc32(openid, "")
    redisData := mpcache.NewRedis().GetConn().HGetAll(redisKey).Val()
    snsTokenKey, ok := redisData["sns_key"]
    if ok && (snsTokenKey == redisKey) {
        expireTime, _ := strconv.ParseInt(redisData["sns_expire"], 10, 64)
        if expireTime >= nowTime {
            return redisData["sns_content"]
        }
    }

    refreshRes := util.refreshUserSnsToken(corpId, openid, persistentCode)
    snsToken := refreshRes["sns_token"].(string)
    expireTime := refreshRes["expires_in"].(int64) + nowTime - 10
    stData := make([]string, 0)
    stData = append(stData, redisKey, "sns_key", redisKey, "sns_content", snsToken, "sns_expire", strconv.FormatInt(expireTime, 10))
    mpcache.NewRedis().DoHmSet(stData)
    mpcache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return snsToken
}
