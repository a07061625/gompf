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

func (util *utilLogistics) SendKdBirdRequest(service api.IAPIOuter, errorCode uint) api.APIResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    resTag, ok := respData["Success"]
    if ok && resTag.(bool) {
        result.Data = respData
    } else {
        result.Code = errorCode
        switch respData["Reason"].(type) {
        case string:
            result.Msg = respData["Reason"].(string)
        default:
            result.Msg = "请求错误"
        }
    }
    return result
}
