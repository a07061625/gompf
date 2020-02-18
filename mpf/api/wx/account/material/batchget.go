/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 9:45
 */
package material

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/api/wx/account"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取素材列表
type batchGet struct {
    wx.BaseWxAccount
    appId        string
    materialType string // 素材类型
    offset       int    // 偏移位置
    count        int    // 条数
}

func (bg *batchGet) SetMaterialType(materialType string) {
    _, ok := account.MaterialTypes[materialType]
    if ok {
        bg.materialType = materialType
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "素材类型不合法", nil))
    }
}

func (bg *batchGet) SetOffset(offset int) {
    if offset >= 0 {
        bg.offset = offset
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "偏移位置不合法", nil))
    }
}

func (bg *batchGet) SetCount(count int) {
    if (count > 0) && (count <= 20) {
        bg.count = count
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "条数不合法", nil))
    }
}

func (bg *batchGet) checkData() {
    if len(bg.materialType) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "素材类型不能为空", nil))
    }
}

func (bg *batchGet) SendRequest() api.APIResult {
    bg.checkData()

    reqData := make(map[string]interface{})
    reqData["type"] = bg.materialType
    reqData["offset"] = bg.offset
    reqData["count"] = bg.count
    reqBody := mpf.JSONMarshal(reqData)
    bg.ReqURI = "https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=" + wx.NewUtilWx().GetSingleAccessToken(bg.appId)
    client, req := bg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := bg.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewBatchGet(appId string) *batchGet {
    bg := &batchGet{wx.NewBaseWxAccount(), "", "", 0, 0}
    bg.appId = appId
    bg.offset = 0
    bg.count = 20
    bg.ReqContentType = project.HTTPContentTypeJSON
    bg.ReqMethod = fasthttp.MethodPost
    return bg
}
