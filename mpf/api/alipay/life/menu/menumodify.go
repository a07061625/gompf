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

// 默认菜单更新接口
type menuModify struct {
    alipay.BaseAliPay
    menuType   string                   // 菜单类型
    buttonList []map[string]interface{} // 菜单列表
}

func (mm *menuModify) SetMenuType(menuType string) {
    if (menuType == "icon") || (menuType == "text") {
        mm.menuType = menuType
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单类型不合法", nil))
    }
}

func (mm *menuModify) SetButtonList(buttonList []map[string]interface{}) {
    if len(buttonList) > 0 {
        mm.buttonList = buttonList
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单列表不合法", nil))
    }
}

func (mm *menuModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(mm.menuType) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单类型不能为空", nil))
    }
    if len(mm.buttonList) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单列表不能为空", nil))
    }
    mm.BizContent["type"] = mm.menuType
    mm.BizContent["button"] = mm.buttonList

    return mm.GetRequest()
}

func NewMenuModify(appId string) *menuModify {
    mm := &menuModify{alipay.NewBase(appId), "", make([]map[string]interface{}, 0)}
    mm.SetMethod("alipay.open.public.menu.modify")
    return mm
}
