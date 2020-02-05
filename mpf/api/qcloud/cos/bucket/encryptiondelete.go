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

// 删除默认加密配置
type encryptionDelete struct {
    qcloud.BaseCos
}

func (ed *encryptionDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    ed.ReqUrl = "http://" + ed.ReqHeader["Host"] + ed.ReqUri + "?encryption"
    return ed.GetRequest()
}

func NewEncryptionDelete() *encryptionDelete {
    ed := &encryptionDelete{qcloud.NewCos()}
    ed.ReqMethod = fasthttp.MethodDelete
    ed.SetParamData("encryption", "")
    return ed
}
