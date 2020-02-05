/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/5 0005
 * Time: 0:17
 */
package wx

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

type singleJsPayConfig struct {
    BaseWxSingle
    timestamp int
    nonceStr  string
    prepayId  string // 预支付交易会话标识
    signType  string // 签名类型
}

func (jpc *singleJsPayConfig) SetTimestamp(timestamp int) {
    if timestamp > 1000000000 {
        jpc.timestamp = timestamp
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "时间戳不合法", nil))
    }
}

func (jpc *singleJsPayConfig) SetPrepayId(prepayId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,64}$`, prepayId)
    if match {
        jpc.ReqData["package"] = "prepay_id=" + prepayId
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "预支付交易会话标识不合法", nil))
    }
}

func (jpc *singleJsPayConfig) checkData() {
    if jpc.timestamp <= 0 {
        panic(mperr.NewWx(errorcode.WxParam, "时间戳不能为空", nil))
    }
    jpc.ReqData["timeStamp"] = strconv.Itoa(jpc.timestamp)
    _, ok := jpc.ReqData["package"]
    if !ok {
        panic(mperr.NewWx(errorcode.WxParam, "交易会话标识不能为空", nil))
    }
}

func (jpc *singleJsPayConfig) GetResult() map[string]string {
    jpc.checkData()
    jpc.ReqData["paySign"] = NewUtilWx().CreateSinglePaySign(jpc.ReqData, jpc.ReqData["appId"], "md5")
    return jpc.ReqData
}

func NewSingleJsPayConfig(appId string) *singleJsPayConfig {
    jpc := &singleJsPayConfig{NewBaseWxSingle(), 0, "", "", ""}
    jpc.ReqData["appId"] = appId
    jpc.ReqData["signType"] = "MD5"
    jpc.ReqData["nonceStr"] = mpf.ToolCreateNonceStr(32, "numlower")
    return jpc
}
