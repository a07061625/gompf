/**
 * 短信配置
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 12:38
 */
package mpsms

import (
    "regexp"
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configYun253 struct {
    appKey     string // app key
    appSecret  string // app密钥
    urlSmsSend string // 短信下发链接
}

func (c *configYun253) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewSmsYun253(errorcode.SmsYun253Param, "app key不合法", nil))
    }
}

func (c *configYun253) GetAppKey() string {
    return c.appKey
}

func (c *configYun253) SetAppSecret(appSecret string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appSecret)
    if match {
        c.appSecret = appSecret
    } else {
        panic(mperr.NewSmsYun253(errorcode.SmsYun253Param, "app密钥不合法", nil))
    }
}

func (c *configYun253) GetAppSecret() string {
    return c.appSecret
}

func (c *configYun253) SetUrlSmsSend(urlSmsSend string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, urlSmsSend)
    if match {
        c.urlSmsSend = urlSmsSend
    } else {
        panic(mperr.NewSmsYun253(errorcode.SmsYun253Param, "短信下发链接不合法", nil))
    }
}

func (c *configYun253) GetUrlSmsSend() string {
    return c.urlSmsSend
}

var (
    onceConfigYun253 sync.Once
    insConfigYun253  *configYun253
)

func init() {
    insConfigYun253 = &configYun253{"", "", ""}
}

func NewConfigYun253() *configYun253 {
    onceConfigYun253.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpsms")
        insConfigYun253.SetAppKey(conf.GetString("yun253." + mpf.EnvProjectKey() + ".app.key"))
        insConfigYun253.SetAppSecret(conf.GetString("yun253." + mpf.EnvProjectKey() + ".app.secret"))
        insConfigYun253.SetUrlSmsSend(conf.GetString("yun253." + mpf.EnvProjectKey() + ".url.smssend"))

    })
    return insConfigYun253
}
