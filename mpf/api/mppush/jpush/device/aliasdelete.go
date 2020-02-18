/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 13:19
 */
package device

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除别名
type aliasDelete struct {
    mppush.BaseJPush
    aliasName string // 别名
}

func (ad *aliasDelete) SetAliasName(aliasName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, aliasName)
    if match {
        ad.aliasName = aliasName
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "别名不合法", nil))
    }
}

func (ad *aliasDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ad.aliasName) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "别名不能为空", nil))
    }
    ad.ServiceUri = "/v3/aliases/" + ad.aliasName

    ad.ReqURI = ad.GetServiceUrl()

    return ad.GetRequest()
}

func NewAliasDelete(key string) *aliasDelete {
    ad := &aliasDelete{mppush.NewBaseJPush(mppush.JPushServiceDomainDevice, key, "app"), ""}
    ad.ReqMethod = fasthttp.MethodDelete
    return ad
}
