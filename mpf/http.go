// Package mpf http
// User: 姜伟
// Time: 2020-02-19 05:30:24
package mpf

import (
    "encoding/xml"
    "net/url"
    "sort"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/valyala/fasthttp"
)

// HTTPParam HTTPParam
type HTTPParam struct {
    Key string
    Val string
}

// HTTPParamSorter HTTPParamSorter
type HTTPParamSorter struct {
    Params []HTTPParam
}

// Len Len
func (ps HTTPParamSorter) Len() int {
    return len(ps.Params)
}

// Swap Swap
func (ps HTTPParamSorter) Swap(i, j int) {
    ps.Params[i], ps.Params[j] = ps.Params[j], ps.Params[i]
}

// HTTPParamKey HTTPParamKey
type HTTPParamKey struct {
    HTTPParamSorter
}

// Less Less
func (pk HTTPParamKey) Less(i, j int) bool {
    return pk.Params[i].Key < pk.Params[j].Key
}

// NewHTTPParamKey NewHTTPParamKey
func NewHTTPParamKey(data map[string]string) *HTTPParamKey {
    pk := &HTTPParamKey{}
    pk.Params = make([]HTTPParam, 0, len(data))
    for k, v := range data {
        pk.Params = append(pk.Params, HTTPParam{k, v})
    }
    return pk
}

// HTTPParamVal HTTPParamVal
type HTTPParamVal struct {
    HTTPParamSorter
}

// Less Less
func (pv HTTPParamVal) Less(i, j int) bool {
    return pv.Params[i].Val < pv.Params[j].Val
}

// NewHTTPParamVal NewHTTPParamVal
func NewHTTPParamVal(data map[string]string) *HTTPParamVal {
    pv := &HTTPParamVal{}
    pv.Params = make([]HTTPParam, 0, len(data))
    for k, v := range data {
        pv.Params = append(pv.Params, HTTPParam{k, v})
    }
    return pv
}

// HTTPResp HTTPResp
type HTTPResp struct {
    RespCode      uint                `json:"resp_code"` // 响应状态码
    Msg           string              `json:"msg"`       // 错误信息
    Content       string              `json:"content"`   // 响应内容
    Body          []byte              // 响应体源数据
    ContentLength int                 `json:"content_length"` // 响应内容长度
    StatusCode    int                 `json:"status_code"`    // HTTP状态码
    Headers       map[string]string   `json:"headers"`        // 响应头
    Cookies       map[string][]string `json:"cookies"`        // cookies
}

// NewHTTPResp NewHTTPResp
func NewHTTPResp() HTTPResp {
    return HTTPResp{0, "", "", []byte(""), 0, 0, make(map[string]string), make(map[string][]string)}
}

// HTTPRespResult HTTPRespResult
type HTTPRespResult struct {
    ReqID string                 `json:"req_id"` // 请求ID
    Time  int64                  `json:"time"`   // 请求时间
    Code  uint                   `json:"code"`   // 状态码
    Msg   string                 `json:"msg"`    // 错误信息
    Data  map[string]interface{} `json:"data"`   // 响应数据
}

// NewHTTPRespResult NewHTTPRespResult
func NewHTTPRespResult() HTTPRespResult {
    nowTime := time.Now().Unix()
    nonceStr := strconv.FormatInt(nowTime, 10) + ToolCreateNonceStr(8, "numlower")
    reqID := HashMd5(nonceStr, "")
    return HTTPRespResult{reqID, nowTime, errorcode.CommonBaseSuccess, "", make(map[string]interface{})}
}

var (
    httpReqID = ""
)

// HTTPReqID 生成请求ID
func HTTPReqID() string {
    if len(httpReqID) != 32 {
        nowTime := time.Now().UnixNano() / 1000
        str := strconv.FormatInt(nowTime, 10) + ToolCreateNonceStr(8, "lower")
        httpReqID = HashMd5(str, "")
    }

    return httpReqID
}

// HTTPSendReq 发送http请求
func HTTPSendReq(client *fasthttp.Client, req *fasthttp.Request, timeout time.Duration) HTTPResp {
    defer fasthttp.ReleaseRequest(req)

    result := NewHTTPResp()
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    err := client.DoTimeout(req, resp, timeout)
    if err == nil {
        result.Body = resp.Body()
        result.Content = string(resp.Body())
        result.ContentLength = resp.Header.ContentLength()
        result.StatusCode = resp.Header.StatusCode()
        resp.Header.VisitAll(func(key, value []byte) {
            headerKey := string(key)
            if (headerKey != project.HTTPHeadKeyContentLength) && (headerKey != project.HTTPHeadKeyCookie) {
                result.Headers[headerKey] = string(value)
            }
        })
        resp.Header.VisitAllCookie(func(key, value []byte) {
            result.Cookies[string(key)] = append(result.Cookies[string(key)], string(value))
        })
    } else {
        mplog.LogError("send http req fail,reason: " + err.Error())
        result.RespCode = errorcode.CommonRequestFail
        result.Msg = err.Error()
    }

    return result
}

// HTTPCreateParams 生成请求参数,生成的结果按照键名的ASCII码排序
//   data map[string]string 原数据
//   sortType string 排序类型 none:不排序 key:键名升序 val:键值升序
//   encodeType int 编码类型
//     1:url_encode(k1=v1&k2=v2&...)
//     2:json
//     3:xml
//     4:k1=v1&k2=v2&...
//     5:k1v1k2v2...
func HTTPCreateParams(data map[string]string, sortType string, encodeType int) string {
    params := make([]HTTPParam, 0)

    switch sortType {
    case "key":
        pk := NewHTTPParamKey(data)
        sort.Sort(pk)
        params = pk.Params
    case "val":
        pv := NewHTTPParamVal(data)
        sort.Sort(pv)
        params = pv.Params
    default:
        for k, v := range data {
            params = append(params, HTTPParam{k, v})
        }
    }
    paramNum := len(params)

    switch encodeType {
    case 1:
        values := url.Values{}
        for i := 0; i < paramNum; i++ {
            values.Add(params[i].Key, params[i].Val)
        }
        return values.Encode()
    case 2:
        sortData := make(map[string]string)
        for i := 0; i < paramNum; i++ {
            sortData[params[i].Key] = params[i].Val
        }
        return JSONMarshal(sortData)
    case 3:
        sortData := make(map[string]string)
        for i := 0; i < paramNum; i++ {
            sortData[params[i].Key] = params[i].Val
        }
        xm := XMLMap{}
        xm = sortData
        res, _ := xml.Marshal(xm)
        return string(res)
    case 4:
        str := ""
        for i := 0; i < paramNum; i++ {
            str += "&" + params[i].Key + "=" + params[i].Val
        }
        return str[1:]
    default:
        str := ""
        for i := 0; i < paramNum; i++ {
            str += params[i].Key + params[i].Val
        }
        return str
    }
}

// HTTPAddReqHeader 添加请求的头信息
func HTTPAddReqHeader(req *fasthttp.Request, headers map[string]string) {
    for k, v := range headers {
        req.Header.Add(k, v)
    }
}
