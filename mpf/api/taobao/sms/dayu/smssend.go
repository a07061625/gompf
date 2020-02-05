/**
 * 大鱼接口-短信发送
 * User: 姜伟
 * Date: 2019/12/24 0024
 * Time: 23:19
 */
package dayu

import (
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/taobao"
    "github.com/a07061625/gompf/mpf/api/taobao/sms"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 短信发送
type smsSend struct {
    taobao.BaseTaoBao
    smsType        string            // 短信类型
    recNumList     []string          // 接收手机号码列表
    signName       string            // 签名名称
    templateCode   string            // 模板ID
    templateParams map[string]string // 模板参数
}

func (s *smsSend) SetRecNumList(recNumList []string) {
    if len(recNumList) == 0 {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "接收号码不能为空", nil))
    }
    if len(recNumList) > 200 {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "接收号码不能超过200个", nil))
    }

    s.recNumList = make([]string, 0)
    for _, recNum := range recNumList {
        match, _ := regexp.MatchString(project.RegexPhone, recNum)
        if match {
            s.recNumList = append(s.recNumList, recNum)
        }
    }
}

func (s *smsSend) SetSignName(signName string) {
    if len(signName) == 0 {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "签名名称不能为空", nil))
    }
    _, ok := badSignNames[signName]
    if ok {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "签名名称不能为系统默认签名", nil))
    }

    s.signName = signName
}

func (s *smsSend) SetTemplateCode(templateCode string) {
    if len(templateCode) > 0 {
        s.templateCode = templateCode
    } else {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "模板ID不能为空", nil))
    }
}

func (s *smsSend) SetTemplateParams(templateParams map[string]string) {
    if len(templateParams) > 0 {
        s.ReqData["sms_param"] = mpf.JsonMarshal(templateParams)
    }
}

func (s *smsSend) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(s.recNumList) == 0 {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "接收号码不能为空", nil))
    }
    if len(s.signName) == 0 {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "签名名称不能为空", nil))
    }
    if len(s.templateCode) == 0 {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "模板ID不能为空", nil))
    }
    s.ReqData["rec_num"] = strings.Join(s.recNumList, ",")
    s.ReqData["sms_free_sign_name"] = s.signName
    s.ReqData["sms_template_code"] = s.templateCode

    return s.GetRequest()
}

var (
    badSignNames map[string]int
)

func init() {
    badSignNames = make(map[string]int)
    badSignNames["大鱼测试"] = 1
    badSignNames["活动验证"] = 1
    badSignNames["变更验证"] = 1
    badSignNames["登录验证"] = 1
    badSignNames["注册验证"] = 1
    badSignNames["身份验证"] = 1
}

func NewSmsSend() *smsSend {
    send := &smsSend{taobao.NewBaseTaoBao(), "", make([]string, 0), "", "", make(map[string]string)}
    conf := sms.NewConfigDaYu()
    send.AppKey = conf.GetAppKey()
    send.AppSecret = conf.GetAppSecret()
    send.ReqData["sms_type"] = "normal"
    send.SetMethod("alibaba.aliqin.fc.mpsms.num.send")
    return send
}
