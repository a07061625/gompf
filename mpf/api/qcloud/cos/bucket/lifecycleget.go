/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 9:43
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 获取存储桶的生命周期设置
type lifeCycleGet struct {
    qcloud.BaseCos
}

func (lcg *lifeCycleGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    lcg.ReqUrl = "http://" + lcg.ReqHeader["Host"] + lcg.ReqUri + "?lifecycle"
    return lcg.GetRequest()
}

func NewLifeCycleGet() *lifeCycleGet {
    lcg := &lifeCycleGet{qcloud.NewCos()}
    lcg.SetParamData("lifecycle", "")
    return lcg
}
