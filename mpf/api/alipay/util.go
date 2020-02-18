/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 19:48
 */
package alipay

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "fmt"
    "hash"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

type IAliPayBase interface {
    api.IAPIOuter
    GetRespTag() string
}

type IAliPayOuter interface {
}

type utilAliPay struct {
    api.UtilAPI
    outer IAliPayOuter
}

// 获取待签名字符串
func (util *utilAliPay) sortData(data map[string]string) string {
    pk := mpf.NewHTTPParamKey(data)
    sort.Sort(pk)
    str := ""
    for _, param := range pk.Params {
        trueVal := strings.TrimSpace(param.Val)
        if (len(trueVal) > 0) && (trueVal[:1] != "@") {
            str += "&" + param.Key + "=" + param.Val
        }
    }
    return str[1:]
}

// 生成签名
//   data map[string]string 待签名数据
//   signType string 签名方式,只支持RSA和RSA2
func (util *utilAliPay) CreateSign(data map[string]string, signType string) string {
    conf := NewConfig().GetAccount(data["app_id"])
    // 解析秘钥
    block, _ := pem.Decode([]byte(conf.GetPriKeyRsa()))
    if block == nil {
        panic(mperr.NewAliPay(errorcode.AliPaySign, "签名错误", nil))
    }

    key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        panic(mperr.NewAliPay(errorcode.AliPaySign, "签名错误", err))
    }

    var (
        h1  hash.Hash
        h2  crypto.Hash
    )
    if signType == "RSA" {
        h1 = sha1.New()
        h2 = crypto.SHA1
    } else {
        h1 = sha256.New()
        h2 = crypto.SHA256
    }

    signStr := util.sortData(data)
    _, err = h1.Write([]byte(signStr))
    if err != nil {
        panic(mperr.NewAliPay(errorcode.AliPaySign, "签名错误", err))
    }
    // 调用算法
    encrypted, err := rsa.SignPKCS1v15(rand.Reader, key, h2, h1.Sum(nil))
    if err != nil {
        panic(mperr.NewAliPay(errorcode.AliPaySign, "签名错误", err))
    }

    return base64.StdEncoding.EncodeToString(encrypted)
}

// 校验签名
//   data map[string]string 待签名数据
//   verifyType int 校验类型 1：不校验数据签名类型 2：校验数据签名类型
//   signType string 签名方式,只支持RSA和RSA2
func (util *utilAliPay) VerifySign(data map[string]string, verifyType int, signType string) bool {
    nowSign, ok := data["sign"]
    if !ok {
        return false
    }

    delete(data, "sign")
    if verifyType == 1 {
        delete(data, "sign_type")
    }
    newSign := util.CreateSign(data, signType)

    signBytes, err := base64.StdEncoding.DecodeString(nowSign)
    if err != nil {
        return false
    }

    conf := NewConfig().GetAccount(data["app_id"])
    // 解析公钥
    block, _ := pem.Decode([]byte(conf.GetPubKeyAli()))
    if block == nil {
        return false
    }
    key, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return false
    }
    publicKey, ok := key.(*rsa.PublicKey)
    if !ok {
        return false
    }

    var h1 crypto.Hash
    if signType == "RSA" {
        h1 = crypto.SHA1
    } else {
        h1 = crypto.SHA256
    }

    h2 := h1.New()
    h2.Write([]byte(newSign))
    err = rsa.VerifyPKCS1v15(publicKey, h1, h2.Sum(nil), signBytes)
    return err != nil
}

// 发送服务请求
func (util *utilAliPay) SendRequest(service IAliPayBase, errorCode uint) api.APIResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    respTag := service.GetRespTag()
    _, ok := respData[respTag]
    if ok {
        resultData := respData[respTag].(map[string]interface{})
        resultCode, ok := resultData["code"]
        if ok && (resultCode.(string) == "10000") {
            result.Data = resultData
        } else {
            result.Code = errorCode
            result.Msg = resultData["sub_msg"].(string)
        }
    } else {
        result.Code = errorCode
        result.Msg = "解析服务数据出错"
    }
    return result
}

// 生成手机网页支付HTML代码
func (util *utilAliPay) CreatePayWapHtml(data map[string]interface{}) string {
    formId := "alipaywap" + mpf.ToolCreateNonceStr(6, "numlower") + strconv.FormatInt(time.Now().Unix(), 10)
    html := "<form id=\"" + formId + "\" name=\"" + formId + "\" action=\"" + UrlGateWay + "?charset=utf-8\" method=\"POST\">"
    for k, v := range data {
        tv := ""
        switch v.(type) {
        case string:
            tv = v.(string)
        case int:
            tv = strconv.Itoa(v.(int))
        case float32:
            tv = fmt.Sprintf("%.2f", v.(float32))
        }
        if len(tv) > 0 {
            html += "<input type=\"hidden\" name=\"" + k + "\" value=\"" + tv + "\" />"
        }
    }
    html += "<input type=\"submit\" value=\"ok\" style=\"display:none;\"/></form><script>document.forms[\"" + formId + "\"].submit();</script>"
    return html
}

var (
    onceUtil sync.Once
    insUtil  *utilAliPay
)

func init() {
    insUtil = &utilAliPay{api.NewUtilAPI(), nil}
}

func LoadUtil(outer IAliPayOuter) {
    onceUtil.Do(func() {
        insUtil.outer = outer
    })
}

func NewUtil() *utilAliPay {
    return insUtil
}
