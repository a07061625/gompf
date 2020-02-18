/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 1:05
 */
package user

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取部门用户
type simpleList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    departId int // 部门id
}

func (sl *simpleList) SetDepartId(departId int) {
    if departId > 0 {
        sl.departId = departId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不合法", nil))
    }
}

func (sl *simpleList) SetLang(lang string) {
    if (lang == "zh_CN") || (lang == "en_US") {
        sl.ReqData["lang"] = lang
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "语言不合法", nil))
    }
}

func (sl *simpleList) SetOffset(offset int) {
    if offset >= 0 {
        sl.ReqData["offset"] = strconv.Itoa(offset)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (sl *simpleList) SetSize(size int) {
    if (size > 0) && (size <= 100) {
        sl.ReqData["size"] = strconv.Itoa(size)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

// 排序 entry_asc: 进入时间升序 entry_desc: 进入时间降序 modify_asc: 修改时间升序 modify_desc: 修改时间降序 custom: 用户定义排序,默认
func (sl *simpleList) SetOrder(order string) {
    switch order {
    case "entry_asc":
        sl.ReqData["order"] = order
    case "entry_desc":
        sl.ReqData["order"] = order
    case "modify_asc":
        sl.ReqData["order"] = order
    case "modify_desc":
        sl.ReqData["order"] = order
    case "custom":
        sl.ReqData["order"] = order
    default:
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "排序不合法", nil))
    }
}

func (sl *simpleList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if sl.departId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不能为空", nil))
    }
    sl.ReqData["department_id"] = strconv.Itoa(sl.departId)
    sl.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(sl.corpId, sl.agentTag, sl.atType)
    sl.ReqUrl = dingtalk.UrlService + "/user/simplelist?" + mpf.HTTPCreateParams(sl.ReqData, "none", 1)

    return sl.GetRequest()
}

func NewSimpleList(corpId, agentTag, atType string) *simpleList {
    sl := &simpleList{dingtalk.NewCorp(), "", "", "", 0}
    sl.corpId = corpId
    sl.agentTag = agentTag
    sl.atType = atType
    sl.ReqData["lang"] = "zh_CN"
    sl.ReqData["offset"] = "0"
    sl.ReqData["size"] = "10"
    sl.ReqData["order"] = "custom"
    return sl
}
