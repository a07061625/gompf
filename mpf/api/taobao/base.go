/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/2 0002
 * Time: 9:30
 */
package taobao

import (
    "crypto/tls"
    "strings"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/valyala/fasthttp"
)

type BaseTaoBao struct {
    api.APIOuter
    respTag   string // 响应标识
    AppKey    string // 应用标识
    AppSecret string // 应用密钥
}

func (b *BaseTaoBao) SetMethod(method string) {
    b.ReqData["method"] = method
    trueMethod := strings.TrimPrefix(method, "taobao.")
    b.respTag = strings.Replace(trueMethod, ".", "_", -1) + "_response"
}

func (b *BaseTaoBao) GetRespTag() string {
    return b.respTag
}

// 生成签名
func (b *BaseTaoBao) createSign() {
    delete(b.ReqData, "sign")
    signStr := mpf.HTTPCreateParams(b.ReqData, "key", 5)
    sign := ""
    if b.ReqData["sign_method"] == "md5" {
        sign = mpf.HashMd5(b.AppSecret+signStr+b.AppSecret, "")
    } else {
        sign = mpf.HashMd5(signStr, b.AppSecret)
    }
    b.ReqData["sign"] = strings.ToUpper(sign)
}

func (b *BaseTaoBao) GetRequest() (*fasthttp.Client, *fasthttp.Request) {
    b.ReqData["app_key"] = b.AppKey
    b.createSign()
    reqBody := mpf.HTTPCreateParams(b.ReqData, "key", 1)

    client := &fasthttp.Client{}
    client.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    req := fasthttp.AcquireRequest()
    req.SetBody([]byte(reqBody))
    req.Header.SetRequestURI(b.ReqURI)
    req.Header.SetMethod(b.ReqMethod)
    mpf.HTTPAddReqHeader(req, b.ReqHeader)

    return client, req
}

func NewBaseTaoBao() BaseTaoBao {
    now := time.Now()
    b := BaseTaoBao{api.NewAPIOuter(), "", "", ""}
    b.ReqData["v"] = "2.0"
    b.ReqData["sign_method"] = "md5"
    b.ReqData["format"] = "json"
    b.ReqData["simplify"] = "1"
    b.ReqData["timestamp"] = now.Format("2006-01-02 03:04:05")
    b.ReqHeader["Expect"] = ""
    b.ReqMethod = fasthttp.MethodPost
    b.ReqURI = TaoBaoUrlEnv
    return b
}
