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

// 获取清单任务列表
type inventoryList struct {
    qcloud.BaseCos
}

func (il *inventoryList) SetToken(token string) {
    if len(token) > 0 {
        il.SetParamData("continuation-token", token)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "令牌不合法", nil))
    }
}

func (il *inventoryList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    delete(il.ReqData, "inventory")
    il.ReqUrl = "http://" + il.ReqHeader["Host"] + il.ReqUri + "?inventory"
    if len(il.ReqData) > 0 {
        il.ReqUrl += "&" + mpf.HTTPCreateParams(il.ReqData, "none", 1)
    }
    return il.GetRequest()
}

func NewInventoryList() *inventoryList {
    il := &inventoryList{qcloud.NewCos()}
    il.SetParamData("inventory", "")
    return il
}
