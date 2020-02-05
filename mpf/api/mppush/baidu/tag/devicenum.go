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

// 查询标签组设备数量
type deviceNum struct {
    mppush.BaseBaiDu
    tag string // 标签名
}

func (dn *deviceNum) SetTag(tag string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tag)
    if match {
        dn.tag = tag
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不合法", nil))
    }
}

func (dn *deviceNum) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dn.tag) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不能为空", nil))
    }
    dn.ReqData["tag"] = dn.tag

    return dn.GetRequest()
}

func NewDeviceNum() *deviceNum {
    dn := &deviceNum{mppush.NewBaseBaiDu(), ""}
    dn.ServiceUri = "/tag/device_num"
    return dn
}
