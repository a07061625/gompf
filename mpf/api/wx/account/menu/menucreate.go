/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 12:12
 */
package menu

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建菜单
type menuCreate struct {
    wx.BaseWxAccount
    appId    string
    menuList []map[string]interface{} // 菜单列表
}

func (mc *menuCreate) AddMenu(menu *menu) {
    if len(mc.menuList) >= 3 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "菜单数量不能超过3个", nil))
    }
    mc.menuList = append(mc.menuList, menu.GetResult())
}

func (mc *menuCreate) checkData() {
    if len(mc.menuList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "菜单列表不能为空", nil))
    }
}

func (mc *menuCreate) SendRequest() api.ApiResult {
    mc.checkData()

    reqData := make(map[string]interface{})
    reqData["button"] = mc.menuList
    reqBody := mpf.JSONMarshal(reqData)
    mc.ReqUrl = "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=" + wx.NewUtilWx().GetSingleAccessToken(mc.appId)
    client, req := mc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := mc.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewMenuCreate(appId string) *menuCreate {
    mc := &menuCreate{wx.NewBaseWxAccount(), "", make([]map[string]interface{}, 0)}
    mc.appId = appId
    mc.ReqContentType = project.HTTPContentTypeJSON
    mc.ReqMethod = fasthttp.MethodPost
    return mc
}
