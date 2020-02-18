/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 22:27
 */
package mppush

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
)

type utilPush struct {
    api.UtilAPI
}

func (util *utilPush) SendXinGeRequest(service api.IAPIOuter, errorCode uint) api.APIResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["ret_code"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorCode
        result.Msg = respData["err_msg"].(string)
    }
    return result
}

func (util *utilPush) SendBaiDuRequest(service api.IAPIOuter, errorCode uint) api.APIResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["error_code"]
    if ok {
        result.Code = errorCode
        result.Msg = respData["error_msg"].(string)
        return result
    }
    _, ok = respData["response_params"]
    if ok {
        result.Data = respData["response_params"]
    } else {
        result.Data = make(map[string]string)
    }
    return result
}

func (util *utilPush) SendJPushRequest(service api.IAPIOuter, errorCode uint) api.APIResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["error"]
    if ok {
        errorInfo := respData["error"].(map[string]interface{})
        result.Code = errorCode
        result.Msg = errorInfo["message"].(string)
    } else {
        result.Data = respData
    }
    return result
}

var (
    insUtil *utilPush
)

func init() {
    insUtil = &utilPush{api.NewUtilAPI()}
}

func NewUtil() *utilPush {
    return insUtil
}
