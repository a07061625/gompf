/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 16:22
 */
package mpcurrency

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
)

func (util *utilCurrency) SendAMJiSuRequest(service api.IApiOuter, errorCode uint) api.ApiResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    status, ok := respData["status"]
    if ok && (status.(int) == 0) {
        result.Data = respData["result"]
    } else {
        result.Code = errorCode
        result.Msg = respData["msg"].(string)
    }
    return result
}
