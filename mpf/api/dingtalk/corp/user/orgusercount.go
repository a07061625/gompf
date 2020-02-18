/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 1:05
 */
package user

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取企业员工人数
type orgUserCount struct {
    dingtalk.BaseCorp
    corpId     string
    agentTag   string
    atType     string
    onlyActive int // 激活钉钉标识 0:包含未激活钉钉的人员数量 1:不包含未激活钉钉的人员数量
}

func (ouc *orgUserCount) SetOnlyActive(onlyActive int) {
    if (onlyActive == 0) || (onlyActive == 1) {
        ouc.onlyActive = onlyActive
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "激活钉钉标识不合法", nil))
    }
}

func (ouc *orgUserCount) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if ouc.onlyActive < 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "激活钉钉标识不能为空", nil))
    }
    ouc.ReqData["onlyActive"] = strconv.Itoa(ouc.onlyActive)
    ouc.ReqData["access_token"] = dingtalk.NewUtil().GetAccessToken(ouc.corpId, ouc.agentTag, ouc.atType)
    ouc.ReqUrl = dingtalk.UrlService + "/user/get_org_user_count?" + mpf.HTTPCreateParams(ouc.ReqData, "none", 1)

    return ouc.GetRequest()
}

func NewOrgUserCount(corpId, agentTag, atType string) *orgUserCount {
    ouc := &orgUserCount{dingtalk.NewCorp(), "", "", "", -1}
    ouc.corpId = corpId
    ouc.agentTag = agentTag
    ouc.atType = atType
    ouc.onlyActive = -1
    return ouc
}
