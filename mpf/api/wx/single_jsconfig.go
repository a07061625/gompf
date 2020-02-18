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
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type singleJsConfig struct {
    BaseWxSingle
    timestamp int64
    nonceStr  string
    url       string
}

func (jc *singleJsConfig) SetTimestamp(timestamp int64) {
    if timestamp > 1000000000 {
        jc.timestamp = timestamp
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "时间戳不合法", nil))
    }
}

func (jc *singleJsConfig) SetNonceStr(nonceStr string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{32}$`, nonceStr)
    if match {
        jc.ReqData["nonceStr"] = nonceStr
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "随机字符串不合法", nil))
    }
}

func (jc *singleJsConfig) SetUrl(url string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, url)
    if match {
        jc.url = url
    } else {
        panic(mperr.NewWx(errorcode.WxParam, "链接地址不合法", nil))
    }
}

func (jc *singleJsConfig) GetResult(getType string) map[string]string {
    ticket := NewUtilWx().GetSingleCache(jc.ReqData["appId"], getType)
    jc.ReqData["timestamp"] = strconv.FormatInt(jc.timestamp, 10)
    signStr := "jsapi_ticket=" + ticket + "&noncestr=" + jc.ReqData["nonceStr"] + "&timestamp=" + jc.ReqData["timestamp"] + "&url=" + jc.url
    jc.ReqData["signature"] = mpf.HashSha1(signStr, "")
    return jc.ReqData
}

func NewSingleJsConfig(appId string) *singleJsConfig {
    conf := NewConfig().GetAccount(appId)
    jc := &singleJsConfig{NewBaseWxSingle(), 0, "", ""}
    jc.ReqData["appId"] = conf.GetAppId()
    jc.ReqData["nonceStr"] = mpf.ToolCreateNonceStr(32, "numlower")
    jc.timestamp = time.Now().Unix()
    jc.url = conf.GetPayAuthUrl()
    return jc
}
