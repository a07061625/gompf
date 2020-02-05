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

type shelfAdd struct {
    wx.BaseWxAccount
    appId       string
    shelfName   string                   // 货架名称
    shelfBanner string                   // 货架招牌图片Url
    shelfData   []map[string]interface{} // 货架信息列表
}

func (sa *shelfAdd) SetShelfName(shelfName string) {
    if len(shelfName) > 0 {
        sa.shelfName = shelfName
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架名称不合法", nil))
    }
}

func (sa *shelfAdd) SetShelfBanner(shelfBanner string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, shelfBanner)
    if match {
        sa.shelfBanner = shelfBanner
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架招牌图片Url不合法", nil))
    }
}

func (sa *shelfAdd) SetShelfData(shelfData []map[string]interface{}) {
    sa.shelfData = shelfData
}

func (sa *shelfAdd) checkData() {
    if len(sa.shelfName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架名称不能为空", nil))
    }
    if len(sa.shelfBanner) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架招牌图片Url不能为空", nil))
    }
    if len(sa.shelfData) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架信息列表不能为空", nil))
    }
}

func (sa *shelfAdd) SendRequest() api.ApiResult {
    sa.checkData()

    moduleInfo := make(map[string]interface{})
    moduleInfo["module_infos"] = sa.shelfData
    reqData := make(map[string]interface{})
    reqData["shelf_name"] = sa.shelfName
    reqData["shelf_banner"] = sa.shelfBanner
    reqData["shelf_data"] = moduleInfo
    reqBody := mpf.JsonMarshal(sa.ReqData)
    sa.ReqUrl = "https://api.weixin.qq.com/merchant/shelf/add?access_token=" + wx.NewUtilWx().GetSingleAccessToken(sa.appId)
    client, req := sa.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sa.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewShelfAdd(appId string) *shelfAdd {
    sa := &shelfAdd{wx.NewBaseWxAccount(), "", "", "", make([]map[string]interface{}, 0)}
    sa.appId = appId
    sa.ReqContentType = project.HttpContentTypeJson
    sa.ReqMethod = fasthttp.MethodPost
    return sa
}
