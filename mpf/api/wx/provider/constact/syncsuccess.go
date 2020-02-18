/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 22:55
 */
package constact

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 设置通讯录同步完成
type syncSuccess struct {
    wx.BaseWxProvider
    accessToken string // 令牌,由查询注册状态接口返回
}

func (ss *syncSuccess) SetAccessToken(accessToken string) {
    if len(accessToken) > 0 {
        ss.accessToken = accessToken
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "令牌不合法", nil))
    }
}

func (ss *syncSuccess) checkData() {
    if len(ss.accessToken) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "令牌不能为空", nil))
    }
}

func (ss *syncSuccess) SendRequest() api.ApiResult {
    ss.checkData()

    ss.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/sync/contact_sync_success?access_token=" + ss.accessToken
    client, req := ss.GetRequest()

    resp, result := ss.SendInner(client, req, errorcode.WxProviderRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxProviderRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewSyncSuccess() *syncSuccess {
    return &syncSuccess{wx.NewBaseWxProvider(), ""}
}
