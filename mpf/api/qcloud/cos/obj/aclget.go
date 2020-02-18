/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 16:42
 */
package obj

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取对象的访问权限
type aclGet struct {
    qcloud.BaseCos
    objectKey string // 对象名称
}

func (ag *aclGet) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        ag.ReqUri = "/" + objectKey
        ag.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (ag *aclGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ag.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }

    ag.ReqURI = "http://" + ag.ReqHeader["Host"] + ag.ReqUri + "?acl"
    return ag.GetRequest()
}

func NewAclGet() *aclGet {
    ag := &aclGet{qcloud.NewCos(), ""}
    ag.SetParamData("acl", "")
    return ag
}
