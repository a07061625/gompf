/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 17:21
 */
package obj

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 初始化分片上传
type mulUploadInitiate struct {
    qcloud.BaseCos
    objectKey string // 对象名称
}

func (mui *mulUploadInitiate) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        mui.ReqUri = "/" + objectKey
        mui.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (mui *mulUploadInitiate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(mui.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }

    mui.ReqURI = "http://" + mui.ReqHeader["Host"] + mui.ReqUri + "?uploads"
    return mui.GetRequest()
}

func NewMulUploadInitiate() *mulUploadInitiate {
    mui := &mulUploadInitiate{qcloud.NewCos(), ""}
    mui.ReqMethod = fasthttp.MethodPost
    mui.SetParamData("uploads", "")
    return mui
}
