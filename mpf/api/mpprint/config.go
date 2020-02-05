/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 19:46
 */
package mpprint

import (
    "regexp"
    "strings"
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configFeYin struct {
    appId      string // 应用id
    appKey     string // API密钥
    memberCode string // 商户编码
}

func (c *configFeYin) SetAppId(appId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appId)
    if match {
        c.appId = appId
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "应用id不合法", nil))
    }
}

func (c *configFeYin) GetAppId() string {
    return c.appId
}

func (c *configFeYin) SetAppKey(appKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appKey)
    if match {
        c.appKey = appKey
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "API密钥不合法", nil))
    }
}

func (c *configFeYin) GetAppKey() string {
    return c.appKey
}

func (c *configFeYin) SetMemberCode(memberCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, memberCode)
    if match {
        c.memberCode = memberCode
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "商户编码不合法", nil))
    }
}

func (c *configFeYin) GetMemberCode() string {
    return c.memberCode
}

var (
    onceConfigFeYin sync.Once
    insConfigsFeYin map[string]*configFeYin
)

func init() {
    insConfigsFeYin = make(map[string]*configFeYin)
}

func NewConfigFeYin(appId string) *configFeYin {
    onceConfigFeYin.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpprint")
        configs := conf.GetStringSlice("feyin." + mpf.EnvProjectKey() + ".list")
        for _, c := range configs {
            cs := strings.Split(c, "_")
            econf := &configFeYin{}
            econf.SetAppId(cs[0])
            econf.SetAppKey(cs[1])
            econf.SetMemberCode(cs[2])
            insConfigsFeYin[cs[0]] = econf
        }
    })

    fy, ok := insConfigsFeYin[appId]
    if ok {
        return fy
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "飞印配置不存在", nil))
    }
}
