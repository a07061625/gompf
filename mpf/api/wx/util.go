/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/25 0025
 * Time: 10:42
 */
package wx

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "encoding/binary"
    "encoding/xml"
    "errors"
    "io"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpencrypt"
    "github.com/a07061625/gompf/mpf/mperr"
)

type IWxOuter interface {
    GetOpenAuthorizeInfo(appId string) *DataOpenAuthorize
    RefreshOpenAuthorizeInfo(appId string, operateType int, data map[string]interface{})
    GetProviderAuthorizeInfo(corpId string) *DataProviderAuthorize
    RefreshProviderAuthorizeInfo(corpId string, operateType int, data map[string]interface{})
}

type utilWx struct {
    outer IWxOuter
}

// 生成SHA1算法安全签名
func (util *utilWx) CreateSha1Sign(token string, nowTime int, nonceStr string, encryptMsg string) string {
    s := []string{token, strconv.Itoa(nowTime), nonceStr, encryptMsg}
    sort.Strings(s)

    return mpf.HashSha1(strings.Join(s, ""), "")
}

// 生成企业签名
func (util *utilWx) CreateCropSign(data map[string]string, acceptKeys []string, agentSecret string) string {
    acceptData := make(map[string]string)
    for _, ak := range acceptKeys {
        eData, ok := data[ak]
        if !ok {
            panic(mperr.NewWxCorp(errorcode.WxCorpParam, "缺少字段"+ak, nil))
        }
        acceptData[ak] = eData
    }

    pk := mpf.NewHttpParamKey(acceptData)
    sort.Sort(pk)

    needStr1 := ""
    for _, param := range pk.Params {
        needStr1 += "&" + param.Key + "=" + param.Val
    }

    needStr2 := needStr1[1:] + agentSecret
    needStr3 := mpf.HashMd5(needStr2, "")
    return strings.ToUpper(needStr3)
}

// 生成企业支付签名
func (util *utilWx) CreateCropPaySign(data map[string]string, payKey string) string {
    pk := mpf.NewHttpParamKey(data)
    sort.Sort(pk)
    needStr1 := ""
    for _, param := range pk.Params {
        if param.Key == "sign" {
            continue
        }
        if len(param.Val) == 0 {
            continue
        }
        needStr1 += "&" + param.Key + "=" + param.Val
    }
    needStr2 := needStr1[1:] + payKey
    needStr3 := mpf.HashMd5(needStr2, "")
    return strings.ToUpper(needStr3)
}

func (util *utilWx) aesDecrypt(cipherData, aesKey []byte) ([]byte, error) {
    k := len(aesKey)
    if len(cipherData)%k != 0 {
        return nil, errors.New("密文长度不是密钥长度的整数倍")
    }

    block, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, err
    }

    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }

    blockMode := cipher.NewCBCDecrypter(block, iv)
    plainData := make([]byte, len(cipherData))
    blockMode.CryptBlocks(plainData, cipherData)
    return plainData, nil
}

// 解密xml数据
//   encryptXml string 加密xml数据
//   encodeAesKey string 编码后的加密密钥
//   appId string 应用ID,企业服务商为服务商corpId,开放平台为第三方平台的appId
func (util *utilWx) DecryptXml(encryptXml string, encodeAesKey string, appId string) string {
    aesKey, _ := base64.StdEncoding.DecodeString(encodeAesKey + "=")
    cipherData, err := base64.StdEncoding.DecodeString(encryptXml)
    if err != nil {
        panic(mperr.NewWx(errorcode.WxParam, "密钥解码失败", nil))
    }

    plainData, err := util.aesDecrypt(cipherData, aesKey)
    if err != nil {
        panic(mperr.NewWx(errorcode.WxParam, "解密xml数据失败", nil))
    }

    buf := bytes.NewBuffer(plainData[16:20])
    var length int32 = 0
    binary.Read(buf, binary.BigEndian, &length)
    appIdStart := 20 + length
    appIdEnd := int(appIdStart) + len(appId)
    id := plainData[appIdStart:appIdEnd]
    if appId != string(id) {
        panic(mperr.NewWx(errorcode.WxParam, "微信应用ID校验失败", nil))
    }

    return string(plainData[20:appIdStart])
}

func (util *utilWx) aesEncrypt(plainData, aesKey []byte) ([]byte, error) {
    k := len(aesKey)
    if len(plainData)%k != 0 {
        plainData = mpencrypt.AesPaddingPKCS7(plainData, k)
    }

    block, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, err
    }

    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }

    cipherData := make([]byte, len(plainData))
    blockMode := cipher.NewCBCEncrypter(block, iv)
    blockMode.CryptBlocks(cipherData, plainData)
    return cipherData, nil
}

// 加密xml数据
//   replyMsg string xml数据
//   appId string 应用ID,企业服务商为服务商corpId,开放平台为第三方平台的appId
//   appToken string 加解密令牌
//   encodeAesKey string 编码后的加密密钥
func (util *utilWx) EncryptXml(replyMsg string, appId string, appToken string, encodeAesKey string) string {
    replyData := []byte(replyMsg)
    buf := new(bytes.Buffer)
    err := binary.Write(buf, binary.BigEndian, int32(len(replyData)))
    if err != nil {
        panic(mperr.NewWx(errorcode.WxParam, "加密xml数据失败", nil))
    }

    bufLength := buf.Bytes()
    nonceStr := mpf.ToolCreateNonceStr(16, "numlower")
    nonceBytes := []byte(nonceStr)
    aesKey, _ := base64.StdEncoding.DecodeString(encodeAesKey + "=")
    plainData := bytes.Join([][]byte{nonceBytes, bufLength, replyData, []byte(appId)}, nil)
    cipherData, err := util.aesEncrypt(plainData, aesKey)
    if err != nil {
        panic(mperr.NewWx(errorcode.WxParam, "加密xml数据失败", nil))
    }

    encryptMsg := base64.StdEncoding.EncodeToString(cipherData)
    nowTime := time.Now().Second()
    sign := util.CreateSha1Sign(appToken, nowTime, nonceStr, encryptMsg)
    resp := &WxResponse{}
    resp.Encrypt = WxCDATAText{"<![CDATA[" + encryptMsg + "]]>"}
    resp.MsgSignature = WxCDATAText{"<![CDATA[" + sign + "]]>"}
    resp.TimeStamp = strconv.Itoa(nowTime)
    resp.Nonce = WxCDATAText{"<![CDATA[" + nonceStr + "]]>"}
    encryptXml, err := xml.MarshalIndent(resp, " ", "  ")
    if err != nil {
        panic(mperr.NewWx(errorcode.WxParam, "加密xml数据失败", nil))
    }

    return string(encryptXml)
}

// 解密小程序用户数据
//   encodeData string 编码后的原始数据
//   encodeKey string 编码后的密钥
//   encodeIv string 编码后的初始向量
//   appId string 小程序应用ID
func (util *utilWx) DecryptMiniData(encodeData, encodeKey, encodeIv, appId string) map[string]interface{} {
    data, err := base64.StdEncoding.DecodeString(encodeData)
    if err != nil {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "解密小程序原始数据失败", nil))
    }
    key, err := base64.StdEncoding.DecodeString(encodeKey)
    if err != nil {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "解密小程序密钥失败", nil))
    }
    iv, err := base64.StdEncoding.DecodeString(encodeIv)
    if err != nil {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "解密小程序初始向量失败", nil))
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "解密小程序数据失败", nil))
    }

    blockSize := block.BlockSize()
    if len(data) < blockSize {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "解密小程序数据失败", nil))
    }
    if len(data)%blockSize != 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "解密小程序数据失败", nil))
    }
    blockMode := cipher.NewCBCDecrypter(block, iv)
    plainData := make([]byte, len(data))
    blockMode.CryptBlocks(plainData, data)

    jsonStr := string(mpencrypt.AesUnPaddingPKCS7(plainData))
    jsonData, _ := mpf.JsonUnmarshalMap(jsonStr)
    watermark, ok := jsonData["watermark"]
    if !ok {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "小程序数据格式错误", nil))
    }
    watermarkData := watermark.(map[string]interface{})
    id, ok := watermarkData["appid"]
    if !ok {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "小程序数据格式错误", nil))
    } else if id.(string) != appId {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "非法的小程序数据", nil))
    }

    return jsonData
}

// 生成公众号或小程序支付签名
func (util *utilWx) CreateSinglePaySign(data map[string]string, appId string, signType string) string {
    pk := mpf.NewHttpParamKey(data)
    sort.Sort(pk)
    needStr1 := ""
    for _, param := range pk.Params {
        if param.Key == "sign" {
            continue
        }
        if len(param.Val) == 0 {
            continue
        }
        needStr1 += param.Key + "=" + param.Val + "&"
    }

    payKey := NewConfig().GetAccount(appId).GetPayKey()
    needStr2 := needStr1 + "key=" + payKey
    needStr3 := ""
    if signType == "md5" {
        needStr3 = mpf.HashMd5(needStr2, "")
    } else {
        needStr3 = mpf.HashSha256(needStr2, payKey)
    }
    return strings.ToUpper(needStr3)
}

// 检验公众号或小程序支付签名
func (util *utilWx) CheckSinglePaySign(data map[string]string, appId string) bool {
    sign, ok := data["sign"]
    if ok {
        nowSign := util.CreateSinglePaySign(data, appId, "md5")
        return sign == nowSign
    }
    return false
}

// 获取公众号或小程序缓存
//   appId string 公众号或小程序应用ID
//   getType string 获取类型
//     single_accesstoken: 获取公众号或小程序访问令牌
//     single_jsticket: 获取公众号或小程序js ticket
//     single_cardticket: 获取公众号或小程序卡券 ticket
//     open_accesstoken: 获取开放平台授权者访问令牌
//     open_jsticket: 获取开放平台授权者js ticket
//     open_cardticket: 获取开放平台授权者卡券 ticket
func (util *utilWx) GetSingleCache(appId, getType string) string {
    switch getType {
    case SingleCacheTypeAccessToken:
        return util.GetSingleAccessToken(appId)
    case SingleCacheTypeJsTicket:
        return util.GetSingleJsTicket(appId)
    case SingleCacheTypeCardTicket:
        return util.GetSingleCardTicket(appId)
    case SingleCacheTypeOpenAccessToken:
        return util.GetOpenAuthorizeAccessToken(appId)
    case SingleCacheTypeOpenJsTicket:
        return util.GetOpenAuthorizeJsTicket(appId)
    case SingleCacheTypeOpenCardTicket:
        return util.GetOpenAuthorizeCardTicket(appId)
    default:
        panic(mperr.NewWx(errorcode.WxParam, "获取类型不支持", nil))
    }
}

// 获取企业号缓存
//   corpId string 企业号
//   agentTag string 应用标识
//   getType string 获取类型
//     corp_accesstoken: 获取企业号访问令牌
//     corp_jsticket: 获取企业号js ticket
//     provider_accesstoken: 获取服务商授权者访问令牌
//     provider_jsticket: 获取服务商授权者js ticket
func (util *utilWx) GetCorpCache(corpId, agentTag, getType string) string {
    switch getType {
    case CorpCacheTypeAccessToken:
        return util.GetCorpAccessToken(corpId, agentTag)
    case CorpCacheTypeJsTicket:
        return util.GetCorpJsTicket(corpId, agentTag)
    case CorpCacheTypeProviderAccessToken:
        return util.GetProviderAuthorizeAccessToken(corpId)
    case CorpCacheTypeProviderJsTicket:
        return util.GetProviderAuthorizeJsTicket(corpId)
    default:
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "获取类型不支持", nil))
    }
}

var (
    onceUtilWx sync.Once
    insUtilWx  *utilWx
)

func init() {
    insUtilWx = &utilWx{}
}

func LoadUtil(outer IWxOuter) {
    onceUtilWx.Do(func() {
        insUtilWx.outer = outer
    })
}

func NewUtilWx() *utilWx {
    return insUtilWx
}
