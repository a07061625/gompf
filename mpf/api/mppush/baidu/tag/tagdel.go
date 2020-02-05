/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 22:30
 */
package tag

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除标签组
type tagDel struct {
    mppush.BaseBaiDu
    tag string // 标签名
}

func (td *tagDel) SetTag(tag string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tag)
    if match {
        td.tag = tag
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不合法", nil))
    }
}

func (td *tagDel) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(td.tag) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不能为空", nil))
    }
    td.ReqData["tag"] = td.tag

    return td.GetRequest()
}

func NewTagDel() *tagDel {
    td := &tagDel{mppush.NewBaseBaiDu(), ""}
    td.ServiceUri = "/app/del_tag"
    return td
}
