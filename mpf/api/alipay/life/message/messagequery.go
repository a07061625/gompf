/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 23:24
 */
package message

import (
    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 生活号查询已发送消息接口
type messageQuery struct {
    alipay.BaseAliPay
    idList []string // 消息id列表
}

func (mq *messageQuery) SetIdList(idList []string) {
    if len(idList) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "消息ID列表不合法", nil))
    } else if len(idList) > 20 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "消息ID列表超过限制", nil))
    }

    mq.idList = make([]string, 0)
    for _, v := range idList {
        if (len(v) > 0) && (len(v) <= 64) {
            mq.idList = append(mq.idList, v)
        }
    }
}

func (mq *messageQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(mq.idList) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "消息ID列表不能为空", nil))
    }
    mq.BizContent["message_ids"] = mq.idList

    return mq.GetRequest()
}

func NewMessageQuery(appId string) *messageQuery {
    mq := &messageQuery{alipay.NewBase(appId), make([]string, 0)}
    mq.SetMethod("alipay.open.public.message.query")
    return mq
}
