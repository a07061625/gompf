/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 9:57
 */
package bucket

import (
    "encoding/base64"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/clbanning/mxj"
    "github.com/valyala/fasthttp"
)

// 设置存储桶Referer的白名单或者黑名单
type refererPut struct {
    qcloud.BaseCos
    refererConfig map[string]interface{} // 配置信息
}

func (rp *refererPut) SetRefererConfig(refererConfig map[string]interface{}) {
    if len(refererConfig) > 0 {
        rp.refererConfig = refererConfig
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "配置信息不合法", nil))
    }
}

func (rp *refererPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rp.refererConfig) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "配置信息不能为空", nil))
    }

    xmlData := mxj.Map(rp.refererConfig)
    xmlStr, _ := xmlData.Xml("RefererConfiguration")

    reqBody := project.DataPrefixXML + string(xmlStr)
    encodeStr := base64.StdEncoding.EncodeToString([]byte(reqBody))
    rp.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    rp.SetHeaderData("Content-Length", strconv.Itoa(len(reqBody)))
    rp.ReqURI = "http://" + rp.ReqHeader["Host"] + rp.ReqUri + "?referer"
    client, req := rp.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewRefererPut() *refererPut {
    rp := &refererPut{qcloud.NewCos(), make(map[string]interface{})}
    rp.ReqMethod = fasthttp.MethodPut
    rp.SetParamData("referer", "")
    return rp
}
