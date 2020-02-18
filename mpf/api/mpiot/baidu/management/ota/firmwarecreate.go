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

// 创建固件包
type firmwareCreate struct {
    mpiot.BaseBaiDu
    schemaId    string // 物模型ID
    description string // 描述
    version     string // 版本号
    fileId      string // 固件包文件ID
}

func (fc *firmwareCreate) SetSchemaId(schemaId string) {
    if len(schemaId) > 0 {
        fc.schemaId = schemaId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "物模型ID不合法", nil))
    }
}

func (fc *firmwareCreate) SetDescription(description string) {
    if len(description) > 0 {
        fc.description = description
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "描述不合法", nil))
    }
}

func (fc *firmwareCreate) SetVersion(version string) {
    if len(version) > 0 {
        fc.version = version
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "版本号不合法", nil))
    }
}

func (fc *firmwareCreate) SetFileId(fileId string) {
    if len(fileId) > 0 {
        fc.fileId = fileId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "固件包文件ID不合法", nil))
    }
}

func (fc *firmwareCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(fc.schemaId) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "物模型ID不能为空", nil))
    }
    if len(fc.description) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "描述不能为空", nil))
    }
    if len(fc.version) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "版本号不能为空", nil))
    }
    if len(fc.fileId) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "固件包文件ID不能为空", nil))
    }
    fc.ExtendData["schemaId"] = fc.schemaId
    fc.ExtendData["description"] = fc.description
    fc.ExtendData["version"] = fc.version
    fc.ExtendData["fileId"] = fc.fileId

    fc.ReqUrl = fc.GetServiceUrl()

    reqBody := mpf.JSONMarshal(fc.ExtendData)
    client, req := fc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewFirmwareCreate() *firmwareCreate {
    fc := &firmwareCreate{mpiot.NewBaseBaiDu(), "", "", "", ""}
    fc.ServiceUri = "/v3/iot/management/ota/firmware"
    fc.ReqContentType = project.HTTPContentTypeJSON
    fc.ReqMethod = fasthttp.MethodPost
    return fc
}
