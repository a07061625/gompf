/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 9:05
 */
package trace

import (
    "time"

    "github.com/a07061625/gompf/mpf/api/taobao"
    "github.com/a07061625/gompf/mpf/api/taobao/logistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 物流订单流转信息推送接口
type orderTracePush struct {
    taobao.BaseTaoBao
    mailNo          string // 快递单号
    ocCureTime      string // 流转节点发生时间
    operateDetail   string // 流转节点详情
    companyName     string // 物流公司名称
    operatorName    string // 快递业务员名称
    operatorContact string // 快递业务员联系方式
    currentCity     string // 流转节点城市
    facilityName    string // 网点名称
    nodeDescription string // 流转节点描述 TMS_ACCEPT:揽收 TMS_DELIVERING:派送 TMS_SIGN:签收
}

func (otp *orderTracePush) SetMailNo(mailNo string) {
    if len(mailNo) > 0 {
        otp.mailNo = mailNo
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "运单号不合法", nil))
    }
}

func (otp *orderTracePush) SetOcCureTime(ocCureTime int) {
    if ocCureTime > 0 {
        t := time.Unix(int64(ocCureTime), 0)
        otp.ocCureTime = t.Format("2006-01-02 03:04:05")
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "流转节点发生时间不合法", nil))
    }
}

func (otp *orderTracePush) SetOperateDetail(operateDetail string) {
    if len(operateDetail) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "流转节点详情不能为空", nil))
    } else if len(operateDetail) > 200 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "流转节点详情长度不能超过200字节", nil))
    }
    otp.operateDetail = operateDetail
}

func (otp *orderTracePush) SetCompanyName(companyName string) {
    if len(companyName) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "物流公司名称不能为空", nil))
    } else if len(companyName) > 20 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "物流公司名称长度不能超过20字节", nil))
    }
    otp.companyName = companyName
}

func (otp *orderTracePush) SetOperatorName(operatorName string) {
    if len(operatorName) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "快递业务员名称不能为空", nil))
    } else if len(operatorName) > 20 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "快递业务员名称长度不能超过20字节", nil))
    }
    otp.ReqData["operator_name"] = operatorName
}

func (otp *orderTracePush) SetOperatorContact(operatorContact string) {
    if len(operatorContact) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "快递业务员联系方式不能为空", nil))
    } else if len(operatorContact) > 20 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "快递业务员联系方式长度不能超过20字节", nil))
    }
    otp.ReqData["operator_contact"] = operatorContact
}

func (otp *orderTracePush) SetCurrentCity(currentCity string) {
    if len(currentCity) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "流转节点城市不能为空", nil))
    } else if len(currentCity) > 20 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "流转节点城市长度不能超过20字节", nil))
    }
    otp.ReqData["current_city"] = currentCity
}

func (otp *orderTracePush) SetFacilityName(facilityName string) {
    if len(facilityName) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "网点名称不能为空", nil))
    } else if len(facilityName) > 100 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "网点名称长度不能超过100字节", nil))
    }
    otp.ReqData["facility_name"] = facilityName
}

func (otp *orderTracePush) SetModeDescription(nodeDescription string) {
    if (nodeDescription == "TMS_ACCEPT") || (nodeDescription == "TMS_DELIVERING") || (nodeDescription == "TMS_SIGN") {
        otp.ReqData["node_description"] = nodeDescription
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "流转节点描述不合法", nil))
    }
}

func (otp *orderTracePush) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(otp.mailNo) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "快递单号不能为空", nil))
    }
    if len(otp.ocCureTime) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "流转节点发生时间不能为空", nil))
    }
    if len(otp.operateDetail) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "流转节点详情不能为空", nil))
    }
    if len(otp.companyName) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "物流公司名称不能为空", nil))
    }
    otp.ReqData["mail_no"] = otp.mailNo
    otp.ReqData["occure_time"] = otp.ocCureTime
    otp.ReqData["operate_detail"] = otp.operateDetail
    otp.ReqData["company_name"] = otp.companyName

    return otp.GetRequest()
}

func NewOrderTracePush() *orderTracePush {
    otp := &orderTracePush{taobao.NewBaseTaoBao(), "", "", "", "", "", "", "", "", ""}
    conf := logistics.NewConfigTaoBao()
    otp.AppKey = conf.GetAppKey()
    otp.AppSecret = conf.GetAppSecret()
    otp.SetMethod("taobao.logistics.ordertrace.push")
    return otp
}
