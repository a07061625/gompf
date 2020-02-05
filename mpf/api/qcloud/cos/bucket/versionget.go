/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 14:31
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 获取存储桶的版本控制配置
type versionGet struct {
    qcloud.BaseCos
}

func (vg *versionGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    vg.ReqUrl = "http://" + vg.ReqHeader["Host"] + vg.ReqUri + "?versioning"
    return vg.GetRequest()
}

func NewVersionGet() *versionGet {
    vg := &versionGet{qcloud.NewCos()}
    vg.SetParamData("versioning", "")
    return vg
}
