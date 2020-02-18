/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 17:10
 */
package shelf

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type shelfModify struct {
    wx.BaseWxAccount
    appId       string
    shelfId     int                      // 货架ID
    shelfName   string                   // 货架名称
    shelfBanner string                   // 货架招牌图片Url
    shelfData   []map[string]interface{} // 货架信息列表
}

func (sm *shelfModify) SetShelfId(shelfId int) {
    if shelfId > 0 {
        sm.shelfId = shelfId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架ID不合法", nil))
    }
}

func (sm *shelfModify) SetShelfName(shelfName string) {
    if len(shelfName) > 0 {
        sm.shelfName = shelfName
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架名称不合法", nil))
    }
}

func (sm *shelfModify) SetShelfBanner(shelfBanner string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, shelfBanner)
    if match {
        sm.shelfBanner = shelfBanner
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架招牌图片Url不合法", nil))
    }
}

func (sm *shelfModify) SetShelfData(shelfData []map[string]interface{}) {
    sm.shelfData = shelfData
}

func (sm *shelfModify) checkData() {
    if sm.shelfId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架ID不能为空", nil))
    }
    if len(sm.shelfName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架名称不能为空", nil))
    }
    if len(sm.shelfBanner) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架招牌图片Url不能为空", nil))
    }
    if len(sm.shelfData) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架信息列表不能为空", nil))
    }
}

func (sm *shelfModify) SendRequest() api.ApiResult {
    sm.checkData()

    moduleInfo := make(map[string]interface{})
    moduleInfo["module_infos"] = sm.shelfData
    reqData := make(map[string]interface{})
    reqData["shelf_id"] = sm.shelfId
    reqData["shelf_name"] = sm.shelfName
    reqData["shelf_banner"] = sm.shelfBanner
    reqData["shelf_data"] = moduleInfo
    reqBody := mpf.JSONMarshal(sm.ReqData)
    sm.ReqUrl = "https://api.weixin.qq.com/merchant/shelf/mod?access_token=" + wx.NewUtilWx().GetSingleAccessToken(sm.appId)
    client, req := sm.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sm.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewShelfModify(appId string) *shelfModify {
    sm := &shelfModify{wx.NewBaseWxAccount(), "", 0, "", "", make([]map[string]interface{}, 0)}
    sm.appId = appId
    sm.ReqContentType = project.HTTPContentTypeJSON
    sm.ReqMethod = fasthttp.MethodPost
    return sm
}
