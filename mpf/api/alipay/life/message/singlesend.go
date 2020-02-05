/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 21:32
 */
package message

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 单发模板消息
type singleSend struct {
    alipay.BaseAliPay
    userId       string                 // 用户ID
    templateInfo map[string]interface{} // 模板信息
}

func (ss *singleSend) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, userId)
    if match {
        ss.userId = userId
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "用户ID不合法", nil))
    }
}

func (ss *singleSend) SetTemplateInfo(templateInfo map[string]interface{}) {
    if len(templateInfo) > 0 {
        ss.templateInfo = templateInfo
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "模板信息不合法", nil))
    }
}

func (ss *singleSend) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ss.userId) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "用户ID不能为空", nil))
    }
    if len(ss.templateInfo) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "模板信息不能为空", nil))
    }
    ss.BizContent["to_user_id"] = ss.userId
    ss.BizContent["template"] = ss.templateInfo

    return ss.GetRequest()
}

func NewSingleSend(appId string) *singleSend {
    ss := &singleSend{alipay.NewBase(appId), "", make(map[string]interface{})}
    ss.SetMethod("alipay.open.public.message.single.send")
    return ss
}
