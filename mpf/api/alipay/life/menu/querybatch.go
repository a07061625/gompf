/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 21:37
 */
package menu

import (
    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/valyala/fasthttp"
)

// 菜单列表查询接口
type queryBatch struct {
    alipay.BaseAliPay
}

func (qb *queryBatch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    return qb.GetRequest()
}

func NewQueryBatch(appId string) *queryBatch {
    qb := &queryBatch{alipay.NewBase(appId)}
    qb.SetMethod("alipay.open.public.menu.batchquery")
    return qb
}
