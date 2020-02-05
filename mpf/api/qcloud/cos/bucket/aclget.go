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

// 获取存储桶的访问权限控制列表
type aclGet struct {
    qcloud.BaseCos
}

func (ag *aclGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    ag.ReqUrl = "http://" + ag.ReqHeader["Host"] + ag.ReqUri + "?acl"
    return ag.GetRequest()
}

func NewAclGet() *aclGet {
    ag := &aclGet{qcloud.NewCos()}
    ag.SetParamData("acl", "")
    return ag
}
