/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 9:12
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 创建存储桶
type bucketPut struct {
    qcloud.BaseCos
}

func (bp *bucketPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    bp.ReqUrl = "http://" + bp.ReqHeader["Host"] + bp.ReqUri
    return bp.GetRequest()
}

func NewBucketPut() *bucketPut {
    bp := &bucketPut{qcloud.NewCos()}
    bp.ReqMethod = fasthttp.MethodPut
    return bp
}
