// Package api api_base
// User: 姜伟
// Time: 2020-02-19 06:43:49
package api

import (
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type iAPICommon interface {
    GetRequest() (*fasthttp.Client, *fasthttp.Request) // 获取请求http操作对象
}

// IAPIOuter 外部请求处理接口,针对有统一接口返回规范的api
type IAPIOuter interface {
    iAPICommon
    GetReqTimeout() time.Duration
    CheckData() (*fasthttp.Client, *fasthttp.Request) // 数据校验
}

// IAPIInner 内部请求处理接口,针对没有统一接口返回规范的api
type IAPIInner interface {
    iAPICommon
    checkData() (*fasthttp.Client, *fasthttp.Request) // 数据校验
    SendRequest() APIResult                           // 发送请求
}

// APIResult api返回结果
type APIResult struct {
    Code uint        `json:"code"` // 结果状态码
    Msg  string      `json:"msg"`  // 错误描述
    Data interface{} `json:"data"` // 结果数据
}

// NewAPIResult NewAPIResult
func NewAPIResult() APIResult {
    return APIResult{0, "", make(map[string]interface{})}
}

type apiCommon struct {
    ReqURI         string            // 请求地址
    ReqTimeout     time.Duration     // 请求超时时间,单位为纳秒
    ReqMethod      string            // 请求方式
    ReqContentType string            // 请求数据类型
    ReqData        map[string]string // 请求数据,如果是json格式请在checkData方法重新定义变量遍历
    ReqHeader      map[string]string // 请求头
}

func newAPI() apiCommon {
    api := apiCommon{}
    api.ReqURI = ""
    api.ReqTimeout = 3 * time.Second
    api.ReqMethod = fasthttp.MethodGet
    api.ReqContentType = project.HTTPContentTypeForm
    api.ReqData = make(map[string]string)
    api.ReqHeader = make(map[string]string)
    return api
}

// APIOuter 外部实现请求处理的结构体
type APIOuter struct {
    apiCommon
}

// GetReqTimeout GetReqTimeout
func (api *APIOuter) GetReqTimeout() time.Duration {
    return api.ReqTimeout
}

// NewAPIOuter NewAPIOuter
func NewAPIOuter() APIOuter {
    return APIOuter{newAPI()}
}

// APIInner 内部实现请求处理的结构体
type APIInner struct {
    apiCommon
}

// SendInner SendInner
func (api *APIInner) SendInner(client *fasthttp.Client, req *fasthttp.Request, errorCode uint) (mpf.HTTPResp, APIResult) {
    resp := mpf.HTTPSendReq(client, req, api.ReqTimeout)
    result := NewAPIResult()
    if resp.RespCode > 0 {
        result.Code = errorCode
        result.Msg = resp.Msg
    }

    return resp, result
}

// NewAPIInner NewAPIInner
func NewAPIInner() APIInner {
    return APIInner{newAPI()}
}
