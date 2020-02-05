/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/23 0023
 * Time: 9:56
 */
package qiniu

import (
    "sync"

    "github.com/qiniu/api.v7/auth/qbox"
)

type utilQiNiu struct {
    kodoMac *qbox.Mac
}

func (util *utilQiNiu) GetKodoMac() *qbox.Mac {
    onceUtil.Do(func() {
        conf := NewConfig().GetKodo()
        util.kodoMac = qbox.NewMac(conf.GetAccessKey(), conf.GetSecretKey())
    })

    return util.kodoMac
}

var (
    onceUtil sync.Once
    insUtil  *utilQiNiu
)

func init() {
    insUtil = &utilQiNiu{}
}

func NewUtil() *utilQiNiu {
    return insUtil
}
