/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 16:30
 */
package account

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 管理员换绑
type adminRebind struct {
    wx.BaseWxOpen
    taskId string // 任务ID
}

func (ar *adminRebind) SetTaskId(taskId string) {
    if len(taskId) > 0 {
        ar.taskId = taskId
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "任务ID不合法", nil))
    }
}

func (ar *adminRebind) checkData() {
    if len(ar.taskId) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "任务ID不能为空", nil))
    }
    ar.ReqData["taskid"] = ar.taskId
}

func (ar *adminRebind) SendRequest() api.ApiResult {
    ar.checkData()

    reqBody := mpf.JsonMarshal(ar.ReqData)
    ar.ReqUrl = "https://api.weixin.qq.com/cgi-bin/account/componentrebindadmin?access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := ar.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ar.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewAdminRebind() *adminRebind {
    ar := &adminRebind{wx.NewBaseWxOpen(), ""}
    ar.ReqContentType = project.HTTPContentTypeJSON
    ar.ReqMethod = fasthttp.MethodPost
    return ar
}
