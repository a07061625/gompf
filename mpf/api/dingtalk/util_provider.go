package dingtalk

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

func (util *utilDingTalk) RefreshProviderSuiteTicket(suiteKey, suiteTicket string) {
    redisKey := project.RedisPrefix(project.RedisPrefixDingTalkProviderSuite) + suiteKey
    stData := make([]string, 0)
    stData = append(stData, redisKey, "unique_key", redisKey, "ticket", suiteTicket)
    cache.NewRedis().DoHmSet(stData)
    cache.NewRedis().GetConn().Expire(redisKey, 3600*time.Second)
}

// 获取服务商套件ticket
func (util *utilDingTalk) GetProviderSuiteTicket() string {
    conf := NewConfig().GetProvider()
    redisKey := project.RedisPrefix(project.RedisPrefixDingTalkProviderSuite) + conf.GetSuiteKey()
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    uniqueKey, ok := redisData["unique_key"]
    if ok && (uniqueKey == redisKey) {
        return redisData["ticket"]
    } else {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "获取服务商套件缓存失败", nil))
    }
}

func (util *utilDingTalk) refreshProviderAuthorizeAccessToken(corpId string) map[string]interface{} {
    nowTime := time.Now().Second()
    conf := NewConfig().GetProvider()
    suiteTicket := util.GetProviderSuiteTicket()
    signData := strconv.Itoa(nowTime) + "\n" + suiteTicket
    atMap := make(map[string]string)
    atMap["timestamp"] = strconv.Itoa(nowTime)
    atMap["accessKey"] = conf.GetSuiteKey()
    atMap["suiteTicket"] = suiteTicket
    atMap["signature"] = util.CreateApiSign(signData, conf.GetSuiteSecret())
    atUrl := UrlService + "/service/get_corp_token?" + mpf.HttpCreateParams(atMap, "none", 1)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(atUrl)
    req.Header.SetContentType(project.HttpContentTypeJson)
    req.Header.SetMethod(fasthttp.MethodPost)

    atData := make(map[string]string)
    atData["auth_corpid"] = corpId
    reqBody := mpf.JsonMarshal(atData)
    req.SetBody([]byte(reqBody))

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "获取服务商授权者访问令牌出错", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["access_token"]
    if !ok {
        panic(mperr.NewDingTalkProvider(errorcode.DingTalkProviderParam, "获取服务商授权者访问令牌出错", nil))
    }

    return respData
}

// 获取授权者访问令牌
func (util *utilDingTalk) GetProviderAuthorizeAccessToken(corpId string) string {
    nowTime := time.Now().Second()
    redisKey := project.RedisPrefix(project.RedisPrefixDingTalkProviderAuthorize) + corpId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    accessTokenKey, ok := redisData["at_key"]
    if ok && (accessTokenKey == redisKey) {
        expireTime, _ := strconv.Atoi(redisData["at_expire"])
        if expireTime >= nowTime {
            return redisData["at_content"]
        }
    }

    refreshRes := util.refreshProviderAuthorizeAccessToken(corpId)
    accessToken := refreshRes["access_token"].(string)
    expireTime := refreshRes["expires_in"].(int) + nowTime - 10
    atData := make([]string, 0)
    atData = append(atData, redisKey, "at_key", redisKey, "at_content", accessToken, "at_expire", strconv.Itoa(expireTime))
    cache.NewRedis().DoHmSet(atData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return accessToken
}

func (util *utilDingTalk) GetProviderSsoToken() string {
    nowTime := time.Now().Second()
    conf := NewConfig().GetProvider()
    redisKey := project.RedisPrefix(project.RedisPrefixDingTalkProviderAccount) + conf.GetCorpId()
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    ssoTokenKey, ok := redisData["sso_key"]
    if ok && (ssoTokenKey == redisKey) {
        expireTime, _ := strconv.Atoi(redisData["sso_expire"])
        if expireTime >= nowTime {
            return redisData["sso_content"]
        }
    }

    refreshRes := util.refreshSsoToken("")
    ssoToken := refreshRes["access_token"].(string)
    expireTime := nowTime + 7000
    stData := make([]string, 0)
    stData = append(stData, redisKey, "sso_key", redisKey, "sso_content", ssoToken, "sso_expire", strconv.Itoa(expireTime))
    cache.NewRedis().DoHmSet(stData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return ssoToken
}

func (util *utilDingTalk) GetProviderSnsToken() string {
    nowTime := time.Now().Second()
    conf := NewConfig().GetProvider()
    redisKey := project.RedisPrefix(project.RedisPrefixDingTalkProviderAccount) + conf.GetCorpId()
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    snsTokenKey, ok := redisData["sns_key"]
    if ok && (snsTokenKey == redisKey) {
        expireTime, _ := strconv.Atoi(redisData["sns_expire"])
        if expireTime >= nowTime {
            return redisData["sns_content"]
        }
    }

    refreshRes := util.refreshSnsToken("")
    snsToken := refreshRes["access_token"].(string)
    expireTime := nowTime + 7000
    stData := make([]string, 0)
    stData = append(stData, redisKey, "sns_key", redisKey, "sns_content", snsToken, "sns_expire", strconv.Itoa(expireTime))
    cache.NewRedis().DoHmSet(stData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return snsToken
}

func (util *utilDingTalk) GetProviderUserSnsToken(openid, persistentCode string) string {
    nowTime := time.Now().Second()
    conf := NewConfig().GetProvider()
    redisKey := project.RedisPrefix(project.RedisPrefixDingTalkProviderAccount) + conf.GetCorpId() + "_" + mpf.HashCrc32(openid, "")
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    snsTokenKey, ok := redisData["sns_key"]
    if ok && (snsTokenKey == redisKey) {
        expireTime, _ := strconv.Atoi(redisData["sns_expire"])
        if expireTime >= nowTime {
            return redisData["sns_content"]
        }
    }

    refreshRes := util.refreshUserSnsToken("", openid, persistentCode)
    snsToken := refreshRes["sns_token"].(string)
    expireTime := refreshRes["expires_in"].(int) + nowTime - 10
    stData := make([]string, 0)
    stData = append(stData, redisKey, "sns_key", redisKey, "sns_content", snsToken, "sns_expire", strconv.Itoa(expireTime))
    cache.NewRedis().DoHmSet(stData)
    cache.NewRedis().GetConn().Expire(redisKey, 8000*time.Second)
    return snsToken
}
