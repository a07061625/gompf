/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 14:14
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

// 创建清单任务
type inventoryPut struct {
    qcloud.BaseCos
    taskId   string                 // 任务ID
    taskInfo map[string]interface{} // 任务信息
}

func (ip *inventoryPut) SetTaskId(taskId string) {
    if len(taskId) > 0 {
        ip.SetParamData("id", taskId)
        ip.taskId = taskId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不合法", nil))
    }
}

func (ip *inventoryPut) SetTaskInfo(taskInfo map[string]interface{}) {
    if len(taskInfo) > 0 {
        ip.taskInfo = taskInfo
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务信息不合法", nil))
    }
}

func (ip *inventoryPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ip.taskId) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不能为空", nil))
    }
    if len(ip.taskInfo) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务信息不能为空", nil))
    }

    xmlData := mxj.Map(ip.taskInfo)
    xmlStr, _ := xmlData.Xml("InventoryConfiguration")
    reqBody := project.DataPrefixXML + string(xmlStr)
    encodeStr := base64.StdEncoding.EncodeToString([]byte(reqBody))
    ip.SetHeaderData("Content-MD5", mpf.HashMd5(encodeStr, ""))
    delete(ip.ReqData, "inventory")
    ip.ReqURI = "http://" + ip.ReqHeader["Host"] + ip.ReqUri + "?inventory&" + mpf.HTTPCreateParams(ip.ReqData, "none", 1)
    client, req := ip.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewInventoryPut() *inventoryPut {
    ip := &inventoryPut{qcloud.NewCos(), "", make(map[string]interface{})}
    ip.ReqMethod = fasthttp.MethodPut
    ip.SetParamData("inventory", "")
    return ip
}
