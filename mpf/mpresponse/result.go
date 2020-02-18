// Package mpresponse result
// User: 姜伟
// Time: 2020-02-19 05:04:11
package mpresponse

import (
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/kataras/iris/v12"
)

type resultBasic struct {
    ReqID string `json:"req_id" ,xml:"ReqId"` // 请求ID
    Code  uint   `json:"code" ,xml:"Code"`    // 状态码
    Time  int64  `json:"time" ,xml:"Time"`    // 当前时间戳
    Msg   string `json:"msg" ,xml:"Msg"`      // 错误信息
}

func newResultBasic() resultBasic {
    r := resultBasic{}
    r.refreshBasic()
    return r
}

func (r *resultBasic) refreshBasic() {
    r.ReqID = mpf.ToolGetReqId()
    r.Code = errorcode.CommonBaseSuccess
    r.Time = time.Now().Unix()
    r.Msg = ""
}

// ResultAPI 接口响应,用于请求正常情况
type ResultAPI struct {
    resultBasic
    Data interface{} `json:"data"  ,xml:"Data"` // 响应数据
}

// Refresh 刷新
func (r *ResultAPI) Refresh() {
    r.refreshBasic()
    r.Data = make(map[string]interface{})
}

// NewResultAPI 实例化
func NewResultAPI() *ResultAPI {
    r := &ResultAPI{}
    r.Refresh()
    return r
}

// ResultProblem 问题响应,用于请求出问题情况
type ResultProblem struct {
    resultBasic
    Tag    string `json:"tag" ,xml:"Tag"`       // 问题标识
    Title  string `json:"title" ,xml:"Title"`   // 问题标题
    Status int    `json:"status" ,xml:"Status"` // 问题状态
}

// Refresh 刷新
func (r *ResultProblem) Refresh() {
    r.refreshBasic()
    r.Tag = ""
    r.Title = ""
    r.Status = iris.StatusInternalServerError
}

// NewResultProblem 实例化
func NewResultProblem() *ResultProblem {
    r := &ResultProblem{}
    r.Refresh()
    return r
}
