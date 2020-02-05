/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 14:39
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 获取跨地域复制规则配置
type replicationGet struct {
    qcloud.BaseCos
}

func (rg *replicationGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    rg.ReqUrl = "http://" + rg.ReqHeader["Host"] + rg.ReqUri + "?replication"
    return rg.GetRequest()
}

func NewReplicationGet() *replicationGet {
    rg := &replicationGet{qcloud.NewCos()}
    rg.SetParamData("replication", "")
    return rg
}
