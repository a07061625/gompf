/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 9:57
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 设置获取存储桶Referer的白名单或者黑名单
type refererGet struct {
    qcloud.BaseCos
}

func (rg *refererGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    rg.ReqUrl = "http://" + rg.ReqHeader["Host"] + rg.ReqUri + "?referer"
    return rg.GetRequest()
}

func NewRefererGet() *refererGet {
    rg := &refererGet{qcloud.NewCos()}
    rg.SetParamData("referer", "")
    return rg
}
