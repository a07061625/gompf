/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 22:13
 */
package datacube

import (
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type statistics struct {
    wx.BaseWxAccount
    appId     string
    beginDate string // 起始日期
    endDate   string // 结束日期
    maxDate   int    // 最大日期天数
}

func (s *statistics) SetStatType(statType string) {
    switch statType {
    case "0100": // 用户每日变动数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getusersummary?access_token="
        s.maxDate = 7
    case "0101": // 用户每日总数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getusercumulate?access_token="
        s.maxDate = 7
    case "0200": // 图文每日发送数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getarticlesummary?access_token="
        s.maxDate = 1
    case "0201": // 图文每日发送总数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getarticletotal?access_token="
        s.maxDate = 1
    case "0202": // 图文每日浏览数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getuserread?access_token="
        s.maxDate = 3
    case "0203": // 图文每小时浏览数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getuserreadhour?access_token="
        s.maxDate = 1
    case "0204": // 图文每日分享数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getusershare?access_token="
        s.maxDate = 7
    case "0205": // 图文每小时分享数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getusersharehour?access_token="
        s.maxDate = 1
    case "0300": // 消息每日发送数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getupstreammsg?access_token="
        s.maxDate = 7
    case "0301": // 消息每小时发送数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getupstreammsghour?access_token="
        s.maxDate = 1
    case "0302": // 消息每周发送数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getupstreammsgweek?access_token="
        s.maxDate = 30
    case "0303": // 消息每月发送数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getupstreammsgmonth?access_token="
        s.maxDate = 30
    case "0304": // 消息每日分布数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getupstreammsgdist?access_token="
        s.maxDate = 15
    case "0305": // 消息每周分布数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getupstreammsgdistweek?access_token="
        s.maxDate = 30
    case "0306": // 消息每月分布数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getupstreammsgdistmonth?access_token="
        s.maxDate = 30
    case "0400": // 接口每日调用数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getinterfacesummary?access_token="
        s.maxDate = 30
    case "0401": // 接口每小时调用数
        s.ReqUrl = "https://api.weixin.qq.com/datacube/getinterfacesummaryhour?access_token="
        s.maxDate = 1
    default:
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "统计类型不支持", nil))
    }
}

func (s *statistics) SetDate(beginTime, endTime string) {
    nowTime := time.Now()
    dayTime := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, nowTime.Location()).Unix()
    beginDay, _ := time.ParseInLocation("2006-01-02", beginTime, nowTime.Location())
    endDay, _ := time.ParseInLocation("2006-01-02", endTime, nowTime.Location())
    if beginDay.Unix() < 1417363200 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "起始时间不能小于2014年12月1日", nil))
    } else if endDay.Unix() < 1417363200 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "结束时间不能小于2014年12月1日", nil))
    } else if beginDay.Unix() > beginDay.Unix() {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "起始时间不能大于结束时间", nil))
    } else if beginDay.Unix() >= dayTime {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "起始时间必须小于今天", nil))
    } else if endDay.Unix() >= dayTime {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "结束时间必须小于今天", nil))
    }

    dayNum := int((endDay.Unix() - beginDay.Unix()) / 86400)
    if dayNum >= s.maxDate {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "时间跨度必须小于最大限定天数", nil))
    }

    s.beginDate = beginTime
    s.endDate = endTime
}

func (s *statistics) checkData() {
    if s.maxDate <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "统计类型不能为空", nil))
    }
    if len(s.beginDate) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "起始时间不能为空", nil))
    }
    s.ReqData["begin_date"] = s.beginDate
    s.ReqData["end_date"] = s.endDate
}

func (s *statistics) SendRequest() api.ApiResult {
    s.checkData()

    reqBody := mpf.JsonMarshal(s.ReqData)
    s.ReqUrl += wx.NewUtilWx().GetSingleAccessToken(s.appId)
    client, req := s.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := s.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewStatistics(appId string) *statistics {
    s := &statistics{wx.NewBaseWxAccount(), "", "", "", 0}
    s.appId = appId
    s.maxDate = 0
    s.ReqContentType = project.HTTPContentTypeJSON
    s.ReqMethod = fasthttp.MethodPost
    return s
}
