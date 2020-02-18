/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 21:27
 */
package obj

import (
    "encoding/base64"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 分块上传
type partUpload struct {
    qcloud.BaseCos
    objectKey     string // 对象名称
    partNumber    int    // 分块编号
    uploadId      string // 上传ID
    uploadContent string // 上传内容
}

func (pu *partUpload) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        pu.ReqUri = "/" + objectKey
        pu.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (pu *partUpload) SetPartNumber(partNumber int) {
    if partNumber > 0 {
        pu.SetParamData("partNumber", strconv.Itoa(partNumber))
        pu.partNumber = partNumber
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "分块编号不合法", nil))
    }
}

func (pu *partUpload) SetUploadId(uploadId string) {
    if len(uploadId) > 0 {
        pu.SetParamData("uploadId", uploadId)
        pu.uploadId = uploadId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传ID不合法", nil))
    }
}

func (pu *partUpload) SetUploadContent(uploadContent string) {
    if len(uploadContent) > 0 {
        pu.uploadContent = uploadContent
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传内容不合法", nil))
    }
}

func (pu *partUpload) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pu.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }
    if pu.partNumber <= 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "分块编号不能为空", nil))
    }
    if len(pu.uploadId) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传ID不能为空", nil))
    }
    if len(pu.uploadContent) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传内容不能为空", nil))
    }
    encodeStr := base64.StdEncoding.EncodeToString([]byte(pu.uploadContent))
    pu.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    pu.SetHeaderData("Content-Length", strconv.Itoa(len(pu.uploadContent)))
    pu.ReqURI = "http://" + pu.ReqHeader["Host"] + pu.ReqUri + "?" + mpf.HTTPCreateParams(pu.ReqData, "none", 1)
    client, req := pu.GetRequest()
    req.SetBody([]byte(pu.uploadContent))

    return client, req
}

func NewPartUpload() *partUpload {
    pu := &partUpload{qcloud.NewCos(), "", 0, "", ""}
    pu.ReqMethod = fasthttp.MethodPut
    pu.SetHeaderData("Expect", "")
    return pu
}
