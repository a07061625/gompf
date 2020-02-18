/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 21:27
 */
package obj

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 复制分块内容
type partCopyUpload struct {
    qcloud.BaseCos
    objectKey  string // 对象名称
    partNumber int    // 分块编号
    uploadId   string // 上传ID
}

func (pcu *partCopyUpload) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        pcu.ReqUri = "/" + objectKey
        pcu.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (pcu *partCopyUpload) SetPartNumber(partNumber int) {
    if partNumber > 0 {
        pcu.SetParamData("partNumber", strconv.Itoa(partNumber))
        pcu.partNumber = partNumber
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "分块编号不合法", nil))
    }
}

func (pcu *partCopyUpload) SetUploadId(uploadId string) {
    if len(uploadId) > 0 {
        pcu.SetParamData("uploadId", uploadId)
        pcu.uploadId = uploadId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传ID不合法", nil))
    }
}

func (pcu *partCopyUpload) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pcu.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }
    if pcu.partNumber <= 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "分块编号不能为空", nil))
    }
    if len(pcu.uploadId) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传ID不能为空", nil))
    }
    _, ok := pcu.ReqHeader["x-cos-copy-source"]
    if !ok {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "源对象URL不能为空", nil))
    }
    pcu.ReqURI = "http://" + pcu.ReqHeader["Host"] + pcu.ReqUri + "?" + mpf.HTTPCreateParams(pcu.ReqData, "none", 1)
    return pcu.GetRequest()
}

func NewPartCopyUpload() *partCopyUpload {
    pcu := &partCopyUpload{qcloud.NewCos(), "", 0, ""}
    pcu.ReqMethod = fasthttp.MethodPut
    return pcu
}
