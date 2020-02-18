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

// 删除跨地域复制规则配置
type replicationDelete struct {
    qcloud.BaseCos
}

func (rd *replicationDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    rd.ReqURI = "http://" + rd.ReqHeader["Host"] + rd.ReqUri + "?replication"
    return rd.GetRequest()
}

func NewReplicationDelete() *replicationDelete {
    rd := &replicationDelete{qcloud.NewCos()}
    rd.SetParamData("replication", "")
    return rd
}
