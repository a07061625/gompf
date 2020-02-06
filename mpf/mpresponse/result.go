/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/23 0023
 * Time: 10:09
 */
package mpresponse

import (
    "regexp"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 请求响应
type result struct {
    ReqId string                 `json:"req_id"`
    Code  uint                   `json:"code"`
    Data  map[string]interface{} `json:"data"`
    Msg   string                 `json:"msg"`
    Time  int64                  `json:"time"`
}

func NewResult(reqId string) *result {
    match := false
    if len(reqId) > 0 {
        match, _ = regexp.MatchString(`^[0-9a-z]{32}$`, reqId)
    }

    nowTime := time.Now().Unix()
    trueReqId := ""
    if match {
        trueReqId = reqId
    } else {
        needStr := mpf.ToolCreateNonceStr(8, "total") + strconv.FormatInt(nowTime, 10)
        trueReqId = mpf.HashMd5(needStr, "")
    }

    r := &result{}
    r.ReqId = trueReqId
    r.Code = errorcode.CommonBaseSuccess
    r.Data = make(map[string]interface{})
    r.Msg = ""
    r.Time = nowTime
    return r
}
