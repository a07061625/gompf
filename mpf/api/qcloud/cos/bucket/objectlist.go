/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 9:16
 */
package bucket

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取存储桶的对象列表
type objectList struct {
    qcloud.BaseCos
}

func (ol *objectList) SetPrefix(prefix string) {
    if len(prefix) > 0 {
        ol.SetParamData("prefix", prefix)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "前缀不合法", nil))
    }
}

func (ol *objectList) SetDelimiter(delimiter string) {
    if len(delimiter) > 0 {
        ol.SetParamData("delimiter", delimiter)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "分隔符不合法", nil))
    }
}

func (ol *objectList) SetMaxNum(maxNum int) {
    if (maxNum > 0) && (maxNum <= 1000) {
        ol.SetParamData("max-keys", strconv.Itoa(maxNum))
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "最大条目数不合法", nil))
    }
}

func (ol *objectList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    ol.ReqURI = "http://" + ol.ReqHeader["Host"] + ol.ReqUri + "?" + mpf.HTTPCreateParams(ol.ReqData, "none", 1)
    return ol.GetRequest()
}

func NewObjectList() *objectList {
    ol := &objectList{qcloud.NewCos()}
    ol.SetParamData("max-keys", "100")
    ol.SetParamData("encoding-type", "url")
    return ol
}
