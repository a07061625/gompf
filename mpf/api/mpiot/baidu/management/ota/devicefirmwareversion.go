/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 16:51
 */
package ota

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询设备使用固件包版本
type deviceFirmwareVersion struct {
    mpiot.BaseBaiDu
    schemaId string // 物模型ID
}

func (dfv *deviceFirmwareVersion) SetSchemaId(schemaId string) {
    if len(schemaId) > 0 {
        dfv.schemaId = schemaId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "物模型ID不合法", nil))
    }
}

func (dfv *deviceFirmwareVersion) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dfv.schemaId) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "物模型ID不能为空", nil))
    }
    dfv.ExtendData["schemaId"] = dfv.schemaId

    dfv.ReqUrl = dfv.GetServiceUrl()

    reqBody := mpf.JSONMarshal(dfv.ExtendData)
    client, req := dfv.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDeviceFirmwareVersion() *deviceFirmwareVersion {
    dfv := &deviceFirmwareVersion{mpiot.NewBaseBaiDu(), ""}
    dfv.ServiceUri = "/v3/iot/management/ota/device-firmware-version-query"
    dfv.ReqContentType = project.HTTPContentTypeJSON
    dfv.ReqMethod = fasthttp.MethodPost
    return dfv
}
