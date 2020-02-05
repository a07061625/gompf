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

// 获取清单任务
type inventoryGet struct {
    qcloud.BaseCos
    taskId string // 任务ID
}

func (ig *inventoryGet) SetTaskId(taskId string) {
    if len(taskId) > 0 {
        ig.SetParamData("id", taskId)
        ig.taskId = taskId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不合法", nil))
    }
}

func (ig *inventoryGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ig.taskId) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不能为空", nil))
    }

    delete(ig.ReqData, "inventory")
    ig.ReqUrl = "http://" + ig.ReqHeader["Host"] + ig.ReqUri + "?inventory&" + mpf.HttpCreateParams(ig.ReqData, "none", 1)
    return ig.GetRequest()
}

func NewInventoryGet() *inventoryGet {
    ig := &inventoryGet{qcloud.NewCos(), ""}
    ig.SetParamData("inventory", "")
    return ig
}
