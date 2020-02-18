/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 17:32
 */
package account

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 小程序名称设置及改名
type nicknameSet struct {
    wx.BaseWxOpen
    appId          string            // 应用ID
    nickname       string            // 昵称
    idCard         string            // 身份证照片
    license        string            // 营业执照
    otherStuffList map[string]string // 其他证明材料列表
}

func (ns *nicknameSet) SetNickname(nickname string) {
    if len(nickname) > 0 {
        ns.nickname = nickname
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "昵称不合法", nil))
    }
}

func (ns *nicknameSet) SetIdCard(idCard string) {
    if len(idCard) > 0 {
        ns.idCard = idCard
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "身份证照片不合法", nil))
    }
}

func (ns *nicknameSet) SetLicense(license string) {
    if len(license) > 0 {
        ns.license = license
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "营业执照不合法", nil))
    }
}

func (ns *nicknameSet) SetOtherStuffList(otherStuffList []string) {
    if len(otherStuffList) > 5 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "其他证明材料最多五个", nil))
    }

    num := 1
    for _, v := range otherStuffList {
        key := "naming_other_stuff_" + strconv.Itoa(num)
        ns.ReqData[key] = v
        num++
    }
}

func (ns *nicknameSet) checkData() {
    if len(ns.nickname) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "昵称不能为空", nil))
    }
    if (len(ns.idCard) == 0) && (len(ns.license) == 0) {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "身份证照片和营业执照至少要填一个", nil))
    }
    ns.ReqData["nick_name"] = ns.nickname
    if len(ns.idCard) > 0 {
        ns.ReqData["id_card"] = ns.idCard
    }
    if len(ns.license) > 0 {
        ns.ReqData["license"] = ns.license
    }
}

func (ns *nicknameSet) SendRequest() api.ApiResult {
    ns.checkData()

    reqBody := mpf.JsonMarshal(ns.ReqData)
    ns.ReqUrl = "https://api.weixin.qq.com/wxa/setnickname?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(ns.appId)
    client, req := ns.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ns.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewNicknameSet(appId string) *nicknameSet {
    ns := &nicknameSet{wx.NewBaseWxOpen(), "", "", "", "", make(map[string]string)}
    ns.appId = appId
    ns.ReqContentType = project.HTTPContentTypeJSON
    ns.ReqMethod = fasthttp.MethodPost
    return ns
}
