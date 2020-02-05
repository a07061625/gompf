/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 10:56
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

// 设置存储桶标签
type taggingPut struct {
    qcloud.BaseCos
    tagList []map[string]interface{} // 标签列表
}

func (tp *taggingPut) SetPartList(tagList []map[string]interface{}) {
    if (len(tagList) > 0) && (len(tagList) <= 10) {
        tp.tagList = tagList
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "标签列表不合法", nil))
    }
}

func (tp *taggingPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tp.tagList) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "标签列表不能为空", nil))
    }

    xmlStr := ""
    for _, v := range tp.tagList {
        xmlData := mxj.Map(v)
        eXml, _ := xmlData.Xml("Tag")
        xmlStr += string(eXml)
    }

    reqBody := project.DataPrefixXml + "<Tagging><TagSet>" + xmlStr + "</TagSet></Tagging>"
    encodeStr := base64.StdEncoding.EncodeToString([]byte(reqBody))
    tp.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    tp.ReqUrl = "http://" + tp.ReqHeader["Host"] + tp.ReqUri + "?tagging"
    client, req := tp.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewTaggingPut() *taggingPut {
    tp := &taggingPut{qcloud.NewCos(), make([]map[string]interface{}, 0)}
    tp.ReqMethod = fasthttp.MethodPut
    tp.SetParamData("tagging", "")
    return tp
}
