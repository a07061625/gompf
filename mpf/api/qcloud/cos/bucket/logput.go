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

// 配置日志记录
type logPut struct {
    qcloud.BaseCos
    logInfo map[string]interface{} // 日志信息
}

func (lp *logPut) SetLogInfo(logInfo map[string]interface{}) {
    if len(logInfo) > 0 {
        lp.logInfo = logInfo
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "日志信息不合法", nil))
    }
}

func (lp *logPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(lp.logInfo) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "日志信息不能为空", nil))
    }

    xmlData := mxj.Map(lp.logInfo)
    xmlStr, _ := xmlData.Xml("BucketLoggingStatus")
    reqBody := project.DataPrefixXML + string(xmlStr)
    encodeStr := base64.StdEncoding.EncodeToString([]byte(reqBody))
    lp.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    lp.ReqUrl = "http://" + lp.ReqHeader["Host"] + lp.ReqUri + "?logging"
    client, req := lp.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewLogPut() *logPut {
    lp := &logPut{qcloud.NewCos(), make(map[string]interface{})}
    lp.ReqMethod = fasthttp.MethodPut
    lp.SetParamData("logging", "")
    return lp
}
