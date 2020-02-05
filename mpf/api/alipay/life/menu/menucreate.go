/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 22:43
 */
package menu

import (
    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 默认菜单创建接口
type menuCreate struct {
    alipay.BaseAliPay
    menuType   string                   // 菜单类型
    buttonList []map[string]interface{} // 菜单列表
}

func (mc *menuCreate) SetMenuType(menuType string) {
    if (menuType == "icon") || (menuType == "text") {
        mc.menuType = menuType
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单类型不合法", nil))
    }
}

func (mc *menuCreate) SetButtonList(buttonList []map[string]interface{}) {
    if len(buttonList) > 0 {
        mc.buttonList = buttonList
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单列表不合法", nil))
    }
}

func (mc *menuCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(mc.buttonList) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单列表不能为空", nil))
    }
    mc.BizContent["type"] = mc.menuType
    mc.BizContent["button"] = mc.buttonList

    return mc.GetRequest()
}

func NewMenuCreate(appId string) *menuCreate {
    mc := &menuCreate{alipay.NewBase(appId), "", make([]map[string]interface{}, 0)}
    mc.menuType = "text"
    mc.SetMethod("alipay.open.public.menu.create")
    return mc
}
