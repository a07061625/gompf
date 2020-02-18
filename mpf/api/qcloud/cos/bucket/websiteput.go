/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 11:07
 */
package bucket

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/clbanning/mxj"
    "github.com/valyala/fasthttp"
)

// 设置存储桶的静态网站配置
type websitePut struct {
    qcloud.BaseCos
    websiteConfig map[string]interface{} // 静态网站配置
}

func (wp *websitePut) SetWebsiteConfig(websiteConfig map[string]interface{}) {
    if len(websiteConfig) > 0 {
        wp.websiteConfig = websiteConfig
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "静态网站配置不合法", nil))
    }
}

func (wp *websitePut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(wp.websiteConfig) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "静态网站配置不能为空", nil))
    }

    xmlData := mxj.Map(wp.websiteConfig)
    xmlStr, _ := xmlData.Xml("WebsiteConfiguration")
    reqBody := project.DataPrefixXML + string(xmlStr)
    wp.ReqUrl = "http://" + wp.ReqHeader["Host"] + wp.ReqUri + "?website"
    client, req := wp.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewWebsitePut() *websitePut {
    wp := &websitePut{qcloud.NewCos(), make(map[string]interface{})}
    wp.ReqMethod = fasthttp.MethodPut
    wp.SetParamData("website", "")
    return wp
}
