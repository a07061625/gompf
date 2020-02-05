/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 16:54
 */
package aliopen

import (
    "sync"

    "github.com/a07061625/gompf/mpf"
)

type configAliOpen struct {
    dySms *configDySms
    push  *configPush
    iot   *configIot
}

func (config *configAliOpen) GetDySms() *configDySms {
    onceConfigDySms.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpsms")
        config.dySms.SetRegionId(conf.GetString("aliyun." + mpf.EnvProjectKey() + ".region.id"))
        config.dySms.SetAppKey(conf.GetString("aliyun." + mpf.EnvProjectKey() + ".app.key"))
        config.dySms.SetAppSecret(conf.GetString("aliyun." + mpf.EnvProjectKey() + ".app.secret"))
    })

    return config.dySms
}

func (config *configAliOpen) GetPush() *configPush {
    onceConfigPush.Do(func() {
        conf := mpf.NewConfig().GetConfig("mppush")
        config.push.SetRegionId(conf.GetString("aliyun." + mpf.EnvProjectKey() + ".region.id"))
        config.push.SetAppKey(conf.GetString("aliyun." + mpf.EnvProjectKey() + ".app.key"))
        config.push.SetAppSecret(conf.GetString("aliyun." + mpf.EnvProjectKey() + ".app.secret"))
    })

    return config.push
}

func (config *configAliOpen) GetIot() *configIot {
    onceConfigIot.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpiot")
        config.iot.SetRegionId(conf.GetString("aliyun." + mpf.EnvProjectKey() + ".region.id"))
        config.iot.SetAppKey(conf.GetString("aliyun." + mpf.EnvProjectKey() + ".app.key"))
        config.iot.SetAppSecret(conf.GetString("aliyun." + mpf.EnvProjectKey() + ".app.secret"))
    })

    return config.iot
}

var (
    onceConfigDySms sync.Once
    onceConfigPush  sync.Once
    onceConfigIot   sync.Once
    insConfig       *configAliOpen
)

func init() {
    insConfig = &configAliOpen{}
    insConfig.dySms = &configDySms{}
    insConfig.push = &configPush{}
    insConfig.iot = &configIot{}
}

func NewConfig() *configAliOpen {
    return insConfig
}
