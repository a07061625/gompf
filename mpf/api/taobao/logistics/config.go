/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/18 0018
 * Time: 22:43
 */
package logistics

import (
    "regexp"
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configTaoBao struct {
    appKey    string // APP KEY
    appSecret string // APP 密钥
}

func (c *configTaoBao) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "app key不合法", nil))
    }
}

func (c *configTaoBao) GetAppKey() string {
    return c.appKey
}

func (c *configTaoBao) SetAppSecret(appSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appSecret)
    if match {
        c.appSecret = appSecret
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "app 密钥不合法", nil))
    }
}

func (c *configTaoBao) GetAppSecret() string {
    return c.appSecret
}

var (
    onceConfigTaoBao sync.Once
    insConfigTaoBao  *configTaoBao
)

func init() {
    insConfigTaoBao = &configTaoBao{}
}

func NewConfigTaoBao() *configTaoBao {
    onceConfigTaoBao.Do(func() {
        conf := mpf.NewConfig().GetConfig("mplogistics")
        insConfigTaoBao.SetAppKey(conf.GetString("taobao." + mpf.EnvProjectKey() + ".app.key"))
        insConfigTaoBao.SetAppSecret(conf.GetString("taobao." + mpf.EnvProjectKey() + ".app.secret"))
    })

    return insConfigTaoBao
}
