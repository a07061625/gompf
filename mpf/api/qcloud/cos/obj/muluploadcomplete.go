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
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/clbanning/mxj"
    "github.com/valyala/fasthttp"
)

// 完成分块上传
type mulUploadComplete struct {
    qcloud.BaseCos
    objectKey string                   // 对象名称
    uploadId  string                   // 上传ID
    partList  []map[string]interface{} // 内容列表
}

func (muc *mulUploadComplete) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        muc.ReqUri = "/" + objectKey
        muc.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (muc *mulUploadComplete) SetUploadId(uploadId string) {
    if len(uploadId) > 0 {
        muc.SetParamData("uploadId", uploadId)
        muc.uploadId = uploadId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传ID不合法", nil))
    }
}

func (muc *mulUploadComplete) SetPartList(partList []map[string]interface{}) {
    if len(partList) > 0 {
        muc.partList = partList
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "内容列表不合法", nil))
    }
}

func (muc *mulUploadComplete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(muc.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }
    if len(muc.uploadId) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "上传ID不能为空", nil))
    }
    if len(muc.partList) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "内容列表不能为空", nil))
    }

    xmlStr := ""
    for _, v := range muc.partList {
        xmlData := mxj.Map(v)
        eXml, _ := xmlData.Xml("Part")
        xmlStr += string(eXml)
    }

    reqBody := project.DataPrefixXML + "<CompleteMultipartUpload>" + xmlStr + "</CompleteMultipartUpload>"
    muc.ReqURI = "http://" + muc.ReqHeader["Host"] + muc.ReqUri + "?" + mpf.HTTPCreateParams(muc.ReqData, "none", 1)
    client, req := muc.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewMulUploadComplete() *mulUploadComplete {
    muc := &mulUploadComplete{qcloud.NewCos(), "", "", make([]map[string]interface{}, 0)}
    muc.ReqMethod = fasthttp.MethodPost
    return muc
}
