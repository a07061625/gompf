/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 23:36
 */
package obj

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

// 批量删除多个对象
type objectMulDelete struct {
    qcloud.BaseCos
    modeType   string                   // 模式类型
    objectList []map[string]interface{} // 对象列表
}

func (omd *objectMulDelete) SetModeType(modeType string) {
    if (modeType == "true") || (modeType == "false") {
        omd.modeType = modeType
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "模式类型不合法", nil))
    }
}

func (omd *objectMulDelete) SetObjectList(objectList []map[string]interface{}) {
    if len(objectList) > 0 {
        omd.objectList = objectList
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象列表不合法", nil))
    }
}

func (omd *objectMulDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(omd.modeType) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "模式类型不能为空", nil))
    }
    if len(omd.objectList) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "对象列表不能为空", nil))
    }

    xmlStr := ""
    for _, v := range omd.objectList {
        xmlData := mxj.Map(v)
        eXml, _ := xmlData.Xml("Object")
        xmlStr += string(eXml)
    }
    reqBody := project.DataPrefixXML + "<Delete><Quiet>" + omd.modeType + "</Quiet>" + xmlStr + "</Delete>"

    encodeStr := base64.StdEncoding.EncodeToString([]byte(reqBody))
    omd.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    omd.ReqURI = "http://" + omd.ReqHeader["Host"] + omd.ReqUri + "?delete"
    client, req := omd.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewObjectMulDelete() *objectMulDelete {
    omd := &objectMulDelete{qcloud.NewCos(), "", make([]map[string]interface{}, 0)}
    omd.ReqMethod = fasthttp.MethodPost
    omd.SetParamData("delete", "")
    return omd
}
