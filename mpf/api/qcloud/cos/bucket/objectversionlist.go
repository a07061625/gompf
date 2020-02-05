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

// 获取存储桶的对象及其历史版本信息象列表
type objectVersionList struct {
    qcloud.BaseCos
}

func (ovl *objectVersionList) SetPrefix(prefix string) {
    if len(prefix) > 0 {
        ovl.SetParamData("prefix", prefix)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "前缀不合法", nil))
    }
}

func (ovl *objectVersionList) SetDelimiter(delimiter string) {
    if len(delimiter) > 0 {
        ovl.SetParamData("delimiter", delimiter)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "分隔符不合法", nil))
    }
}

func (ovl *objectVersionList) SetMaxNum(maxNum int) {
    if (maxNum > 0) && (maxNum <= 1000) {
        ovl.SetParamData("max-keys", strconv.Itoa(maxNum))
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "最大条目数不合法", nil))
    }
}

func (ovl *objectVersionList) SetKeyMarker(keyMarker string) {
    if len(keyMarker) > 0 {
        ovl.SetParamData("key-marker", keyMarker)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "起始对象键标记不合法", nil))
    }
}

func (ovl *objectVersionList) SetVersionIdMarker(versionIdMarker string) {
    if len(versionIdMarker) > 0 {
        ovl.SetParamData("version-id-marker", versionIdMarker)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "起始版本ID标记不合法", nil))
    }
}

func (ovl *objectVersionList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    ovl.ReqUrl = "http://" + ovl.ReqHeader["Host"] + ovl.ReqUri + "?" + mpf.HttpCreateParams(ovl.ReqData, "none", 1)
    return ovl.GetRequest()
}

func NewObjectVersionList() *objectVersionList {
    ovl := &objectVersionList{qcloud.NewCos()}
    ovl.SetParamData("max-keys", "100")
    ovl.SetParamData("encoding-type", "url")
    return ovl
}
