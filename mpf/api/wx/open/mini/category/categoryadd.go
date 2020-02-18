/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 23:00
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

// 添加类目
type categoryAdd struct {
    wx.BaseWxOpen
    appId      string                   // 应用ID
    categories []map[string]interface{} // 类目信息列表
}

func (ca *categoryAdd) SetCategories(categories []map[string]interface{}) {
    ca.categories = categories
}

func (ca *categoryAdd) checkData() {
    if len(ca.categories) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "类目信息不能为空", nil))
    }
}

func (ca *categoryAdd) SendRequest() api.ApiResult {
    ca.checkData()

    reqData := make(map[string]interface{})
    reqData["categories"] = ca.categories
    reqBody := mpf.JSONMarshal(reqData)
    ca.ReqUrl = "https://api.weixin.qq.com/cgi-bin/wxopen/addcategory?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(ca.appId)
    client, req := ca.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ca.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewCategoryAdd(appId string) *categoryAdd {
    ca := &categoryAdd{wx.NewBaseWxOpen(), "", make([]map[string]interface{}, 0)}
    ca.appId = appId
    ca.ReqContentType = project.HTTPContentTypeJSON
    ca.ReqMethod = fasthttp.MethodPost
    return ca
}
