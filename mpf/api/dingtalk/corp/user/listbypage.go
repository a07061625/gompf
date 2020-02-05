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

// 获取部门用户详情
type listByPage struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    departId int // 部门id
}

func (lp *listByPage) SetDepartId(departId int) {
    if departId > 0 {
        lp.departId = departId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不合法", nil))
    }
}

func (lp *listByPage) SetLang(lang string) {
    if (lang == "zh_CN") || (lang == "en_US") {
        lp.ReqData["lang"] = lang
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "语言不合法", nil))
    }
}

func (lp *listByPage) SetOffset(offset int) {
    if offset >= 0 {
        lp.ReqData["offset"] = strconv.Itoa(offset)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (lp *listByPage) SetSize(size int) {
    if (size > 0) && (size <= 100) {
        lp.ReqData["size"] = strconv.Itoa(size)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

// 排序 entry_asc: 进入时间升序 entry_desc: 进入时间降序 modify_asc: 修改时间升序 modify_desc: 修改时间降序 custom: 用户定义排序,默认
func (lp *listByPage) SetOrder(order string) {
    switch order {
    case "entry_asc":
        lp.ReqData["order"] = order
    case "entry_desc":
        lp.ReqData["order"] = order
    case "modify_asc":
        lp.ReqData["order"] = order
    case "modify_desc":
        lp.ReqData["order"] = order
    case "custom":
        lp.ReqData["order"] = order
    default:
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "排序不合法", nil))
    }
}

func (lp *listByPage) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if lp.departId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "部门id不能为空", nil))
    }
    lp.ReqData["department_id"] = strconv.Itoa(lp.departId)
    lp.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(lp.corpId, lp.agentTag, lp.atType)
    lp.ReqUrl = dingtalk.UrlService + "/user/listbypage?" + mpf.HttpCreateParams(lp.ReqData, "none", 1)

    return lp.GetRequest()
}

func NewListByPage(corpId, agentTag, atType string) *listByPage {
    lp := &listByPage{dingtalk.NewCorp(), "", "", "", 0}
    lp.corpId = corpId
    lp.agentTag = agentTag
    lp.atType = atType
    lp.ReqData["lang"] = "zh_CN"
    lp.ReqData["offset"] = "0"
    lp.ReqData["size"] = "10"
    lp.ReqData["order"] = "custom"
    return lp
}
