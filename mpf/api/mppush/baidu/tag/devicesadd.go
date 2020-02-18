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

// 添加设备到标签组
type devicesAdd struct {
    mppush.BaseBaiDu
    tag         string   // 标签名
    channelList []string // 设备列表
}

func (da *devicesAdd) SetTag(tag string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tag)
    if match {
        da.tag = tag
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不合法", nil))
    }
}

func (da *devicesAdd) SetChannelList(channelList []string) {
    if len(channelList) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备列表不能为空", nil))
    } else if len(channelList) > 10 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备列表不能超过10个", nil))
    }
    da.channelList = make([]string, 0)
    for _, v := range channelList {
        if len(v) > 0 {
            da.channelList = append(da.channelList, v)
        }
    }
}

func (da *devicesAdd) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(da.tag) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "标签名不能为空", nil))
    }
    if len(da.channelList) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "设备列表不能为空", nil))
    }
    da.ReqData["tag"] = da.tag
    da.ReqData["channel_ids"] = mpf.JSONMarshal(da.channelList)

    return da.GetRequest()
}

func NewDevicesAdd() *devicesAdd {
    da := &devicesAdd{mppush.NewBaseBaiDu(), "", make([]string, 0)}
    da.ServiceUri = "/tag/add_devices"
    return da
}
