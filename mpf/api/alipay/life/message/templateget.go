/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 21:24
 */
package message

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 消息模板领取接口
type templateGet struct {
    alipay.BaseAliPay
    templateId string // 模板ID
}

func (tg *templateGet) SetTemplateId(templateId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,20}$`, templateId)
    if match {
        tg.templateId = templateId
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "模板ID不合法", nil))
    }
}

func (tg *templateGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tg.templateId) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "模板ID不能为空", nil))
    }
    tg.BizContent["template_id"] = tg.templateId

    return tg.GetRequest()
}

func NewTemplateGet(appId string) *templateGet {
    tg := &templateGet{alipay.NewBase(appId), ""}
    tg.SetMethod("alipay.open.public.template.message.get")
    return tg
}
