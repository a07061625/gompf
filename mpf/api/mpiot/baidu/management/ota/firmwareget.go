/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 16:51
 */
package ota

import (
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取固件包详情
type firmwareGet struct {
    mpiot.BaseBaiDu
    firmwareId string // 固件包ID
}

func (fg *firmwareGet) SetFirmwareId(firmwareId string) {
    if len(firmwareId) > 0 {
        fg.firmwareId = firmwareId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "固件包ID不合法", nil))
    }
}

func (fg *firmwareGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(fg.firmwareId) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "固件包ID不能为空", nil))
    }
    fg.ServiceUri = "/v3/iot/management/ota/firmware/" + fg.firmwareId

    fg.ReqURI = fg.GetServiceUrl()

    return fg.GetRequest()
}

func NewFirmwareGet() *firmwareGet {
    fg := &firmwareGet{mpiot.NewBaseBaiDu(), ""}
    return fg
}
