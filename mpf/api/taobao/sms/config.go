/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/18 0018
 * Time: 22:15
 */
package sms

import (
    "regexp"
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configDaYu struct {
    appKey    string // app key
    appSecret string // app密钥
}

func (c *configDaYu) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "app key不合法", nil))
    }
}

func (c *configDaYu) GetAppKey() string {
    return c.appKey
}

func (c *configDaYu) SetAppSecret(appSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appSecret)
    if match {
        c.appSecret = appSecret
    } else {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "app密钥不合法", nil))
    }
}

func (c *configDaYu) GetAppSecret() string {
    return c.appSecret
}

var (
    onceConfigDaYu sync.Once
    insConfigDaYu  *configDaYu
)

func init() {
    insConfigDaYu = &configDaYu{"", ""}
}

func NewConfigDaYu() *configDaYu {
    onceConfigDaYu.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpsms")
        insConfigDaYu.SetAppKey(conf.GetString("dayu." + mpf.EnvProjectKey() + ".app.key"))
        insConfigDaYu.SetAppSecret(conf.GetString("dayu." + mpf.EnvProjectKey() + ".app.secret"))
    })
    return insConfigDaYu
}
