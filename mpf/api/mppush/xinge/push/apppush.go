/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 16:53
 */
package push

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 消息推送
type appPush struct {
    mppush.BaseXinGe
    audienceType    string                 // 推送目标
    servicePlatform string                 // 平台类型
    messageInfo     map[string]interface{} // 消息体
    messageType     string                 // 消息类型
    pushId          int                    // 推送ID
}

func (ap *appPush) SetAudienceType(audienceType string) {
    switch audienceType {
    case "all":
        ap.audienceType = audienceType
    case "tag":
        ap.audienceType = audienceType
    case "token":
        ap.audienceType = audienceType
    case "token_list":
        ap.audienceType = audienceType
    case "account":
        ap.audienceType = audienceType
    case "account_list":
        ap.audienceType = audienceType
    default:
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "推送目标不合法", nil))
    }
}

func (ap *appPush) SetServicePlatform(servicePlatform string) {
    switch servicePlatform {
    case mppush.XinGePlatformTypeAll:
        ap.servicePlatform = servicePlatform
    case mppush.XinGePlatformTypeAndroid:
        ap.servicePlatform = servicePlatform
    case mppush.XinGePlatformTypeIOS:
        ap.servicePlatform = servicePlatform
    default:
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "平台类型不合法", nil))
    }
}

func (ap *appPush) SetMessageInfo(messageInfo map[string]interface{}) {
    if len(messageInfo) > 0 {
        ap.messageInfo = messageInfo
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "消息体不合法", nil))
    }
}

func (ap *appPush) SetMessageType(messageType string) {
    if (messageType == "notify") || (messageType == "message") {
        ap.messageType = messageType
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "消息类型不合法", nil))
    }
}

func (ap *appPush) SetPushId(pushId int) {
    if pushId >= 0 {
        ap.pushId = pushId
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "推送ID不合法", nil))
    }
}

func (ap *appPush) SetExpireTime(expireTime int) {
    if (expireTime > 0) && (expireTime <= 259200) {
        ap.ExtendData["expire_time"] = expireTime
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "消息离线存储时间不合法", nil))
    }
}

func (ap *appPush) SetSendTime(sendTime int) {
    if sendTime < time.Now().Second() {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "发送时间不合法", nil))
    }

    st := time.Unix(int64(sendTime), 0)
    ap.ExtendData["send_time"] = st.Format("2006-01-02 03:04:05")
}

func (ap *appPush) SetMultiPkg(multiPkg bool) {
    ap.ExtendData["multi_pkg"] = multiPkg
}

func (ap *appPush) SetLoopTimes(loopTimes int) {
    if (loopTimes > 0) && (loopTimes <= 15) {
        ap.ExtendData["loop_times"] = loopTimes
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "循环任务重复次数不合法", nil))
    }
}

func (ap *appPush) SetLoopInterval(loopInterval int) {
    if (loopInterval > 0) && (loopInterval <= 14) {
        ap.ExtendData["loop_interval"] = loopInterval
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "循环任务间隔天数不合法", nil))
    }
}

func (ap *appPush) SetEnvironment(environment string) {
    if (environment == "dev") || (environment == "product") {
        ap.ExtendData["environment"] = environment
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "推送环境不合法", nil))
    }
}

func (ap *appPush) SetBadgeType(badgeType int) {
    if badgeType >= -2 {
        ap.ExtendData["badge_type"] = badgeType
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "角标数字不合法", nil))
    }
}

func (ap *appPush) SetStatTag(statTag string) {
    if len(statTag) > 0 {
        ap.ExtendData["stat_tag"] = statTag
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "统计标签不合法", nil))
    }
}

func (ap *appPush) SetSeq(seq int) {
    if seq >= 0 {
        ap.ExtendData["seq"] = seq
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "请求ID不合法", nil))
    }
}

func (ap *appPush) SetTagList(tagList map[string]interface{}) {
    if len(tagList) > 0 {
        ap.ExtendData["tag_list"] = tagList
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "标签列表不合法", nil))
    }
}

func (ap *appPush) SetAccountList(accountList []string) {
    if len(accountList) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "账号列表不合法", nil))
    } else if len(accountList) > 1000 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "账号列表不能超过1000个", nil))
    }

    accounts := make([]string, 0)
    for _, v := range accountList {
        if len(v) > 0 {
            accounts = append(accounts, v)
        }
    }
    if len(accounts) > 0 {
        ap.ExtendData["account_list"] = accounts
    }
}

func (ap *appPush) SetAccountPushType(accountPushType int) {
    if (accountPushType == 0) || (accountPushType == 1) {
        ap.ExtendData["account_push_type"] = accountPushType
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "账号推送类型不合法", nil))
    }
}

func (ap *appPush) SetAccountType(accountType int) {
    if accountType >= 0 {
        ap.ExtendData["account_type"] = accountType
    } else {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "账号类型不合法", nil))
    }
}

func (ap *appPush) SetTokenList(tokenList []string) {
    if len(tokenList) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "设备列表不合法", nil))
    } else if len(tokenList) > 1000 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "设备列表不能超过1000个", nil))
    }

    tokens := make([]string, 0)
    for _, v := range tokenList {
        if len(v) > 0 {
            tokens = append(tokens, v)
        }
    }
    if len(tokens) > 0 {
        ap.ExtendData["token_list"] = tokens
    }
}

func (ap *appPush) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ap.audienceType) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "推送目标不能为空", nil))
    }
    if len(ap.servicePlatform) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "平台类型不能为空", nil))
    }
    if len(ap.messageInfo) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "消息体不能为空", nil))
    }
    if len(ap.messageType) == 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "消息类型不能为空", nil))
    }
    if ap.pushId < 0 {
        panic(mperr.NewPushXinGe(errorcode.PushXinGeParam, "推送ID不能为空", nil))
    }
    ap.ExtendData["audience_type"] = ap.audienceType
    ap.ExtendData["platform"] = ap.servicePlatform
    ap.ExtendData["message"] = ap.messageInfo
    ap.ExtendData["message_type"] = ap.messageType
    ap.ExtendData["push_id"] = strconv.Itoa(ap.pushId)

    return ap.GetRequest()
}

func NewAppPush(platform string) *appPush {
    ap := &appPush{mppush.NewBaseXinGe(platform), "", "", make(map[string]interface{}), "", -1}
    ap.pushId = -1
    ap.ServiceUri = "push/app"
    ap.ExtendData["environment"] = "product"
    return ap
}
