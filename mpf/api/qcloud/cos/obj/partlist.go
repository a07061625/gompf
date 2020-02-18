/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 21:50
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

// 获取特定上传事件中已上传的块
type partList struct {
    qcloud.BaseCos
    objectKey string // 对象名称
    uploadId  string // 上传ID
}

func (pl *partList) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        pl.ReqUri = "/" + objectKey
        pl.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (pl *partList) SetUploadId(uploadId string) {
    if len(uploadId) > 0 {
        pl.SetParamData("uploadId", uploadId)
        pl.uploadId = uploadId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传ID不合法", nil))
    }
}

func (pl *partList) SetEncodingType(encodingType string) {
    if len(encodingType) > 0 {
        pl.SetParamData("encoding-type", encodingType)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "编码方式不合法", nil))
    }
}

func (pl *partList) SetMaxPart(maxPart int) {
    if (maxPart > 0) && (maxPart <= 1000) {
        pl.SetParamData("max-parts", strconv.Itoa(maxPart))
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "最大条目数不合法", nil))
    }
}

func (pl *partList) SetNumberMarker(numberMarker string) {
    if len(numberMarker) > 0 {
        pl.SetParamData("part-number-marker", numberMarker)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "起始索引不合法", nil))
    }
}

func (pl *partList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pl.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }
    if len(pl.uploadId) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传ID不能为空", nil))
    }
    pl.ReqURI = "http://" + pl.ReqHeader["Host"] + pl.ReqUri + "?" + mpf.HTTPCreateParams(pl.ReqData, "none", 1)
    return pl.GetRequest()
}

func NewPartList() *partList {
    pl := &partList{qcloud.NewCos(), "", ""}
    return pl
}
