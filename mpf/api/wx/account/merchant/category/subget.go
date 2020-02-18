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

type subGet struct {
    wx.BaseWxAccount
    appId      string
    categoryId int // 分类ID
}

func (sg *subGet) SetCategoryId(categoryId int) {
    if categoryId > 0 {
        sg.categoryId = categoryId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分类ID不合法", nil))
    }
}

func (sg *subGet) checkData() {
    if sg.categoryId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分类ID不能为空", nil))
    }
}

func (sg *subGet) SendRequest() api.ApiResult {
    sg.checkData()

    reqData := make(map[string]interface{})
    reqData["cate_id"] = sg.categoryId
    reqBody := mpf.JsonMarshal(reqData)
    sg.ReqUrl = "https://api.weixin.qq.com/merchant/category/getsub?access_token=" + wx.NewUtilWx().GetSingleAccessToken(sg.appId)
    client, req := sg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sg.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewSubGet(appId string) *subGet {
    sg := &subGet{wx.NewBaseWxAccount(), "", 0}
    sg.appId = appId
    sg.ReqContentType = project.HTTPContentTypeJSON
    sg.ReqMethod = fasthttp.MethodPost
    return sg
}
