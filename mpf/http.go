/**
 * http操作
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 11:05
 */
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

type HttpParam struct {
    Key string
    Val string
}

type HttpParamSorter struct {
    Params []HttpParam
}

func (ps HttpParamSorter) Len() int {
    return len(ps.Params)
}

func (ps HttpParamSorter) Swap(i, j int) {
    ps.Params[i], ps.Params[j] = ps.Params[j], ps.Params[i]
}

type HttpParamKey struct {
    HttpParamSorter
}

func (pk HttpParamKey) Less(i, j int) bool {
    return pk.Params[i].Key < pk.Params[j].Key
}

func NewHttpParamKey(data map[string]string) *HttpParamKey {
    pk := &HttpParamKey{}
    pk.Params = make([]HttpParam, 0, len(data))
    for k, v := range data {
        pk.Params = append(pk.Params, HttpParam{k, v})
    }
    return pk
}

type HttpParamVal struct {
    HttpParamSorter
}

func (pv HttpParamVal) Less(i, j int) bool {
    return pv.Params[i].Val < pv.Params[j].Val
}

func NewHttpParamVal(data map[string]string) *HttpParamVal {
    pv := &HttpParamVal{}
    pv.Params = make([]HttpParam, 0, len(data))
    for k, v := range data {
        pv.Params = append(pv.Params, HttpParam{k, v})
    }
    return pv
}

type HttpResp struct {
    RespCode      uint                `json:"resp_code"` // 响应状态码
    Msg           string              `json:"msg"`       // 错误信息
    Content       string              `json:"content"`   // 响应内容
    Body          []byte              // 响应体源数据
    ContentLength int                 `json:"content_length"` // 响应内容长度
    StatusCode    int                 `json:"status_code"`    // HTTP状态码
    Headers       map[string]string   `json:"headers"`        // 响应头
    Cookies       map[string][]string `json:"cookies"`        // cookies
}

func NewHttpResp() HttpResp {
    return HttpResp{0, "", "", []byte(""), 0, 0, make(map[string]string), make(map[string][]string)}
}

type HttpRespResult struct {
    ReqId string                 `json:"req_id"` // 请求ID
    Time  int64                  `json:"time"`   // 请求时间
    Code  uint                   `json:"code"`   // 状态码
    Msg   string                 `json:"msg"`    // 错误信息
    Data  map[string]interface{} `json:"data"`   // 响应数据
}

func NewHttpRespResult() HttpRespResult {
    nowTime := time.Now().Unix()
    nonceStr := strconv.FormatInt(nowTime, 10) + ToolCreateNonceStr(8, "numlower")
    reqId := HashMd5(nonceStr, "")
    return HttpRespResult{reqId, nowTime, errorcode.CommonBaseSuccess, "", make(map[string]interface{})}
}

var (
    httpReqId = ""
)

// 生成请求ID
func HttpReqId() string {
    if len(httpReqId) != 32 {
        nowTime := time.Now().UnixNano() / 1000
        str := strconv.FormatInt(nowTime, 10) + ToolCreateNonceStr(8, "lower")
        httpReqId = HashMd5(str, "")
    }

    return httpReqId
}

// 发送http请求
func HttpSendReq(client *fasthttp.Client, req *fasthttp.Request, timeout time.Duration) HttpResp {
    defer fasthttp.ReleaseRequest(req)

    result := NewHttpResp()
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
            if (headerKey != project.HttpHeadKeyContentLength) && (headerKey != project.HttpHeadKeyCookie) {
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

// 生成请求参数,生成的结果按照键名的ASCII码排序
//   data map[string]string 原数据
//   sortType string 排序类型 none:不排序 key:键名升序 val:键值升序
//   encodeType int 编码类型
//     1:url_encode(k1=v1&k2=v2&...)
//     2:json
//     3:xml
//     4:k1=v1&k2=v2&...
//     5:k1v1k2v2...
func HttpCreateParams(data map[string]string, sortType string, encodeType int) string {
    params := make([]HttpParam, 0)

    switch sortType {
    case "key":
        pk := NewHttpParamKey(data)
        sort.Sort(pk)
        params = pk.Params
    case "val":
        pv := NewHttpParamVal(data)
        sort.Sort(pv)
        params = pv.Params
    default:
        for k, v := range data {
            params = append(params, HttpParam{k, v})
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
        return JsonMarshal(sortData)
    case 3:
        sortData := make(map[string]string)
        for i := 0; i < paramNum; i++ {
            sortData[params[i].Key] = params[i].Val
        }
        xm := XmlMap{}
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

// 添加请求的头信息
func HttpAddReqHeader(req *fasthttp.Request, headers map[string]string) {
    for k, v := range headers {
        req.Header.Add(k, v)
    }
}
