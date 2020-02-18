/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 17:52
 */
package obj

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 获取进行中的分块上传列表
type mulUploadList struct {
    qcloud.BaseCos
}

func (mul *mulUploadList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    delete(mul.ReqData, "uploads")
    mul.ReqURI = "http://" + mul.ReqHeader["Host"] + mul.ReqUri + "?uploads&" + mpf.HTTPCreateParams(mul.ReqData, "none", 1)
    return mul.GetRequest()
}

func NewMulUploadList() *mulUploadList {
    mul := &mulUploadList{qcloud.NewCos()}
    mul.SetParamData("uploads", "")
    mul.SetParamData("max-uploads", "100")
    mul.SetParamData("encoding-type", "url")
    return mul
}
