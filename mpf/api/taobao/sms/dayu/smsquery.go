/**
 * 大鱼接口-短信查询
 * User: 姜伟
 * Date: 2019/12/23 0023
 * Time: 10:31
 */
package dayu

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf/api/taobao"
    "github.com/a07061625/gompf/mpf/api/taobao/sms"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 短信查询
type smsQuery struct {
    taobao.BaseTaoBao
    bizId     string // 流水号
    recNum    string // 接收号码
    queryDate string // 发送日期
    page      uint   // 页码
    limit     uint   // 每页数量
}

func (q *smsQuery) SetBizId(bizId string) {
    if len(bizId) > 0 {
        q.ReqData["biz_id"] = bizId
    }
}

func (q *smsQuery) SetRecNum(recNum string) {
    match, _ := regexp.MatchString(project.RegexPhone, recNum)
    if match {
        q.recNum = recNum
    } else {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "接收号码不合法", nil))
    }
}

func (q *smsQuery) SetQueryDate(queryDate string) {
    if len(queryDate) == 8 {
        q.queryDate = queryDate
    } else {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "发送日期不合法", nil))
    }
}

func (q *smsQuery) SetPage(page uint) {
    if page > 0 {
        q.ReqData["current_page"] = strconv.Itoa(int(page))
    } else {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "页码必须大于0", nil))
    }
}

func (q *smsQuery) SetLimit(limit uint) {
    if (limit > 0) && (limit <= 50) {
        q.ReqData["page_size"] = strconv.Itoa(int(limit))
    } else {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "每页数量必须在1-50之间", nil))
    }
}

func (q *smsQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(q.recNum) == 0 {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "接收号码必须填写", nil))
    }
    if len(q.queryDate) == 0 {
        panic(mperr.NewSmsDaYu(errorcode.SmsDaYuParam, "发送日期必须填写", nil))
    }
    q.ReqData["rec_num"] = q.recNum
    q.ReqData["query_date"] = q.queryDate

    return q.GetRequest()
}

func NewSmsQuery() *smsQuery {
    query := &smsQuery{taobao.NewBaseTaoBao(), "", "", "", 1, 10}
    conf := sms.NewConfigDaYu()
    query.AppKey = conf.GetAppKey()
    query.AppSecret = conf.GetAppSecret()
    query.ReqData["current_page"] = strconv.Itoa(int(query.page))
    query.ReqData["page_size"] = strconv.Itoa(int(query.limit))
    query.SetMethod("alibaba.aliqin.fc.mpsms.num.query")
    return query
}
