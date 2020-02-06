/**
 * 外部接口公共结构体
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 9:13
 */
package api

import (
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type iApiCommon interface {
    GetRequest() (*fasthttp.Client, *fasthttp.Request) // 获取请求http操作对象
}

// 外部请求处理接口,针对有统一接口返回规范的api
type IApiOuter interface {
    iApiCommon
    GetReqTimeout() time.Duration
    CheckData() (*fasthttp.Client, *fasthttp.Request) // 数据校验
}

// 内部请求处理接口,针对没有统一接口返回规范的api
type IApiInner interface {
    iApiCommon
    checkData() (*fasthttp.Client, *fasthttp.Request) // 数据校验
    SendRequest() ApiResult                           // 发送请求
}

// api返回结果
type ApiResult struct {
    Code uint        `json:"code"` // 结果状态码
    Msg  string      `json:"msg"`  // 错误描述
    Data interface{} `json:"data"` // 结果数据
}

func NewApiResult() ApiResult {
    return ApiResult{0, "", make(map[string]interface{})}
}

type apiCommon struct {
    ReqUrl         string            // 请求地址
    ReqTimeout     time.Duration     // 请求超时时间,单位为纳秒
    ReqMethod      string            // 请求方式
    ReqContentType string            // 请求数据类型
    ReqData        map[string]string // 请求数据,如果是json格式请在checkData方法重新定义变量遍历
    ReqHeader      map[string]string // 请求头
}

func newApi() apiCommon {
    api := apiCommon{}
    api.ReqUrl = ""
    api.ReqTimeout = 3 * time.Second
    api.ReqMethod = fasthttp.MethodGet
    api.ReqContentType = project.HttpContentTypeForm
    api.ReqData = make(map[string]string)
    api.ReqHeader = make(map[string]string)
    return api
}

// 外部实现请求处理的结构体
type ApiOuter struct {
    apiCommon
}

func (api *ApiOuter) GetReqTimeout() time.Duration {
    return api.ReqTimeout
}

func NewApiOuter() ApiOuter {
    return ApiOuter{newApi()}
}

// 内部实现请求处理的结构体
type ApiInner struct {
    apiCommon
}

func (api *ApiInner) SendInner(client *fasthttp.Client, req *fasthttp.Request, errorCode uint) (mpf.HttpResp, ApiResult) {
    resp := mpf.HttpSendReq(client, req, api.ReqTimeout)
    result := NewApiResult()
    if resp.RespCode > 0 {
        result.Code = errorCode
        result.Msg = resp.Msg
    }

    return resp, result
}

func NewApiInner() ApiInner {
    return ApiInner{newApi()}
}
