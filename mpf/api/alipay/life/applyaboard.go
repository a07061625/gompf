/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 21:50
 */
package life

import (
    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/valyala/fasthttp"
)

// 上架生活号
type applyAboard struct {
    alipay.BaseAliPay
}

func (aa *applyAboard) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    return aa.GetRequest()
}

func NewApplyAboard(appId string) *applyAboard {
    aa := &applyAboard{alipay.NewBase(appId)}
    aa.SetMethod("alipay.open.public.life.aboard.apply")
    return aa
}
