/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 0:26
 */
package address

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/taobao"
    "github.com/a07061625/gompf/mpf/api/taobao/logistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 批量判定服务是否可达
type reachableBatch struct {
    taobao.BaseTaoBao
    addressList []map[string]interface{} // 地址列表
}

func (rb *reachableBatch) SetAddressList(addressList []map[string]interface{}) {
    if len(addressList) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "地址列表不能为空", nil))
    } else if len(addressList) > 20 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "地址列表长度不能超过20个", nil))
    }
    rb.addressList = addressList
}

func (rb *reachableBatch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rb.addressList) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "地址列表不能为空", nil))
    }
    rb.ReqData["address_list"] = mpf.JSONMarshal(rb.addressList)

    return rb.GetRequest()
}

func NewReachableBatch() *reachableBatch {
    rb := &reachableBatch{taobao.NewBaseTaoBao(), make([]map[string]interface{}, 0)}
    conf := logistics.NewConfigTaoBao()
    rb.AppKey = conf.GetAppKey()
    rb.AppSecret = conf.GetAppSecret()
    rb.SetMethod("taobao.logistics.address.reachablebatch.get")
    return rb
}
