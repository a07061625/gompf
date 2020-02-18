/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 10:57
 */
package api

import (
    "github.com/a07061625/gompf/mpf"
)

type UtilApi struct {
}

func (util *UtilApi) SendOuter(service IApiOuter, errorCode uint) (mpf.HTTPResp, ApiResult) {
    client, req := service.CheckData()
    resp := mpf.HTTPSendReq(client, req, service.GetReqTimeout())
    result := NewApiResult()
    if resp.RespCode > 0 {
        result.Code = errorCode
        result.Msg = resp.Msg
    }

    return resp, result
}

func NewUtilApi() UtilApi {
    return UtilApi{}
}
