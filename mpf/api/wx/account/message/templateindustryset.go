/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 18:37
 */
package message

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 设置模板消息所属行业
type templateIndustrySet struct {
    wx.BaseWxAccount
    appId       string
    industryId1 int // 行业编号1
    industryId2 int // 行业编号2
}

func (tis *templateIndustrySet) SetIndustryId1(industryId1 int) {
    if industryId1 > 0 {
        tis.industryId1 = industryId1
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "行业编号1不合法", nil))
    }
}

func (tis *templateIndustrySet) SetIndustryId2(industryId2 int) {
    if industryId2 > 0 {
        tis.industryId2 = industryId2
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "行业编号2不合法", nil))
    }
}

func (tis *templateIndustrySet) checkData() {
    if tis.industryId1 <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "行业编号1不能为空", nil))
    }
    if tis.industryId2 <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "行业编号2不能为空", nil))
    }
}

func (tis *templateIndustrySet) SendRequest() api.ApiResult {
    tis.checkData()

    reqData := make(map[string]interface{})
    reqData["industry_id1"] = tis.industryId1
    reqData["industry_id2"] = tis.industryId2
    reqBody := mpf.JsonMarshal(reqData)
    tis.ReqUrl = "https://api.weixin.qq.com/cgi-bin/template/api_set_industry?access_token=" + wx.NewUtilWx().GetSingleAccessToken(tis.appId)
    client, req := tis.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := tis.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewTemplateIndustrySet(appId string) *templateIndustrySet {
    tis := &templateIndustrySet{wx.NewBaseWxAccount(), "", 0, 0}
    tis.appId = appId
    tis.ReqContentType = project.HTTPContentTypeJSON
    tis.ReqMethod = fasthttp.MethodPost
    return tis
}
