/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 13:19
 */
package device

import (
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/valyala/fasthttp"
)

// 查询标签列表
type tagList struct {
    mppush.BaseJPush
}

func (tl *tagList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    tl.ReqUrl = tl.GetServiceUrl()

    return tl.GetRequest()
}

func NewTagList(key string) *tagList {
    tl := &tagList{mppush.NewBaseJPush(mppush.JPushServiceDomainDevice, key, "app")}
    tl.ServiceUri = "/v3/tags"
    return tl
}
