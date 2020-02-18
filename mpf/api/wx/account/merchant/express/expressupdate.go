/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 17:46
 */
package express

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type expressUpdate struct {
    wx.BaseWxAccount
    appId      string
    templateId string                   // 模板ID
    name       string                   // 模板名称
    payMode    int                      // 支付方式(0-买家承担运费 1-卖家承担运费)
    valuation  int                      // 计费单位(0-按件计费 1-按重量计费 2-按体积计费 目前只支持按件计费,默认为0)
    topFee     []map[string]interface{} // 运费计算列表
}

func (eu *expressUpdate) SetTemplateId(templateId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, templateId)
    if match {
        eu.templateId = templateId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板ID不合法", nil))
    }
}

func (eu *expressUpdate) SetName(name string) {
    if len(name) > 0 {
        eu.name = name
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板名称不合法", nil))
    }
}

func (eu *expressUpdate) SetPayMode(payMode int) {
    if (payMode == 0) || (payMode == 1) {
        eu.payMode = payMode
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "支付方式不合法", nil))
    }
}

func (eu *expressUpdate) SetTopFee(topFee []map[string]interface{}) {
    eu.topFee = topFee
}

func (eu *expressUpdate) checkData() {
    if len(eu.templateId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板ID不能为空", nil))
    }
    if len(eu.name) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板名称不能为空", nil))
    }
    if eu.payMode < 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "支付方式不能为空", nil))
    }
    if len(eu.topFee) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "运费计算列表不能为空", nil))
    }
}

func (eu *expressUpdate) SendRequest() api.ApiResult {
    eu.checkData()

    templateInfo := make(map[string]interface{})
    templateInfo["Valuation"] = eu.valuation
    templateInfo["Name"] = eu.name
    templateInfo["Assumer"] = eu.payMode
    templateInfo["TopFee"] = eu.topFee
    reqData := make(map[string]interface{})
    reqData["template_id"] = eu.templateId
    reqData["delivery_template"] = templateInfo
    reqBody := mpf.JsonMarshal(reqData)
    eu.ReqUrl = "https://api.weixin.qq.com/merchant/express/update?access_token=" + wx.NewUtilWx().GetSingleAccessToken(eu.appId)
    client, req := eu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := eu.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewExpressUpdate(appId string) *expressUpdate {
    eu := &expressUpdate{wx.NewBaseWxAccount(), "", "", "", 0, 0, make([]map[string]interface{}, 0)}
    eu.appId = appId
    eu.payMode = -1
    eu.valuation = 0
    eu.ReqContentType = project.HTTPContentTypeJSON
    eu.ReqMethod = fasthttp.MethodPost
    return eu
}
