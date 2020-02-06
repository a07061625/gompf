/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/2 0002
 * Time: 9:53
 */
package qcloud

import (
    "crypto/tls"
    "net/url"
    "sort"
    "strconv"
    "strings"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type baseQCloud struct {
    api.ApiOuter
}

func newBase() baseQCloud {
    return baseQCloud{api.NewApiOuter()}
}

type BaseCos struct {
    baseQCloud
    ReqUri         string            // 请求uri
    signParams     map[string]string // 参与签名的请求参数列表
    signHeaders    map[string]string // 参与签名的请求头列表
    signExpireTime int               // 签名有效时间,单位为秒
}

func (cos *BaseCos) SetParamData(key, val string) {
    lowerKey := strings.ToLower(key)
    cos.ReqData[key] = val
    cos.signParams[lowerKey] = val
}

func (cos *BaseCos) SetHeaderData(key, val string) {
    lowerKey := strings.ToLower(key)
    cos.ReqHeader[key] = val
    cos.signHeaders[lowerKey] = val
}

func (cos *BaseCos) SetSignExpireTime(signExpireTime int) {
    if signExpireTime > 0 {
        cos.signExpireTime = signExpireTime
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "签名有效时间不合法", nil))
    }
}

func (cos *BaseCos) createSign() {
    pkParams := mpf.NewHttpParamKey(cos.signParams)
    sort.Sort(pkParams)
    sortParams := pkParams.Params
    paramLength := len(sortParams)
    paramList := make([]string, 0)
    paramValues := url.Values{}
    for i := 0; i < paramLength; i++ {
        paramList = append(paramList, sortParams[i].Key)
        paramValues.Add(sortParams[i].Key, sortParams[i].Val)
    }

    pkHeaders := mpf.NewHttpParamKey(cos.signHeaders)
    sort.Sort(pkHeaders)
    sortHeaders := pkHeaders.Params
    headerLength := len(sortHeaders)
    headerList := make([]string, 0)
    headerValues := url.Values{}
    for j := 0; j < headerLength; j++ {
        headerList = append(headerList, sortHeaders[j].Key)
        headerValues.Add(sortHeaders[j].Key, sortHeaders[j].Val)
    }

    nowTime := time.Now().Second()
    endTime := nowTime + cos.signExpireTime
    signTime := strconv.Itoa(nowTime) + ";" + strconv.Itoa(endTime)

    conf := NewConfig().GetCos()
    signKey := mpf.HashSha1(signTime, conf.GetSecretKey())
    httpStr := strings.ToLower(cos.ReqMethod) + "\n" + cos.ReqUri + "\n" + paramValues.Encode() + "\n" + headerValues.Encode() + "\n"
    signStr := "sha1\n" + signTime + "\n" + mpf.HashSha1(httpStr, "") + "\n"

    cos.ReqHeader["Authorization"] = "q-sign-algorithm=sha1&q-ak=" + conf.GetSecretId()
    cos.ReqHeader["Authorization"] += "&q-sign-time=" + signTime
    cos.ReqHeader["Authorization"] += "&q-key-time=" + signTime
    cos.ReqHeader["Authorization"] += "&q-header-list=" + strings.Join(headerList, ";")
    cos.ReqHeader["Authorization"] += "&q-url-param-list=" + strings.Join(paramList, ";")
    cos.ReqHeader["Authorization"] += "&q-signature=" + mpf.HashSha1(signStr, signKey)
}

func (cos *BaseCos) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    cos.createSign()
    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.Header.SetRequestURI(cos.ReqUrl)
    req.Header.SetContentType(cos.ReqContentType)
    req.Header.SetMethod(cos.ReqMethod)
    mpf.HttpAddReqHeader(req, cos.ReqHeader)

    return client, req
}

func NewCos() BaseCos {
    conf := NewConfig().GetCos()
    cos := BaseCos{newBase(), "", make(map[string]string), make(map[string]string), 0}
    cos.ReqUri = "/"
    cos.signExpireTime = 30
    cos.SetHeaderData("Host", conf.GetBucketHost())
    cos.ReqContentType = project.HttpContentTypeXml
    cos.ReqMethod = fasthttp.MethodGet
    return cos
}
