/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 16:47
 */
package obj

import (
    "encoding/base64"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/clbanning/mxj"
    "github.com/valyala/fasthttp"
)

// 设置对象的访问权限
type aclPut struct {
    qcloud.BaseCos
    objectKey  string                 // 对象名称
    policyInfo map[string]interface{} // 权限信息
}

func (ap *aclPut) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        ap.ReqUri = "/" + objectKey
        ap.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (ap *aclPut) SetPolicyInfo(policyInfo map[string]interface{}) {
    if len(policyInfo) > 0 {
        ap.policyInfo = policyInfo
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "权限信息不合法", nil))
    }
}

func (ap *aclPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ap.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }
    if len(ap.policyInfo) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "权限信息不能为空", nil))
    }
    xmlData := mxj.Map(ap.policyInfo)
    xmlInfo, _ := xmlData.Xml("AccessControlPolicy")
    reqBody := project.DataPrefixXML + string(xmlInfo) + ""
    encodeStr := base64.StdEncoding.EncodeToString([]byte(reqBody))
    ap.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    ap.ReqURI = "http://" + ap.ReqHeader["Host"] + ap.ReqUri + "?acl"
    client, req := ap.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewAclPut() *aclPut {
    ap := &aclPut{qcloud.NewCos(), "", make(map[string]interface{})}
    ap.ReqMethod = fasthttp.MethodPut
    ap.SetParamData("acl", "")
    return ap
}
