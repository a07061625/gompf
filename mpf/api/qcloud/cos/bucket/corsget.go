/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 22:56
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 获取存储桶的跨域访问配置信息
type corsGet struct {
    qcloud.BaseCos
}

func (cg *corsGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    cg.ReqURI = "http://" + cg.ReqHeader["Host"] + cg.ReqUri + "?cors"
    return cg.GetRequest()
}

func NewCorsGet() *corsGet {
    cg := &corsGet{qcloud.NewCos()}
    cg.SetParamData("cors", "")
    return cg
}
