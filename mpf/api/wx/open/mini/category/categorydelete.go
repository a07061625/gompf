/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 23:27
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

// 删除类目
type categoryDelete struct {
    wx.BaseWxOpen
    appId  string // 应用ID
    first  int    // 一级类目ID
    second int    // 二级类目ID
}

func (cd *categoryDelete) SetFirst(first int) {
    if first > 0 {
        cd.first = first
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "一级类目ID不合法", nil))
    }
}

func (cd *categoryDelete) SetSecond(second int) {
    if second > 0 {
        cd.second = second
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "二级类目ID不合法", nil))
    }
}

func (cd *categoryDelete) checkData() {
    if cd.first <= 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "一级类目ID不能为空", nil))
    }
    if cd.second <= 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "二级类目ID不能为空", nil))
    }
}

func (cd *categoryDelete) SendRequest() api.ApiResult {
    cd.checkData()

    reqData := make(map[string]interface{})
    reqData["first"] = cd.first
    reqData["second"] = cd.second
    reqBody := mpf.JsonMarshal(reqData)
    cd.ReqUrl = "https://api.weixin.qq.com/cgi-bin/wxopen/deletecategory?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(cd.appId)
    client, req := cd.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := cd.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewCategoryDelete(appId string) *categoryDelete {
    cd := &categoryDelete{wx.NewBaseWxOpen(), "", 0, 0}
    cd.appId = appId
    cd.ReqContentType = project.HTTPContentTypeJSON
    cd.ReqMethod = fasthttp.MethodPost
    return cd
}
