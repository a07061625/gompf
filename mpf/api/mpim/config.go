/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 9:58
 */
package mpim

import (
    "os"
    "regexp"
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configTencent struct {
    appId          string // 应用id
    accountAdmin   string // 管理员帐号
    accountType    string // 账号类型
    filePrivateKey string // 私钥文件,全路径
    fileCommand    string // 签名命令文件,全路径
}

func (c *configTencent) SetAppId(appId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, appId)
    if match {
        c.appId = appId
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "应用id不合法", nil))
    }
}

func (c *configTencent) GetAppId() string {
    return c.appId
}

func (c *configTencent) SetAccountAdmin(accountAdmin string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, accountAdmin)
    if match {
        c.accountAdmin = accountAdmin
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "管理员帐号不合法", nil))
    }
}

func (c *configTencent) GetAccountAdmin() string {
    return c.accountAdmin
}

func (c *configTencent) SetAccountType(accountType string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, accountType)
    if match {
        c.accountType = accountType
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "账号类型不合法", nil))
    }
}

func (c *configTencent) GetAccountType() string {
    return c.accountType
}

func (c *configTencent) SetFilePrivateKey(filePrivateKey string) {
    f, err := os.Stat(filePrivateKey)
    if err != nil {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "私钥文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "私钥文件不能是目录", nil))
    }
    c.filePrivateKey = filePrivateKey
}

func (c *configTencent) GetFilePrivateKey() string {
    return c.filePrivateKey
}

func (c *configTencent) SetFileCommand(fileCommand string) {
    f, err := os.Stat(fileCommand)
    if err != nil {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "命令文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "命令文件不能是目录", nil))
    }
    c.fileCommand = fileCommand
}

func (c *configTencent) GetFileCommand() string {
    return c.fileCommand
}

var (
    onceConfigTencent sync.Once
    insConfigTencent  *configTencent
)

func init() {
    insConfigTencent = &configTencent{}
}

func NewConfigTencent() *configTencent {
    onceConfigTencent.Do(func() {
        conf := mpf.NewConfig().GetConfig("mpim")
        insConfigTencent.SetAppId(conf.GetString("tencent." + mpf.EnvProjectKey() + ".app.id"))
        insConfigTencent.SetAccountAdmin(conf.GetString("tencent." + mpf.EnvProjectKey() + ".account.admin"))
        insConfigTencent.SetAccountType(conf.GetString("tencent." + mpf.EnvProjectKey() + ".account.type"))
        insConfigTencent.SetFilePrivateKey(conf.GetString("tencent." + mpf.EnvProjectKey() + ".file.privatekey"))
        insConfigTencent.SetFileCommand(conf.GetString("tencent." + mpf.EnvProjectKey() + ".file.command"))
    })
    return insConfigTencent
}
