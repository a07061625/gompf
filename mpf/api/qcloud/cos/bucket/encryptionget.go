/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 15:07
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/valyala/fasthttp"
)

// 获取默认加密配置
type encryptionGet struct {
    qcloud.BaseCos
}

func (eg *encryptionGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    eg.ReqUrl = "http://" + eg.ReqHeader["Host"] + eg.ReqUri + "?encryption"
    return eg.GetRequest()
}

func NewEncryptionGet() *encryptionGet {
    eg := &encryptionGet{qcloud.NewCos()}
    eg.SetParamData("encryption", "")
    return eg
}
