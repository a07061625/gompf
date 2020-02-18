/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 22:32
 */
package category

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type propertyGet struct {
    wx.BaseWxAccount
    appId      string
    categoryId int // 分类ID
}

func (pg *propertyGet) SetCategoryId(categoryId int) {
    if categoryId > 0 {
        pg.categoryId = categoryId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分类ID不合法", nil))
    }
}

func (pg *propertyGet) checkData() {
    if pg.categoryId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分类ID不能为空", nil))
    }
}

func (pg *propertyGet) SendRequest() api.ApiResult {
    pg.checkData()

    reqData := make(map[string]interface{})
    reqData["cate_id"] = pg.categoryId
    reqBody := mpf.JsonMarshal(reqData)
    pg.ReqUrl = "https://api.weixin.qq.com/merchant/category/getproperty?access_token=" + wx.NewUtilWx().GetSingleAccessToken(pg.appId)
    client, req := pg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := pg.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewPropertyGet(appId string) *propertyGet {
    pg := &propertyGet{wx.NewBaseWxAccount(), "", 0}
    pg.appId = appId
    pg.ReqContentType = project.HTTPContentTypeJSON
    pg.ReqMethod = fasthttp.MethodPost
    return pg
}
