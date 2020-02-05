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

func (util *utilLogistics) SendAMAliRequest(service api.IApiOuter, errorCode uint) api.ApiResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    resCode, ok := respData["showapi_res_code"]
    if ok && (resCode.(int) == 0) {
        resBody := respData["showapi_res_body"].(map[string]interface{})
        retCode, ok := resBody["ret_code"]
        if ok && (retCode.(int) == 0) {
            result.Data = resBody
        } else {
            result.Code = errorCode
            result.Msg = resBody["msg"].(string)
        }
    } else {
        result.Code = errorCode
        result.Msg = respData["showapi_res_error"].(string)
    }
    return result
}
