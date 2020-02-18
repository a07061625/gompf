/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 10:11
 */
package bucket

import (
    "encoding/base64"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 设置存储桶的权限策略
type policyPut struct {
    qcloud.BaseCos
    policyConfig map[string]interface{} // 权限策略配置
}

func (pp *policyPut) SetPolicyConfig(policyConfig map[string]interface{}) {
    if len(policyConfig) > 0 {
        pp.policyConfig = policyConfig
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "权限策略配置不合法", nil))
    }
}

func (pp *policyPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pp.policyConfig) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "权限策略配置不能为空", nil))
    }

    reqBody := mpf.JsonMarshal(pp.policyConfig)
    encodeStr := base64.StdEncoding.EncodeToString([]byte(reqBody))
    pp.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    pp.ReqUrl = "http://" + pp.ReqHeader["Host"] + pp.ReqUri + "?policy"
    client, req := pp.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewPolicyPut() *policyPut {
    pp := &policyPut{qcloud.NewCos(), make(map[string]interface{})}
    pp.ReqMethod = fasthttp.MethodPut
    pp.ReqContentType = project.HTTPContentTypeJSON
    pp.SetParamData("policy", "")
    return pp
}
