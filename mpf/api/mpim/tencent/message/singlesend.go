/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 14:52
 */
package message

import (
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpim"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type singleSend struct {
    mpim.BaseTencent
    syncFlag    int                      // 同步标识 1:同步 2:不同步
    fromAccount string                   // 发送方帐号
    toAccount   string                   // 接收方帐号
    lifeTime    int                      // 离线保存秒数,0为不保存
    body        []map[string]interface{} // 消息内容
    offlinePush map[string]interface{}   // 离线推送数据
}

func (ss *singleSend) SetSyncFlag(syncFlag int) {
    if (syncFlag == 0) || (syncFlag == 1) {
        ss.syncFlag = syncFlag
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "同步标识不合法", nil))
    }
}

func (ss *singleSend) SetFromAccount(fromAccount string) {
    ss.fromAccount = fromAccount
}

func (ss *singleSend) SetToAccount(toAccount string) {
    if len(toAccount) > 0 {
        ss.toAccount = toAccount
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "接收方帐号不合法", nil))
    }
}

func (ss *singleSend) SetLifeTime(lifeTime int) {
    if (lifeTime >= 0) || (lifeTime <= 604800) {
        ss.lifeTime = lifeTime
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "离线保存时间不合法", nil))
    }
}

func (ss *singleSend) SetBody(body []map[string]interface{}) {
    if len(body) > 0 {
        ss.body = body
    } else {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "消息内容不合法", nil))
    }
}

func (ss *singleSend) SetOfflinePush(offlinePush map[string]interface{}) {
    ss.offlinePush = offlinePush
}

func (ss *singleSend) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ss.toAccount) == 0 {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "接收方帐号不能为空", nil))
    }
    if len(ss.body) == 0 {
        panic(mperr.NewIMTencent(errorcode.IMTencentParam, "消息内容不能为空", nil))
    }
    ss.ExtendData["SyncOtherMachine"] = ss.syncFlag
    ss.ExtendData["To_Account"] = ss.toAccount
    ss.ExtendData["MsgLifeTime"] = ss.lifeTime
    ss.ExtendData["MsgBody"] = ss.body
    if len(ss.fromAccount) > 0 {
        ss.ExtendData["From_Account"] = ss.fromAccount
    }
    if len(ss.offlinePush) > 0 {
        ss.ExtendData["OfflinePushInfo"] = ss.offlinePush
    }
    reqBody := mpf.JsonMarshal(ss.ExtendData)

    client, req := ss.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewSingleSend() *singleSend {
    ss := &singleSend{mpim.NewBaseTencent(), 0, "", "", 0, make([]map[string]interface{}, 0), make(map[string]interface{})}
    ss.syncFlag = 1
    ss.lifeTime = 604800
    ss.ServiceUri = "/openim/sendmsg"
    ss.ExtendData["MsgRandom"] = mpf.ToolCreateRandNum(10000000, 89999999)
    ss.ExtendData["MsgTimeStamp"] = time.Now().Unix()
    ss.ReqContentType = project.HttpContentTypeJson
    ss.ReqMethod = fasthttp.MethodPost
    return ss
}
