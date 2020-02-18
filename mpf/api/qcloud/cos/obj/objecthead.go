/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 23:26
 */
package obj

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取对象的meta信息数据
type objectHead struct {
    qcloud.BaseCos
    objectKey string // 对象名称
}

func (oh *objectHead) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        oh.objectKey = "/" + objectKey
        oh.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (oh *objectHead) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(oh.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }

    oh.ReqURI = "http://" + oh.ReqHeader["Host"] + oh.ReqUri
    return oh.GetRequest()
}

func NewObjectHead() *objectHead {
    oh := &objectHead{qcloud.NewCos(), ""}
    oh.ReqMethod = fasthttp.MethodHead
    return oh
}
