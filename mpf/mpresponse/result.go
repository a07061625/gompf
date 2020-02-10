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
    "github.com/kataras/iris/v12"
)

type resultBasic struct {
    ReqId string `json:"req_id",xml:"ReqId"` // 请求ID
    Code  uint   `json:"code",xml:"Code"`    // 状态码
    Time  int64  `json:"time",xml:"Time"`    // 当前时间戳
    Msg   string `json:"msg",xml:"Msg"`      // 错误信息
}

func newResultBasic() resultBasic {
    r := resultBasic{}
    r.ReqId = mpf.ToolGetReqId()
    r.Code = errorcode.CommonBaseSuccess
    r.Time = time.Now().Unix()
    r.Msg = ""
    return r
}

// 接口响应,用于请求正常情况
type ResultApi struct {
    resultBasic
    Data interface{} `json:"data",xml:"Data"` // 响应数据
}

func NewResultApi() *ResultApi {
    return &ResultApi{newResultBasic(), make(map[string]interface{})}
}

// 问题响应,用于请求出问题情况
type ResultProblem struct {
    resultBasic
    Type   string `json:"type",xml:"Type"`     // 问题类型
    Title  string `json:"title",xml:"Title"`   // 问题标题
    Status int    `json:"status",xml:"Status"` // 问题状态
}

func NewResultProblem() *ResultProblem {
    return &ResultProblem{newResultBasic(), "", "", iris.StatusInternalServerError}
}
