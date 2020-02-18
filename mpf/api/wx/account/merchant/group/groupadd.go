/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 16:10
 */
package group

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type groupAdd struct {
    wx.BaseWxAccount
    appId       string
    groupName   string   // 分组名称
    productList []string // 商品ID列表
}

func (ga *groupAdd) SetGroupName(groupName string) {
    if len(groupName) > 0 {
        ga.groupName = groupName
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组名称不合法", nil))
    }
}

func (ga *groupAdd) SetProductList(productList []string) {
    ga.productList = make([]string, 0)
    for _, v := range productList {
        if len(v) > 0 {
            ga.productList = append(ga.productList, v)
        }
    }
}

func (ga *groupAdd) checkData() {
    if len(ga.groupName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分组名称不能为空", nil))
    }
    if len(ga.productList) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID列表不能为空", nil))
    }
}

func (ga *groupAdd) SendRequest() api.ApiResult {
    ga.checkData()

    groupInfo := make(map[string]interface{})
    groupInfo["group_name"] = ga.groupName
    groupInfo["product_list"] = ga.productList
    reqData := make(map[string]interface{})
    reqData["group_detail"] = groupInfo
    reqBody := mpf.JSONMarshal(reqData)
    ga.ReqUrl = "https://api.weixin.qq.com/merchant/group/add?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ga.appId)
    client, req := ga.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ga.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewGroupAdd(appId string) *groupAdd {
    ga := &groupAdd{wx.NewBaseWxAccount(), "", "", make([]string, 0)}
    ga.appId = appId
    ga.ReqContentType = project.HTTPContentTypeJSON
    ga.ReqMethod = fasthttp.MethodPost
    return ga
}
