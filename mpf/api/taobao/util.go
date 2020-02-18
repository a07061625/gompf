/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/2 0002
 * Time: 9:32
 */
package taobao

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
)

type IBaseTaoBao interface {
    api.IApiOuter
    SetMethod(method string)
    GetRespTag() string
}

type utilTaoBao struct {
    api.UtilApi
}

// 发送请求
func (util *utilTaoBao) SendRequest(service IBaseTaoBao, errorCode uint) api.ApiResult {
    resp, result := util.SendOuter(service, errorCode)
    if result.Code > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    data, ok := respData[service.GetRespTag()]
    if ok {
        result.Data = data
        return result
    }

    result.Code = errorCode
    errResp, ok := respData["error_response"]
    if !ok {
        result.Msg = "解析服务数据出错"
        return result
    }

    errInfo := errResp.(map[string]interface{})
    errMsg, ok := errInfo["sub_msg"]
    if ok {
        result.Msg = errMsg.(string)
    } else {
        result.Msg = errInfo["msg"].(string)
    }

    return result
}

var (
    insUtilTaoBao *utilTaoBao
)

func init() {
    insUtilTaoBao = &utilTaoBao{api.NewUtilApi()}
}

func NewUtilTaoBao() *utilTaoBao {
    return insUtilTaoBao
}
