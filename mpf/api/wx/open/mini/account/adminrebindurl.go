/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 17:07
 */
package account

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/wx"
)

// 管理员换绑链接
type adminRebindUrl struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (aru *adminRebindUrl) GetResult() map[string]string {
    conf := wx.NewConfig().GetOpen()
    aru.ReqData["appid"] = aru.appId
    aru.ReqData["component_appid"] = conf.GetAppId()
    aru.ReqData["redirect_uri"] = conf.GetUrlMiniRebindAdmin()

    result := make(map[string]string)
    result["url"] = "https://mp.weixin.qq.com/wxopen/componentrebindadmin?" + mpf.HTTPCreateParams(aru.ReqData, "none", 1)
    return result
}

func NewAdminRebindUrl(appId string) *adminRebindUrl {
    aru := &adminRebindUrl{wx.NewBaseWxOpen(), ""}
    aru.appId = appId
    return aru
}
