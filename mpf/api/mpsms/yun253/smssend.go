/**
 * 253云接口-短信发送
 * User: 姜伟
 * Date: 2019/12/23 0023
 * Time: 10:32
 */
package yun253

import (
    "regexp"
    "strings"
    "time"

    "github.com/a07061625/gompf/mpf/api/mpsms"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type smsSend struct {
    mpsms.BaseYun253
    phoneList []string // 接收手机号码列表
    signName  string   // 签名名称
    msg       string   // 短信内容
    sendTime  string   // 发送短信时间
    report    string   // 状态报告标识
}

func (s *smsSend) SetPhoneList(phoneList []string) {
    if len(phoneList) == 0 {
        panic(mperr.NewSmsYun253(errorcode.SmsYun253Param, "接收号码不能为空", nil))
    }
    if len(phoneList) > 200 {
        panic(mperr.NewSmsYun253(errorcode.SmsYun253Param, "接收号码不能超过200个", nil))
    }

    s.phoneList = make([]string, 0)
    for _, phone := range phoneList {
        match, _ := regexp.MatchString(project.RegexPhone, phone)
        if match {
            s.phoneList = append(s.phoneList, phone)
        }
    }
}

func (s *smsSend) SetSignNameAndMsg(signName string, msg string) {
    if len(signName) == 0 {
        panic(mperr.NewSmsYun253(errorcode.SmsYun253Param, "签名名称不能为空", nil))
    }
    if len(msg) == 0 {
        panic(mperr.NewSmsYun253(errorcode.SmsYun253Param, "短信内容不能为空", nil))
    }
    s.msg = "【" + signName + "】" + msg
}

func (s *smsSend) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(s.phoneList) == 0 {
        panic(mperr.NewSmsYun253(errorcode.SmsYun253Param, "接收号码不能为空", nil))
    }
    if len(s.msg) == 0 {
        panic(mperr.NewSmsYun253(errorcode.SmsYun253Param, "短信内容不能为空", nil))
    }
    s.ReqData["phone"] = strings.Join(s.phoneList, ",")
    s.ReqData["msg"] = s.msg

    return s.GetRequest()
}

func NewSmsSend() *smsSend {
    now := time.Now()
    send := &smsSend{mpsms.NewBaseYun253(), make([]string, 0), "", "", "", ""}
    conf := mpsms.NewConfigYun253()
    send.ReqUrl = conf.GetUrlSmsSend()
    send.ReqData["account"] = conf.GetAppKey()
    send.ReqData["password"] = conf.GetAppSecret()
    send.ReqData["report"] = "false"
    send.ReqData["sendtime"] = now.Format("200601020304")
    return send
}
