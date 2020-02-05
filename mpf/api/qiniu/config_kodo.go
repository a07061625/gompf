/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/23 0023
 * Time: 9:55
 */
package qiniu

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type configKodo struct {
    accessKey string // 访问账号
    secretKey string // 密钥
}

func (c *configKodo) SetAccessKey(accessKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, accessKey)
    if match {
        c.accessKey = accessKey
    } else {
        panic(mperr.NewQiNiuKodo(errorcode.QiNiuKodoParam, "访问账号不合法", nil))
    }
}

func (c *configKodo) GetAccessKey() string {
    return c.accessKey
}

func (c *configKodo) SetSecretKey(secretKey string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, secretKey)
    if match {
        c.secretKey = secretKey
    } else {
        panic(mperr.NewQiNiuKodo(errorcode.QiNiuKodoParam, "密钥不合法", nil))
    }
}

func (c *configKodo) GetSecretKey() string {
    return c.secretKey
}
