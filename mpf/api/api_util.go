// Package api api_util
// User: 姜伟
// Time: 2020-02-19 06:51:27
package api

import (
    "github.com/a07061625/gompf/mpf"
)

// UtilAPI UtilAPI
type UtilAPI struct {
}

// SendOuter SendOuter
func (util *UtilAPI) SendOuter(service IAPIOuter, errorCode uint) (mpf.HTTPResp, APIResult) {
    client, req := service.CheckData()
    resp := mpf.HTTPSendReq(client, req, service.GetReqTimeout())
    result := NewAPIResult()
    if resp.RespCode > 0 {
        result.Code = errorCode
        result.Msg = resp.Msg
    }

    return resp, result
}

// NewUtilAPI NewUtilAPI
func NewUtilAPI() UtilAPI {
    return UtilAPI{}
}
