/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 17:10
 */
package aliopen

import (
    "sync"

    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
    "github.com/aliyun/alibaba-cloud-sdk-go/services/iot"
    "github.com/aliyun/alibaba-cloud-sdk-go/services/push"
)

type utilAliOpen struct {
    clientDySms *dysmsapi.Client
    clientPush  *push.Client
    clientIot   *iot.Client
}

func (util *utilAliOpen) GetClientDySms() *dysmsapi.Client {
    onceClientDySms.Do(func() {
        conf := NewConfig().GetDySms()
        cli, err := dysmsapi.NewClientWithAccessKey(conf.GetRegionId(), conf.GetAppKey(), conf.GetAppSecret())
        if err != nil {
            panic(mperr.NewSmsAliYun(errorcode.SmsAliYunParam, "创建请求客户端出错", err))
        }

        util.clientDySms = cli
    })
    return util.clientDySms
}

func (util *utilAliOpen) GetClientPush() *push.Client {
    onceClientPush.Do(func() {
        conf := NewConfig().GetPush()
        cli, err := push.NewClientWithAccessKey(conf.GetRegionId(), conf.GetAppKey(), conf.GetAppSecret())
        if err != nil {
            panic(mperr.NewPushAliYun(errorcode.PushAliYunParam, "创建请求客户端出错", err))
        }

        util.clientPush = cli
    })
    return util.clientPush
}

func (util *utilAliOpen) GetClientIot() *iot.Client {
    onceClientIot.Do(func() {
        conf := NewConfig().GetIot()
        cli, err := iot.NewClientWithAccessKey(conf.GetRegionId(), conf.GetAppKey(), conf.GetAppSecret())
        if err != nil {
            panic(mperr.NewIotAliYun(errorcode.IotAliYunParam, "创建请求客户端出错", err))
        }

        util.clientIot = cli
    })
    return util.clientIot
}

var (
    onceClientDySms sync.Once
    onceClientPush  sync.Once
    onceClientIot   sync.Once
    insUtil         *utilAliOpen
)

func init() {
    insUtil = &utilAliOpen{}
}

func NewUtil() *utilAliOpen {
    return insUtil
}
