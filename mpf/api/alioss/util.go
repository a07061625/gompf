/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 23:26
 */
package alioss

import (
    "encoding/base64"

    "github.com/a07061625/gompf/mpf"
)

type utilAliOss struct {
}

// 生成前端配置签名
func (util *utilAliOss) CreatePolicySign(config map[string]interface{}) map[string]string {
    conf := NewConfig()
    jsonStr := mpf.JsonMarshal(config)
    signStr := base64.StdEncoding.EncodeToString([]byte(jsonStr))
    sha1Str := mpf.HashSha1(signStr, conf.GetAccessKeySecret())
    sign := base64.StdEncoding.EncodeToString([]byte(sha1Str))

    result := make(map[string]string)
    result["policy_sign"] = signStr
    result["signature"] = sign
    return result
}

var (
    insUtil *utilAliOss
)

func init() {
    insUtil = &utilAliOss{}
}

func NewUtil() *utilAliOss {
    return insUtil
}
