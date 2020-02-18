/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 8:52
 */
package register

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/wx"
)

type miniFastAuth struct {
    wx.BaseWxOpen
    appId string
}

func (mfa *miniFastAuth) GetResult() map[string]string {
    mfa.ReqData["appid"] = mfa.appId

    result := make(map[string]string)
    result["url"] = "https://mp.weixin.qq.com/cgi-bin/fastregisterauth?" + mpf.HTTPCreateParams(mfa.ReqData, "none", 1)
    return result
}

func NewMiniFastAuth(appId string) *miniFastAuth {
    conf := wx.NewConfig().GetOpen()
    mfa := &miniFastAuth{wx.NewBaseWxOpen(), ""}
    mfa.appId = appId
    mfa.ReqData["component_appid"] = conf.GetAppId()
    mfa.ReqData["copy_wx_verify"] = "1"
    mfa.ReqData["redirect_uri"] = conf.GetUrlMiniFastRegister()
    return mfa
}
