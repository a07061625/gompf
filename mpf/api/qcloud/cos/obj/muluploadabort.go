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
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 舍弃并删除一个上传分块
type mulUploadAbort struct {
    qcloud.BaseCos
    objectKey string // 对象名称
    uploadId  string // 上传ID
}

func (mua *mulUploadAbort) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        mua.ReqUri = "/" + objectKey
        mua.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (mua *mulUploadAbort) SetUploadId(uploadId string) {
    if len(uploadId) > 0 {
        mua.SetParamData("uploadId", uploadId)
        mua.uploadId = uploadId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传ID不合法", nil))
    }
}

func (mua *mulUploadAbort) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(mua.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }
    if len(mua.uploadId) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传ID不能为空", nil))
    }

    mua.ReqUrl = "http://" + mua.ReqHeader["Host"] + mua.ReqUri + "?" + mpf.HttpCreateParams(mua.ReqData, "none", 1)
    return mua.GetRequest()
}

func NewMulUploadAbort() *mulUploadAbort {
    mua := &mulUploadAbort{qcloud.NewCos(), "", ""}
    mua.ReqMethod = fasthttp.MethodDelete
    return mua
}
