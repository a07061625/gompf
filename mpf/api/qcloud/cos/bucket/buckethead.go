/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 9:23
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 确认存储桶是否存在,是否有权限访问
type bucketHead struct {
    qcloud.BaseCos
}

func (bh *bucketHead) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    bh.ReqUrl = "http://" + bh.ReqHeader["Host"] + bh.ReqUri
    return bh.GetRequest()
}

func NewBucketHead() *bucketHead {
    bh := &bucketHead{qcloud.NewCos()}
    bh.ReqMethod = fasthttp.MethodHead
    return bh
}
