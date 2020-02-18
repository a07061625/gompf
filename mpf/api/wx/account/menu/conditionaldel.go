/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 12:30
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

// 删除个性化菜单
type conditionalDel struct {
    wx.BaseWxAccount
    appId  string
    menuId string // 菜单ID
}

func (cd *conditionalDel) SetMenuId(menuId string) {
    if len(menuId) > 0 {
        cd.menuId = menuId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "菜单ID不合法", nil))
    }
}

func (cd *conditionalDel) checkData() {
    if len(cd.menuId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "菜单ID不能为空", nil))
    }
    cd.ReqData["menuid"] = cd.menuId
}

func (cd *conditionalDel) SendRequest() api.ApiResult {
    cd.checkData()

    reqBody := mpf.JsonMarshal(cd.ReqData)
    cd.ReqUrl = "https://api.weixin.qq.com/cgi-bin/menu/delconditional?access_token=" + wx.NewUtilWx().GetSingleAccessToken(cd.appId)
    client, req := cd.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := cd.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewConditionalDel(appId string) *conditionalDel {
    cd := &conditionalDel{wx.NewBaseWxAccount(), "", ""}
    cd.appId = appId
    cd.ReqContentType = project.HTTPContentTypeJSON
    cd.ReqMethod = fasthttp.MethodPost
    return cd
}
