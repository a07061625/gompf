/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 17:46
 */
package express

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type expressAdd struct {
    wx.BaseWxAccount
    appId     string
    name      string                   // 模板名称
    payMode   int                      // 支付方式(0-买家承担运费 1-卖家承担运费)
    valuation int                      // 计费单位(0-按件计费 1-按重量计费 2-按体积计费 目前只支持按件计费,默认为0)
    topFee    []map[string]interface{} // 运费计算列表
}

func (ea *expressAdd) SetName(name string) {
    if len(name) > 0 {
        ea.name = name
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板名称不合法", nil))
    }
}

func (ea *expressAdd) SetPayMode(payMode int) {
    if (payMode == 0) || (payMode == 1) {
        ea.payMode = payMode
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "支付方式不合法", nil))
    }
}

func (ea *expressAdd) SetTopFee(topFee []map[string]interface{}) {
    ea.topFee = topFee
}

func (ea *expressAdd) checkData() {
    if len(ea.name) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板名称不能为空", nil))
    }
    if ea.payMode < 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "支付方式不能为空", nil))
    }
    if len(ea.topFee) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "运费计算列表不能为空", nil))
    }
}

func (ea *expressAdd) SendRequest() api.ApiResult {
    ea.checkData()

    templateInfo := make(map[string]interface{})
    templateInfo["Valuation"] = ea.valuation
    templateInfo["Name"] = ea.name
    templateInfo["Assumer"] = ea.payMode
    templateInfo["TopFee"] = ea.topFee
    reqData := make(map[string]interface{})
    reqData["delivery_template"] = templateInfo
    reqBody := mpf.JSONMarshal(reqData)
    ea.ReqUrl = "https://api.weixin.qq.com/merchant/express/add?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ea.appId)
    client, req := ea.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ea.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewExpressAdd(appId string) *expressAdd {
    ea := &expressAdd{wx.NewBaseWxAccount(), "", "", 0, 0, make([]map[string]interface{}, 0)}
    ea.appId = appId
    ea.payMode = -1
    ea.valuation = 0
    ea.ReqContentType = project.HTTPContentTypeJSON
    ea.ReqMethod = fasthttp.MethodPost
    return ea
}
