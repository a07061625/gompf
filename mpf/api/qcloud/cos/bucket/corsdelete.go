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

// 删除存储桶的跨域访问配置信息
type corsDelete struct {
    qcloud.BaseCos
}

func (cd *corsDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    cd.ReqURI = "http://" + cd.ReqHeader["Host"] + cd.ReqUri + "?cors"
    return cd.GetRequest()
}

func NewCorsDelete() *corsDelete {
    cd := &corsDelete{qcloud.NewCos()}
    cd.SetParamData("cors", "")
    cd.ReqMethod = fasthttp.MethodDelete
    return cd
}
