/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 23:18
 */
package obj

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 下载单个对象
type objectGet struct {
    qcloud.BaseCos
    objectKey string // 对象名称
}

func (og *objectGet) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        og.objectKey = "/" + objectKey
        og.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (og *objectGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(og.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }

    og.ReqURI = "http://" + og.ReqHeader["Host"] + og.ReqUri
    return og.GetRequest()
}

func NewObjectGet() *objectGet {
    og := &objectGet{qcloud.NewCos(), ""}
    return og
}
