/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 9:27
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 删除指定的存储桶
type bucketDelete struct {
    qcloud.BaseCos
}

func (bd *bucketDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    bd.ReqUrl = "http://" + bd.ReqHeader["Host"] + bd.ReqUri
    return bd.GetRequest()
}

func NewBucketDelete() *bucketDelete {
    bd := &bucketDelete{qcloud.NewCos()}
    bd.ReqMethod = fasthttp.MethodDelete
    return bd
}
