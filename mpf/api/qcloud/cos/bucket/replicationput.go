/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 14:39
 */
package bucket

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

// 配置跨地域复制规则
type replicationPut struct {
    qcloud.BaseCos
    ruleInfo map[string]interface{} // 规则信息
}

func (rp *replicationPut) SetRuleInfo(ruleInfo map[string]interface{}) {
    if len(ruleInfo) > 0 {
        rp.ruleInfo = ruleInfo
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "规则信息不合法", nil))
    }
}

func (rp *replicationPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rp.ruleInfo) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "规则信息不能为空", nil))
    }

    xmlData := mxj.Map(rp.ruleInfo)
    xmlStr, _ := xmlData.Xml("ReplicationConfiguration")
    reqBody := project.DataPrefixXML + string(xmlStr)
    encodeStr := base64.StdEncoding.EncodeToString([]byte(reqBody))
    rp.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    rp.ReqUrl = "http://" + rp.ReqHeader["Host"] + rp.ReqUri + "?replication"
    client, req := rp.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewReplicationPut() *replicationPut {
    rp := &replicationPut{qcloud.NewCos(), make(map[string]interface{})}
    rp.ReqMethod = fasthttp.MethodPut
    rp.SetParamData("replication", "")
    return rp
}
