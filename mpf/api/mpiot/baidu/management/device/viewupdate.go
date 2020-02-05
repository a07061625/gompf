/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:29
 */
package device

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 修改设备View信息
type viewUpdate struct {
    mpiot.BaseBaiDu
    deviceName     string // 设备名称
    profileVersion int    // 版本号
    reported       map[string]interface{}
    desired        map[string]interface{}
}

func (vu *viewUpdate) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        vu.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (vu *viewUpdate) SetProfileVersion(profileVersion int) {
    if profileVersion > 0 {
        vu.profileVersion = profileVersion
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "版本号不合法", nil))
    }
}

func (vu *viewUpdate) SetReported(reported map[string]interface{}) {
    vu.reported = reported
}

func (vu *viewUpdate) SetDesired(desired map[string]interface{}) {
    vu.desired = desired
}

func (vu *viewUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(vu.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    if vu.profileVersion <= 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "版本号不能为空", nil))
    }
    if (len(vu.reported) == 0) && (len(vu.desired) == 0) {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "reported信息和desired信息不能都为空", nil))
    }
    vu.ServiceUri = "/v3/iot/management/deviceView/" + vu.deviceName
    vu.ExtendData["profileVersion"] = vu.profileVersion
    if len(vu.reported) > 0 {
        vu.ExtendData["reported"] = vu.reported
    }
    if len(vu.desired) > 0 {
        vu.ExtendData["desired"] = vu.desired
    }

    vu.ReqUrl = vu.GetServiceUrl() + "?updateView"

    reqBody := mpf.JsonMarshal(vu.ExtendData)
    client, req := vu.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewViewUpdate() *viewUpdate {
    vu := &viewUpdate{mpiot.NewBaseBaiDu(), "", 0, make(map[string]interface{}), make(map[string]interface{})}
    vu.ReqData["updateView"] = ""
    vu.ReqContentType = project.HttpContentTypeJson
    vu.ReqMethod = fasthttp.MethodPut
    return vu
}
