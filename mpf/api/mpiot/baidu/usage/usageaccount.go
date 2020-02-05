/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:05
 */
package usage

import (
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/valyala/fasthttp"
)

// 获取当前账单月内使用量
type usageAccount struct {
    mpiot.BaseBaiDu
}

func (ua *usageAccount) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    ua.ReqUrl = ua.GetServiceUrl()

    return ua.GetRequest()
}

func NewUsageAccount() *usageAccount {
    ua := &usageAccount{mpiot.NewBaseBaiDu()}
    ua.ServiceUri = "/v1/usage"
    return ua
}
