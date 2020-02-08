/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/23 0023
 * Time: 10:09
 */
package mpresponse

import (
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 请求响应
type ResultBasic struct {
    ReqId string                 `json:"req_id"`
    Code  uint                   `json:"code"`
    Data  map[string]interface{} `json:"data"`
    Msg   string                 `json:"msg"`
    Time  int64                  `json:"time"`
}

func NewResultBasic() *ResultBasic {
    r := &ResultBasic{}
    r.ReqId = mpf.ToolGetReqId()
    r.Code = errorcode.CommonBaseSuccess
    r.Data = make(map[string]interface{})
    r.Msg = ""
    r.Time = time.Now().Unix()
    return r
}
