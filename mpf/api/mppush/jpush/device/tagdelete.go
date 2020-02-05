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

// 删除标签
type tagDelete struct {
    mppush.BaseJPush
    tag string // 标签名
}

func (td *tagDelete) SetTag(tag string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tag)
    if match {
        td.tag = tag
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标签名不合法", nil))
    }
}

func (td *tagDelete) SetPlatformList(platformList []string) {
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
        td.ReqData["platform"] = strings.Join(platforms, ",")
    }
}

func (td *tagDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(td.tag) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标签名不能为空", nil))
    }
    td.ServiceUri = "/v3/tags/" + td.tag

    td.ReqUrl = td.GetServiceUrl()
    if len(td.ReqData) > 0 {
        td.ReqUrl += "?" + mpf.HttpCreateParams(td.ReqData, "none", 1)
    }

    return td.GetRequest()
}

func NewTagDelete(key string) *tagDelete {
    td := &tagDelete{mppush.NewBaseJPush(mppush.JPushServiceDomainDevice, key, "app"), ""}
    td.ReqMethod = fasthttp.MethodDelete
    return td
}
