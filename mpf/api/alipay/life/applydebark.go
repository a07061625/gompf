/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 21:52
 */
package life

import (
    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/valyala/fasthttp"
)

// 下架生活号
type applyDebark struct {
    alipay.BaseAliPay
}

func (ad *applyDebark) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    return ad.GetRequest()
}

func NewApplyDebark(appId string) *applyDebark {
    ad := &applyDebark{alipay.NewBase(appId)}
    ad.SetMethod("alipay.open.public.life.debark.apply")
    return ad
}
