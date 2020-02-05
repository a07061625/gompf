/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 16:38
 */
package mplogistics

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
)

func (util *utilLogistics) SendKd100Request(service api.IApiOuter, errorCode uint) api.ApiResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    resStatus, ok := respData["status"]
    if ok && (resStatus.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorCode
        result.Msg = respData["message"].(string)
    }
    return result
}
