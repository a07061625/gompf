/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 16:29
 */
package mpprint

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

func (util *utilPrint) refreshFeYinAccessToken(appId string) map[string]interface{} {
    conf := NewConfigFeYin(appId)
    atMap := make(map[string]string)
    atMap["appid"] = conf.GetAppId()
    atMap["secret"] = conf.GetAppKey()
    atMap["code"] = conf.GetMemberCode()
    atUrl := FeYinServiceDomain + "/token?" + mpf.HttpCreateParams(atMap, "key", 1)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(atUrl)
    req.Header.SetContentType(project.HttpContentTypeForm)
    req.Header.SetMethod(fasthttp.MethodGet)

    resp := mpf.HttpSendReq(client, req, 3*time.Second)
    if resp.RespCode > 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinRequest, "获取access token失败", nil))
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["access_token"]
    if !ok {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinRequest, respData["errmsg"].(string), nil))
    }

    return respData
}

func (util *utilPrint) GetFeYinAccessToken(appId string) string {
    nowTime := time.Now().Unix()
    redisKey := project.RedisPrefix(project.RedisPrefixPrintFeiYinAccount) + appId
    redisData := cache.NewRedis().GetConn().HGetAll(redisKey).Val()
    uniqueKey, ok := redisData["unique_key"]
    if ok && (uniqueKey == redisKey) {
        eTime, _ := strconv.ParseInt(redisData["expire_time"], 10, 64)
        if eTime >= nowTime {
            return redisData["access_token"]
        }
    }

    refreshRes := util.refreshFeYinAccessToken(appId)
    expireTime := refreshRes["expires_in"].(int64) + nowTime
    activeTime := refreshRes["expires_in"].(time.Duration) + 100
    atData := make([]string, 0)
    atData = append(atData, redisKey, "access_token", refreshRes["access_token"].(string), "expire_time", strconv.FormatInt(expireTime, 10), "unique_key", redisKey)
    cache.NewRedis().DoHmSet(atData)
    cache.NewRedis().GetConn().Expire(redisKey, activeTime*time.Second)

    return refreshRes["access_token"].(string)
}
