/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 23:11
 */
package obj

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除单个对象
type objectDelete struct {
    qcloud.BaseCos
    objectKey string // 对象名称
}

func (od *objectDelete) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        od.ReqUri = "/" + objectKey
        od.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (od *objectDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(od.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }

    od.ReqUrl = "http://" + od.ReqHeader["Host"] + od.ReqUri
    return od.GetRequest()
}

func NewObjectDelete() *objectDelete {
    od := &objectDelete{qcloud.NewCos(), ""}
    od.ReqMethod = fasthttp.MethodDelete
    return od
}
