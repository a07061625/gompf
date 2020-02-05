/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 23:36
 */
package service

import (
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 获取推广注册链接
type registerUrl struct {
    wx.BaseWxProvider
    registerCode string // 注册码
}

func (ru *registerUrl) SetRegisterCode(registerCode string) {
    if len(registerCode) > 0 {
        ru.registerCode = registerCode
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "注册码不合法", nil))
    }
}

func (ru *registerUrl) checkData() {
    if len(ru.registerCode) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "注册码不能为空", nil))
    }
}

func (ru *registerUrl) GetResult() map[string]string {
    ru.checkData()

    result := make(map[string]string)
    result["url"] = "https://open.work.weixin.qq.com/3rdservice/wework/register?register_code=" + ru.registerCode
    return result
}

func NewRegisterUrl() *registerUrl {
    return &registerUrl{wx.NewBaseWxProvider(), ""}
}
