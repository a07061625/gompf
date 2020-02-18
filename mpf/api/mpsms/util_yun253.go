/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 16:22
 */
package mpsms

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
)

func (util *utilSms) SendYun253Request(service api.IAPIOuter, errorCode uint) api.APIResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["code"]
    if ok {
        result.Data = respData
        return result
    }

    result.Code = errorCode
    errMsg, ok := respData["errorMsg"]
    if ok {
        result.Msg = errMsg.(string)
    } else {
        result.Msg = "解析服务数据出错"
    }
    return result
}
