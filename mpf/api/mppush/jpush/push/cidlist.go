/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 10:35
 */
package push

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/cache"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 推送唯一标识符列表
type cidList struct {
    mppush.BaseJPush
}

func (cl *cidList) SetCount(count int) {
    if (count > 0) && (count <= 1000) {
        cl.ReqData["count"] = strconv.Itoa(count)
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标识符数量不合法", nil))
    }
}

func (cl *cidList) SetCidType(cidType string) {
    if (cidType == "push") || (cidType == "schedule") {
        cl.ReqData["type"] = cidType
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标识符类型不合法", nil))
    }
}

func (cl *cidList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    cl.ReqUrl = cl.GetServiceUrl() + "?" + mpf.HttpCreateParams(cl.ReqData, "none", 1)

    return cl.GetRequest()
}

func NewCidList(key string) *cidList {
    cl := &cidList{mppush.NewBaseJPush(mppush.JPushServiceDomainApi, key, "app")}
    cl.ServiceUri = "/v3/push/cid"
    cl.ReqData["count"] = "1"
    cl.ReqData["type"] = "push"
    return cl
}

// 获取APP唯一标识符
//   key: 应用标识
//   cidType: 类型 push:推送 schedule:定时任务
func GetAppCid(key, cidType string) string {
    redisKey := ""
    if cidType == "push" {
        redisKey = project.RedisPrefix(project.RedisPrefixJpushUidPush) + key
    } else {
        redisKey = project.RedisPrefix(project.RedisPrefixJpushUidSchedule) + key
    }

    cid, err := cache.NewRedis().GetConn().LPop(redisKey).Result()
    if err != nil {
        return cid
    }

    cl := NewCidList(key)
    cl.SetCidType(cidType)
    cl.SetCount(800)
    sendRes := mppush.NewUtil().SendJPushRequest(cl, errorcode.PushJPushRequestGet)
    if sendRes.Code > 0 {
        panic(mperr.NewPushJPush(sendRes.Code, sendRes.Msg, nil))
    }

    resData := sendRes.Data.(map[string]interface{})
    idList := resData["idlist"].([]string)
    idNum := len(idList)

    p := cache.NewRedis().GetConn().Pipeline()
    defer p.Close()
    for i := 0; i < idNum; i++ {
        if i > 0 {
            p.RPush(redisKey, idList[i])
        } else {
            cid = idList[i]
        }
    }
    p.Exec()
    cache.NewRedis().GetConn().Expire(redisKey, 86400*time.Second)

    return cid
}
