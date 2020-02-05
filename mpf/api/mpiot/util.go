/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 15:15
 */
package mpiot

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
)

type utilIot struct {
    api.UtilApi
}

func (util *utilIot) SendBaiDuRequest(service api.IApiOuter, errorCode uint) api.ApiResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, err := mpf.JsonUnmarshalMap(resp.Content)
    if err != nil {
        info := make(map[string]string)
        if len(resp.Content) > 0 {
            info["result"] = resp.Content
        } else {
            info["result"] = "success"
        }
        result.Data = info
        return result
    }

    _, ok := respData["code"]
    if ok {
        result.Code = errorCode
        result.Msg = respData["message"].(string)
    } else {
        result.Data = respData
    }
    return result
}

func (util *utilIot) SendTencentRequest(service api.IApiOuter, errorCode uint) api.ApiResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["Response"]
    if ok {
        respInfo := respData["Response"].(map[string]interface{})
        _, ok = respInfo["Error"]
        if ok {
            errData := respInfo["Error"].(map[string]interface{})
            result.Code = errorCode
            result.Msg = errData["Message"].(string)
        } else {
            result.Data = respInfo
        }
    } else {
        result.Code = errorCode
        result.Msg = "解析服务数据出错"
    }
    return result
}

var (
    insUtil *utilIot
)

func init() {
    insUtil = &utilIot{api.NewUtilApi()}
}

func NewUtil() *utilIot {
    return insUtil
}
