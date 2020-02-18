/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 23:51
 */
package obj

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/clbanning/mxj"
    "github.com/valyala/fasthttp"
)

// 恢复archive类型的对象
type objectRestore struct {
    qcloud.BaseCos
    objectKey   string                 // 对象名称
    restoreInfo map[string]interface{} // 恢复信息
}

func (or *objectRestore) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        or.ReqUri = "/" + objectKey
        or.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (or *objectRestore) SetRestoreInfo(restoreInfo map[string]interface{}) {
    if len(restoreInfo) > 0 {
        or.restoreInfo = restoreInfo
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "恢复信息不合法", nil))
    }
}

func (or *objectRestore) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(or.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }
    if len(or.restoreInfo) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "恢复信息不能为空", nil))
    }

    xmlData := mxj.Map(or.restoreInfo)
    xmlStr, _ := xmlData.Xml("RestoreRequest")
    reqBody := project.DataPrefixXML + string(xmlStr)
    or.ReqUrl = "http://" + or.ReqHeader["Host"] + or.ReqUri + "?restore"
    client, req := or.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewObjectRestore() *objectRestore {
    or := &objectRestore{qcloud.NewCos(), "", make(map[string]interface{})}
    or.ReqMethod = fasthttp.MethodPost
    or.SetParamData("restore", "")
    return or
}
