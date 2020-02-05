/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 18:21
 */
package alioss

import (
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
    onceClient sync.Once
    insClient  *oss.Client
)

func init() {
    insClient = nil
}

func NewClient() *oss.Client {
    onceClient.Do(func() {
        clientOptions := make([]oss.ClientOption, 0)
        clientOptions = append(clientOptions, oss.EnableCRC(true))
        clientOptions = append(clientOptions, oss.Timeout(10, 120))
        conf1 := mpf.NewConfig().GetConfig("alioss")
        initType := conf1.GetInt(mpf.EnvProjectKey() + ".inittype")
        switch initType {
        case 1:
        case 2:
            clientOptions = append(clientOptions, oss.UseCname(true))
        case 3:
            securityToken := conf1.GetString(mpf.EnvProjectKey() + ".security.token")
            clientOptions = append(clientOptions, oss.SecurityToken(securityToken))
        case 4:
            proxyUrl := conf1.GetString(mpf.EnvProjectKey() + ".proxy.url")
            clientOptions = append(clientOptions, oss.Proxy(proxyUrl))
        case 5:
            proxyUrl := conf1.GetString(mpf.EnvProjectKey() + ".proxy.url")
            proxyUser := conf1.GetString(mpf.EnvProjectKey() + ".proxy.user")
            proxyPassword := conf1.GetString(mpf.EnvProjectKey() + ".proxy.password")
            clientOptions = append(clientOptions, oss.AuthProxy(proxyUrl, proxyUser, proxyPassword))
        default:
            panic(mperr.NewAliOss(errorcode.AliOssParam, "初始化类型不支持", nil))
        }

        conf2 := NewConfig()
        client, err := oss.New(conf2.GetEndpoint(), conf2.GetAccessKeyId(), conf2.GetAccessKeySecret(), clientOptions...)
        if err != nil {
            panic(mperr.NewAliOss(errorcode.AliOssParam, "创建客户端失败", nil))
        }
        insClient = client
    })

    return insClient
}
