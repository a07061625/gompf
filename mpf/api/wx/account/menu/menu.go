/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 11:48
 */
package menu

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type menu struct {
    name       string                   // 菜单标题
    subButton  []map[string]interface{} // 子菜单
    actionType string                   // 响应动作类型
    key        string                   // 菜单KEY值，用于消息接口推送
    url        string                   // 网页链接，用户点击菜单可打开链接
    mediaId    string                   // 媒体ID
}

func (m *menu) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        m.name = string(trueName[:5])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "菜单名称不合法", nil))
    }
}

func (m *menu) AddSubButton(subButton map[string]interface{}) {
    if len(subButton) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "子菜单不能为空", nil))
    } else if len(m.subButton) >= 5 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "子菜单不能超过5个", nil))
    }
    m.subButton = append(m.subButton, subButton)
}

func (m *menu) SetActionType(actionType string) {
    _, ok := menuTypes[actionType]
    if ok {
        m.actionType = actionType
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "响应动作类型不合法", nil))
    }
}

func (m *menu) SetKey(key string) {
    if (len(key) > 0) && (len(key) <= 128) {
        m.key = key
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "菜单KEY不合法", nil))
    }
}

func (m *menu) SetUrl(url string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, url)
    if match {
        m.url = url
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "网页链接不合法", nil))
    }
}

func (m *menu) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        m.mediaId = mediaId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "媒体ID不合法", nil))
    }
}

func (m *menu) checkData() {
    if len(m.name) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "菜单名称不能为空", nil))
    }
    if len(m.actionType) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "响应动作类型不能为空", nil))
    }
}

func (m *menu) GetResult() map[string]interface{} {
    m.checkData()

    result := make(map[string]interface{})
    result["name"] = m.name
    result["type"] = m.actionType
    result["sub_button"] = m.subButton
    if len(m.key) > 0 {
        result["key"] = m.key
    }
    if len(m.url) > 0 {
        result["url"] = m.url
    }
    if len(m.mediaId) > 0 {
        result["media_id"] = m.mediaId
    }
    return result
}

func NewMenu(appId string) *menu {
    return &menu{"", make([]map[string]interface{}, 0), "", "", "", ""}
}

var (
    menuTypes map[string]int
)

func init() {
    menuTypes = make(map[string]int)
    menuTypes["pic_weixin"] = 1
    menuTypes["pic_sysphoto"] = 1
    menuTypes["pic_photo_or_album"] = 1
    menuTypes["view"] = 1
    menuTypes["view_limited"] = 1
    menuTypes["click"] = 1
    menuTypes["media_id"] = 1
    menuTypes["location_select"] = 1
    menuTypes["scancode_push"] = 1
    menuTypes["scancode_waitmsg"] = 1
}
