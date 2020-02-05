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

// 个性化菜单创建
type personalCreate struct {
    alipay.BaseAliPay
    menuType   string                   // 菜单类型
    buttonList []map[string]interface{} // 菜单列表
    labelRule  []map[string]interface{} // 标签规则
}

func (pc *personalCreate) SetMenuType(menuType string) {
    if (menuType == "icon") || (menuType == "text") {
        pc.menuType = menuType
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单类型不合法", nil))
    }
}

func (pc *personalCreate) SetButtonList(buttonList []map[string]interface{}) {
    if len(buttonList) > 0 {
        pc.buttonList = buttonList
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单列表不合法", nil))
    }
}

func (pc *personalCreate) SetLabelRule(labelRule []map[string]interface{}) {
    if len(labelRule) > 0 {
        pc.labelRule = labelRule
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "标签规则不合法", nil))
    }
}

func (pc *personalCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pc.buttonList) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单列表不能为空", nil))
    }
    if len(pc.labelRule) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "标签规则不能为空", nil))
    }
    pc.BizContent["type"] = pc.menuType
    pc.BizContent["button"] = pc.buttonList
    pc.BizContent["label_rule"] = pc.labelRule

    return pc.GetRequest()
}

func NewPersonalCreate(appId string) *personalCreate {
    pc := &personalCreate{alipay.NewBase(appId), "", make([]map[string]interface{}, 0), make([]map[string]interface{}, 0)}
    pc.menuType = "text"
    pc.SetMethod("alipay.open.public.personalized.menu.create")
    return pc
}
