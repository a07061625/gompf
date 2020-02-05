/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 23:14
 */
package menu

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 个性化菜单删除
type personalDelete struct {
    alipay.BaseAliPay
    menuKey string // 菜单key
}

func (pd *personalDelete) SetMenuKey(menuKey string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, menuKey)
    if match {
        pd.menuKey = menuKey
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单key不合法", nil))
    }
}

func (pd *personalDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pd.menuKey) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "菜单key不能为空", nil))
    }
    pd.BizContent["menu_key"] = pd.menuKey

    return pd.GetRequest()
}

func NewPersonalDelete(appId string) *personalDelete {
    pd := &personalDelete{alipay.NewBase(appId), ""}
    pd.SetMethod("alipay.open.public.personalized.menu.delete")
    return pd
}
