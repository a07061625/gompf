/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 9:19
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

// 设置存储桶的跨域访问配置信息
type corsPut struct {
    qcloud.BaseCos
    corsConfig []map[string]interface{} // 跨域访问配置
}

func (cp *corsPut) SetCorsConfig(corsConfig []map[string]interface{}) {
    if (len(corsConfig) > 0) && (len(corsConfig) <= 100) {
        cp.corsConfig = corsConfig
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "跨域访问配置不合法", nil))
    }
}

func (cp *corsPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cp.corsConfig) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "跨域访问配置不能为空", nil))
    }

    xmlStr := ""
    for _, v := range cp.corsConfig {
        xmlData := mxj.Map(v)
        eXml, _ := xmlData.Xml("CORSRule")
        xmlStr += string(eXml)
    }

    reqBody := project.DataPrefixXml + "<CORSConfiguration>" + xmlStr + "</CORSConfiguration>"
    encodeStr := base64.StdEncoding.EncodeToString([]byte(reqBody))
    cp.ReqUrl = "http://" + cp.ReqHeader["Host"] + cp.ReqUri + "?cors"
    cp.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    client, req := cp.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewCorsPut() *corsPut {
    cp := &corsPut{qcloud.NewCos(), make([]map[string]interface{}, 0)}
    cp.SetParamData("cors", "")
    cp.ReqMethod = fasthttp.MethodPut
    return cp
}
