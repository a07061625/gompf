/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 14:14
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除清单任务
type inventoryDelete struct {
    qcloud.BaseCos
    taskId string // 任务ID
}

func (id *inventoryDelete) SetTaskId(taskId string) {
    if len(taskId) > 0 {
        id.SetParamData("id", taskId)
        id.taskId = taskId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不合法", nil))
    }
}

func (id *inventoryDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(id.taskId) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不能为空", nil))
    }

    delete(id.ReqData, "inventory")
    id.ReqUrl = "http://" + id.ReqHeader["Host"] + id.ReqUri + "?inventory&" + mpf.HTTPCreateParams(id.ReqData, "none", 1)
    return id.GetRequest()
}

func NewInventoryDelete() *inventoryDelete {
    id := &inventoryDelete{qcloud.NewCos(), ""}
    id.ReqMethod = fasthttp.MethodDelete
    id.SetParamData("inventory", "")
    return id
}
