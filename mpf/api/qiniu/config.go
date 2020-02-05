/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/23 0023
 * Time: 9:55
 */
package qiniu

import (
    "sync"

    "github.com/a07061625/gompf/mpf"
)

type configQiNiu struct {
    kodo *configKodo
}

func (c *configQiNiu) GetKodo() *configKodo {
    onceConfigKodo.Do(func() {
        conf := mpf.NewConfig().GetConfig("qiniu")
        c.kodo.SetAccessKey(conf.GetString("kodo." + mpf.EnvProjectKey() + ".access.key"))
        c.kodo.SetSecretKey(conf.GetString("kodo." + mpf.EnvProjectKey() + ".secret.key"))
    })

    return c.kodo
}

var (
    onceConfigKodo sync.Once
    insConfig      *configQiNiu
)

func init() {
    insConfig = &configQiNiu{}
    insConfig.kodo = &configKodo{}
}

func NewConfig() *configQiNiu {
    return insConfig
}
