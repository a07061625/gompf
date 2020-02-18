/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 9:34
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

// 设置存储桶的生命周期
type lifeCyclePut struct {
    qcloud.BaseCos
    lifeConfig []map[string]interface{} // 生命周期配置
}

func (lcp *lifeCyclePut) SetLifeConfig(lifeConfig []map[string]interface{}) {
    if len(lifeConfig) > 0 {
        lcp.lifeConfig = lifeConfig
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "生命周期配置不合法", nil))
    }
}

func (lcp *lifeCyclePut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(lcp.lifeConfig) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "生命周期配置不能为空", nil))
    }

    xmlStr := ""
    for _, v := range lcp.lifeConfig {
        xmlData := mxj.Map(v)
        eXml, _ := xmlData.Xml("Rule")
        xmlStr += string(eXml)
    }

    reqBody := project.DataPrefixXML + "<LifecycleConfiguration>" + xmlStr + "</LifecycleConfiguration>"
    encodeStr := base64.StdEncoding.EncodeToString([]byte(reqBody))
    lcp.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    lcp.ReqURI = "http://" + lcp.ReqHeader["Host"] + lcp.ReqUri + "?lifecycle"
    client, req := lcp.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewLifeCyclePut() *lifeCyclePut {
    lcp := &lifeCyclePut{qcloud.NewCos(), make([]map[string]interface{}, 0)}
    lcp.ReqMethod = fasthttp.MethodPut
    lcp.SetParamData("lifecycle", "")
    return lcp
}
