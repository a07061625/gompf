/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 15:07
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/clbanning/mxj"
    "github.com/valyala/fasthttp"
)

// 设置默认加密配置
type encryptionPut struct {
    qcloud.BaseCos
    configInfo map[string]interface{} // 配置信息
}

func (ep *encryptionPut) SetConfigInfo(configInfo map[string]interface{}) {
    if len(configInfo) > 0 {
        ep.configInfo = configInfo
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "配置信息不合法", nil))
    }
}

func (ep *encryptionPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ep.configInfo) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "配置信息不能为空", nil))
    }

    xmlData := mxj.Map(ep.configInfo)
    xmlStr, _ := xmlData.Xml("ServerSideEncryptionConfiguration")
    reqBody := project.DataPrefixXML + string(xmlStr)
    ep.ReqUrl = "http://" + ep.ReqHeader["Host"] + ep.ReqUri + "?encryption"
    client, req := ep.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewEncryptionPut() *encryptionPut {
    ep := &encryptionPut{qcloud.NewCos(), make(map[string]interface{})}
    ep.ReqMethod = fasthttp.MethodPut
    ep.SetParamData("encryption", "")
    return ep
}
