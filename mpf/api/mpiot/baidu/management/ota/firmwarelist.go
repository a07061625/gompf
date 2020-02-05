/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 23:49
 */
package ota

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取固件包列表
type firmwareList struct {
    mpiot.BaseBaiDu
}

func (fl *firmwareList) SetOrderBy(orderBy string) {
    if (orderBy == "asc") || (orderBy == "desc") {
        fl.ReqData["orderBy"] = orderBy
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序方式不合法", nil))
    }
}

func (fl *firmwareList) SetPageNo(pageNo int) {
    if pageNo > 0 {
        fl.ReqData["pageNo"] = strconv.Itoa(pageNo)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "页码不合法", nil))
    }
}

func (fl *firmwareList) SetPageSize(pageSize int) {
    if (pageSize > 0) && (pageSize <= 100) {
        fl.ReqData["pageSize"] = strconv.Itoa(pageSize)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "每页个数不合法", nil))
    }
}

func (fl *firmwareList) SetFirmwareId(firmwareId string) {
    if len(firmwareId) > 0 {
        fl.ReqData["firmwareId"] = firmwareId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "固件包ID不合法", nil))
    }
}

func (fl *firmwareList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    fl.ReqUrl = fl.GetServiceUrl() + "?" + mpf.HttpCreateParams(fl.ReqData, "none", 1)

    return fl.GetRequest()
}

func NewFirmwareList() *firmwareList {
    fl := &firmwareList{mpiot.NewBaseBaiDu()}
    fl.ServiceUri = "/v3/iot/management/ota/firmware"
    fl.ReqData["order"] = "createdAt"
    fl.ReqData["orderBy"] = "desc"
    fl.ReqData["pageNo"] = "1"
    fl.ReqData["pageSize"] = "10"
    return fl
}
