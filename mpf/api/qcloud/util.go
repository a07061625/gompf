/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/2 0002
 * Time: 9:53
 */
package qcloud

import (
    "encoding/base64"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/clbanning/mxj"
)

type utilQCloud struct {
    api.UtilApi
}

// 生成权限策略签名
func (util *utilQCloud) CreatePolicySign(policyConfig map[string]interface{}) map[string]string {
    nowTime := time.Now().Unix()
    endTime := nowTime + 259200

    conf := NewConfig().GetCos()
    result := make(map[string]string)
    result["q-sign-algorithm"] = "sha1"
    result["q-ak"] = conf.GetSecretId()
    result["q-key-time"] = strconv.FormatInt(nowTime, 10) + ";" + strconv.FormatInt(endTime, 10)

    policyInfo := make(map[string]interface{})
    et := time.Unix(int64(endTime), 0)
    policyInfo["expiration"] = et.Format("2006-01-02T03:04:05.000Z")
    policyInfo["conditions"] = policyConfig
    policy := mpf.JsonMarshal(policyInfo)
    result["policy"] = base64.StdEncoding.EncodeToString([]byte(policy))

    policySign := mpf.HashSha1(policy, "")
    signKey := mpf.HashSha1(result["q-key-time"], conf.GetSecretKey())
    result["q-signature"] = mpf.HashSha1(policySign, signKey)
    return result
}

func (util *utilQCloud) SendCosRequest(service api.IApiOuter, errorCode uint) api.ApiResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    saveRes := make(map[string]interface{})
    saveRes["code"] = resp.StatusCode
    saveRes["headers"] = resp.Headers
    if (resp.ContentLength == 0) || (resp.Headers["Content-Type"] != "application/xml") {
        saveRes["content"] = resp.Content
        result.Data = saveRes
        return result
    }

    xmlData, _ := mxj.NewMapXml(resp.Body)
    _, ok := xmlData["Error"]
    if ok {
        errInfo := xmlData["Error"].(map[string]interface{})
        result.Code = errorCode
        result.Msg = errInfo["Message"].(string)
    } else {
        saveRes["content"] = xmlData
        result.Data = saveRes
    }
    return result
}

var (
    insUtil *utilQCloud
)

func init() {
    insUtil = &utilQCloud{api.NewUtilApi()}
}

func NewUtil() *utilQCloud {
    return insUtil
}
