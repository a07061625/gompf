/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 22:30
 */
package tag

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 将设备从标签组中移除
type devicesDel struct {
    mppush.BaseBaiDu
    tag         string   // 标签名
    channelList []string // 设备列表
}

func (dd *devicesDel) SetTag(tag string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tag)
    if match {
        dd.tag = tag
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不合法", nil))
    }
}

func (dd *devicesDel) SetChannelList(channelList []string) {
    if len(channelList) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备列表不能为空", nil))
    } else if len(channelList) > 10 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备列表不能超过10个", nil))
    }
    dd.channelList = make([]string, 0)
    for _, v := range channelList {
        if len(v) > 0 {
            dd.channelList = append(dd.channelList, v)
        }
    }
}

func (dd *devicesDel) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dd.tag) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不能为空", nil))
    }
    if len(dd.channelList) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备列表不能为空", nil))
    }
    dd.ReqData["tag"] = dd.tag
    dd.ReqData["channel_ids"] = mpf.JSONMarshal(dd.channelList)

    return dd.GetRequest()
}

func NewDevicesDel() *devicesDel {
    dd := &devicesDel{mppush.NewBaseBaiDu(), "", make([]string, 0)}
    dd.ServiceUri = "/tag/del_devices"
    return dd
}
