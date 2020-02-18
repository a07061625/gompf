/**
 * 短信基础工具
 * User: 姜伟
 * Date: 2019/12/23 0023
 * Time: 10:25
 */
package mpsms

import (
    "github.com/a07061625/gompf/mpf/api"
)

type utilSms struct {
    api.UtilAPI
}

var (
    insUtil *utilSms
)

func init() {
    insUtil = &utilSms{api.NewUtilAPI()}
}

func NewUtilSms() *utilSms {
    return insUtil
}
