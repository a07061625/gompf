/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 21:45
 */
package life

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取关注者列表
type followQueryBatch struct {
    alipay.BaseAliPay
}

func (fqb *followQueryBatch) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, userId)
    if match {
        fqb.BizContent["next_user_id"] = userId
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "用户ID不合法", nil))
    }
}

func (fqb *followQueryBatch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    return fqb.GetRequest()
}

func NewFollowQueryBatch(appId string) *followQueryBatch {
    fqb := &followQueryBatch{alipay.NewBase(appId)}
    fqb.SetMethod("alipay.open.public.follow.batchquery")
    return fqb
}
