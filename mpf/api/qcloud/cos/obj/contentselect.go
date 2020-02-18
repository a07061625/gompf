/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 0:02
 */
package obj

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/clbanning/mxj"
    "github.com/valyala/fasthttp"
)

// 检索指定对象(CSV,JSON)的内容
type contentSelect struct {
    qcloud.BaseCos
    objectKey  string                 // 对象名称
    selectInfo map[string]interface{} // 检索信息
}

func (cs *contentSelect) SetObjectKey(objectKey string) {
    if len(objectKey) > 0 {
        cs.ReqUri = "/" + objectKey
        cs.objectKey = objectKey
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不合法", nil))
    }
}

func (cs *contentSelect) SetSelectInfo(selectInfo map[string]interface{}) {
    if len(selectInfo) > 0 {
        cs.selectInfo = selectInfo
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "检索信息不合法", nil))
    }
}

func (cs *contentSelect) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cs.objectKey) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象名称不能为空", nil))
    }
    if len(cs.selectInfo) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "检索信息不能为空", nil))
    }

    xmlData := mxj.Map(cs.selectInfo)
    xmlStr, _ := xmlData.Xml("SelectRequest")
    reqBody := project.DataPrefixXML + string(xmlStr)
    delete(cs.ReqData, "select")
    cs.ReqUrl = "http://" + cs.ReqHeader["Host"] + cs.ReqUri + "?select&" + mpf.HTTPCreateParams(cs.ReqData, "none", 1)
    client, req := cs.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewContentSelect() *contentSelect {
    cs := &contentSelect{qcloud.NewCos(), "", make(map[string]interface{})}
    cs.ReqMethod = fasthttp.MethodPost
    cs.SetParamData("select", "")
    cs.SetParamData("select-type", "2")
    return cs
}
