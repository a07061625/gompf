/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 11:28
 */
package invoice

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询电子发票
type infoGet struct {
    wx.BaseWxCorp
    corpId      string
    agentTag    string
    cardId      string // 发票id
    encryptCode string // 加密密码
}

func (ig *infoGet) SetCardId(cardId string) {
    if len(cardId) > 0 {
        ig.cardId = cardId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "发票id不合法", nil))
    }
}

func (ig *infoGet) SetEncryptCode(encryptCode string) {
    if len(encryptCode) > 0 {
        ig.encryptCode = encryptCode
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "加密密码不合法", nil))
    }
}

func (ig *infoGet) checkData() {
    if len(ig.cardId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "发票id不能为空", nil))
    }
    if len(ig.encryptCode) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "加密密码不能为空", nil))
    }
    ig.ReqData["card_id"] = ig.cardId
    ig.ReqData["encrypt_code"] = ig.encryptCode
}

func (ig *infoGet) SendRequest(getType string) api.ApiResult {
    ig.checkData()
    reqBody := mpf.JSONMarshal(ig.ReqData)

    ig.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/card/invoice/reimburse/getinvoiceinfo?access_token=" + wx.NewUtilWx().GetCorpCache(ig.corpId, ig.agentTag, getType)
    client, req := ig.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ig.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewInfoGet(corpId, agentTag string) *infoGet {
    ig := &infoGet{wx.NewBaseWxCorp(), "", "", "", ""}
    ig.corpId = corpId
    ig.agentTag = agentTag
    ig.ReqContentType = project.HTTPContentTypeJSON
    ig.ReqMethod = fasthttp.MethodPost
    return ig
}
