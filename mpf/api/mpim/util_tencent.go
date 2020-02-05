/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 15:27
 */
package mpim

import (
    "os/exec"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

func (util *utilIM) GetTencentAccountSign(account string) string {
    return util.cache.GetTencentAccountSign(account)
}

func (util *utilIM) DelTencentAccountSign(account string) {
    util.cache.DelTencentAccountSign(account)
}

func (util *utilIM) CreateTencentSign(userTag string) string {
    conf := NewConfigTencent()
    cmd := exec.Command(conf.GetFileCommand(), conf.GetFilePrivateKey(), conf.GetAppId(), userTag)
    output, err := cmd.Output()
    if err != nil {
        panic(mperr.NewIMTencent(errorcode.IMTencentSign, "生成即时通讯签名失败", nil))
    }
    outList := strings.Split(string(output), "\n")
    return outList[0]
}

func (util *utilIM) SendTencentRequest(service api.IApiOuter, errorCode uint) api.ApiResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    actionStatus, ok := respData["ActionStatus"]
    if ok && (actionStatus.(string) == "OK") {
        result.Data = respData
    } else {
        result.Code = errorCode
        result.Msg = respData["ErrorInfo"].(string)
    }
    return result
}
