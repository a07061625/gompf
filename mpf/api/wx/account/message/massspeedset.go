/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 16:16
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

// 控制群发速度
type massSpeedSet struct {
    wx.BaseWxAccount
    appId string
    speed int // 群发速度级别
}

func (mss *massSpeedSet) SetSpeed(speed int) {
    if (speed >= 0) && (speed <= 4) {
        mss.speed = speed
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "群发速度级别不合法", nil))
    }
}

func (mss *massSpeedSet) checkData() {
    if mss.speed < 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "群发速度级别不能为空", nil))
    }
}

func (mss *massSpeedSet) SendRequest() api.ApiResult {
    mss.checkData()

    reqData := make(map[string]interface{})
    reqData["speed"] = mss.speed
    reqBody := mpf.JsonMarshal(reqData)
    mss.ReqUrl = "https://api.weixin.qq.com/cgi-bin/message/mass/speed/set?access_token=" + wx.NewUtilWx().GetSingleAccessToken(mss.appId)
    client, req := mss.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := mss.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["speed"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMassSpeedSet(appId string) *massSpeedSet {
    mss := &massSpeedSet{wx.NewBaseWxAccount(), "", 0}
    mss.appId = appId
    mss.speed = -1
    mss.ReqContentType = project.HTTPContentTypeJSON
    mss.ReqMethod = fasthttp.MethodPost
    return mss
}
