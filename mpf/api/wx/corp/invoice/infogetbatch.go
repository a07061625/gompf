/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 12:22
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

// 批量查询电子发票
type infoGetBatch struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    itemList []map[string]string // 发票列表
}

func (igb *infoGetBatch) SetItemList(itemList []map[string]string) {
    if len(itemList) > 0 {
        igb.itemList = itemList
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "发票列表不合法", nil))
    }
}

func (igb *infoGetBatch) checkData() {
    if len(igb.itemList) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "发票列表不能为空", nil))
    }
}

func (igb *infoGetBatch) SendRequest(getType string) api.APIResult {
    igb.checkData()
    reqData := make(map[string]interface{})
    reqData["item_list"] = igb.itemList
    reqBody := mpf.JSONMarshal(reqData)

    igb.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/card/invoice/reimburse/getinvoiceinfobatch?access_token=" + wx.NewUtilWx().GetCorpCache(igb.corpId, igb.agentTag, getType)
    client, req := igb.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := igb.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewInfoGetBatch(corpId, agentTag string) *infoGetBatch {
    igb := &infoGetBatch{wx.NewBaseWxCorp(), "", "", make([]map[string]string, 0)}
    igb.corpId = corpId
    igb.agentTag = agentTag
    igb.ReqContentType = project.HTTPContentTypeJSON
    igb.ReqMethod = fasthttp.MethodPost
    return igb
}
