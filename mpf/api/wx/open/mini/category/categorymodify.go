/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 23:40
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

// 修改类目
type categoryModify struct {
    wx.BaseWxOpen
    appId      string                   // 应用ID
    first      int                      // 一级类目ID
    second     int                      // 二级类目ID
    categories []map[string]interface{} // 类目信息列表
}

func (cm *categoryModify) SetFirst(first int) {
    if first > 0 {
        cm.first = first
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "一级类目ID不合法", nil))
    }
}

func (cm *categoryModify) SetSecond(second int) {
    if second > 0 {
        cm.second = second
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "二级类目ID不合法", nil))
    }
}

func (cm *categoryModify) SetCategories(categories []map[string]interface{}) {
    cm.categories = categories
}

func (cm *categoryModify) checkData() {
    if cm.first <= 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "一级类目ID不能为空", nil))
    }
    if cm.second <= 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "二级类目ID不能为空", nil))
    }
    if len(cm.categories) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "类目信息不能为空", nil))
    }
}

func (cm *categoryModify) SendRequest() api.APIResult {
    cm.checkData()

    reqData := make(map[string]interface{})
    reqData["first"] = cm.first
    reqData["second"] = cm.second
    reqData["categories"] = cm.categories
    reqBody := mpf.JSONMarshal(cm.ReqData)
    cm.ReqURI = "https://api.weixin.qq.com/cgi-bin/wxopen/modifycategory?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(cm.appId)
    client, req := cm.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := cm.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewCategoryModify(appId string) *categoryModify {
    cm := &categoryModify{wx.NewBaseWxOpen(), "", 0, 0, make([]map[string]interface{}, 0)}
    cm.appId = appId
    cm.ReqContentType = project.HTTPContentTypeJSON
    cm.ReqMethod = fasthttp.MethodPost
    return cm
}
