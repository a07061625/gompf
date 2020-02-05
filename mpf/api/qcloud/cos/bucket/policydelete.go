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

// 删除存储桶的权限策略
type policyDelete struct {
    qcloud.BaseCos
}

func (pd *policyDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    pd.ReqUrl = "http://" + pd.ReqHeader["Host"] + pd.ReqUri + "?policy"
    return pd.GetRequest()
}

func NewPolicyDelete() *policyDelete {
    pd := &policyDelete{qcloud.NewCos()}
    pd.ReqMethod = fasthttp.MethodDelete
    pd.SetParamData("policy", "")
    return pd
}
