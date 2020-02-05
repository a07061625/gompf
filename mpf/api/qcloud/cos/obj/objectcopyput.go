/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 22:03
 */
package obj

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 复制文件
type objectCopyPut struct {
    qcloud.BaseCos
    objectKey string // 对象名称
}

func (ocp *objectCopyPut) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        ocp.objectKey = "/" + objectKey
        ocp.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (ocp *objectCopyPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ocp.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }
    _, ok := ocp.ReqHeader["x-cos-copy-source"]
    if !ok {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "源对象URL不能为空", nil))
    }

    ocp.ReqUrl = "http://" + ocp.ReqHeader["Host"] + ocp.ReqUri
    return ocp.GetRequest()
}

func NewObjectCopyPut() *objectCopyPut {
    ocp := &objectCopyPut{qcloud.NewCos(), ""}
    ocp.ReqMethod = fasthttp.MethodPut
    return ocp
}
