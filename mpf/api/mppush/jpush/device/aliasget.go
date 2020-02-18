/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 13:19
 */
package device

import (
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询别名
type aliasGet struct {
    mppush.BaseJPush
    aliasName string // 别名
}

func (ag *aliasGet) SetAliasName(aliasName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, aliasName)
    if match {
        ag.aliasName = aliasName
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "别名不合法", nil))
    }
}

func (ag *aliasGet) SetPlatformList(platformList []string) {
    if len(platformList) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "平台类型不合法", nil))
    }

    platforms := make([]string, 0)
    for _, v := range platformList {
        _, ok := mppush.JPushPlatformTypes[v]
        if ok {
            platforms = append(platforms, v)
        }
    }
    if len(platforms) > 0 {
        ag.ReqData["platform"] = strings.Join(platforms, ",")
    }
}

func (ag *aliasGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ag.aliasName) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "别名不能为空", nil))
    }
    ag.ServiceUri = "/v3/aliases/" + ag.aliasName

    ag.ReqURI = ag.GetServiceUrl()
    if len(ag.ReqData) > 0 {
        ag.ReqURI += "?" + mpf.HTTPCreateParams(ag.ReqData, "none", 1)
    }

    return ag.GetRequest()
}

func NewAliasGet(key string) *aliasGet {
    ag := &aliasGet{mppush.NewBaseJPush(mppush.JPushServiceDomainDevice, key, "app"), ""}
    return ag
}
