/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 9:58
 */
package mpim

import (
    "sync"

    "github.com/a07061625/gompf/mpf/api"
)

type ICache interface {
    GetTencentAccountSign(account string) string
    DelTencentAccountSign(account string)
}

type utilIM struct {
    api.UtilApi
    cache ICache
}

var (
    onceUtil sync.Once
    insUtil  *utilIM
)

func init() {
    insUtil = &utilIM{api.NewUtilApi(), nil}
}

func LoadUtil(cache ICache) {
    onceUtil.Do(func() {
        insUtil.cache = cache
    })
}

func NewUtilIM() *utilIM {
    return insUtil
}
