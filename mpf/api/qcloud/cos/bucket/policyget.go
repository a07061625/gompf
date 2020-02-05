/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 10:11
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 获取存储桶的权限策略
type policyGet struct {
    qcloud.BaseCos
}

func (pg *policyGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    pg.ReqUrl = "http://" + pg.ReqHeader["Host"] + pg.ReqUri + "?policy"
    return pg.GetRequest()
}

func NewPolicyGet() *policyGet {
    pg := &policyGet{qcloud.NewCos()}
    pg.SetParamData("policy", "")
    return pg
}
