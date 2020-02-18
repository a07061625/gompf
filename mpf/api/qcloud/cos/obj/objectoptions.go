/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 23:30
 */
package obj

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 对象跨域访问预请求
type objectOptions struct {
    qcloud.BaseCos
    objectKey string // 对象名称
}

func (oo *objectOptions) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        oo.objectKey = "/" + objectKey
        oo.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (oo *objectOptions) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(oo.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }
    _, ok := oo.ReqHeader["Origin"]
    if !ok {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "跨域请求域名不能为空", nil))
    }
    _, ok = oo.ReqHeader["Access-Control-Request-Method"]
    if !ok {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "跨域请求方法不能为空", nil))
    }

    oo.ReqURI = "http://" + oo.ReqHeader["Host"] + oo.ReqUri
    return oo.GetRequest()
}

func NewObjectOptions() *objectOptions {
    oo := &objectOptions{qcloud.NewCos(), ""}
    oo.ReqMethod = fasthttp.MethodOptions
    return oo
}
