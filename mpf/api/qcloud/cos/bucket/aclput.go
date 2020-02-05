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

// 设置存储桶的访问权限控制
type aclPut struct {
    qcloud.BaseCos
}

func (ap *aclPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    ap.ReqUrl = "http://" + ap.ReqHeader["Host"] + ap.ReqUri + "?acl"
    return ap.GetRequest()
}

func NewAclPut() *aclPut {
    ap := &aclPut{qcloud.NewCos()}
    ap.SetParamData("acl", "")
    ap.ReqMethod = fasthttp.MethodPut
    return ap
}
