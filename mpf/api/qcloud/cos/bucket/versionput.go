/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/22 0022
 * Time: 14:31
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

// 启用或者暂停存储桶的版本控制
type versionPut struct {
    qcloud.BaseCos
    versionConfig map[string]interface{} // 版本配置
}

func (vp *versionPut) SetVersionConfig(versionConfig map[string]interface{}) {
    if len(versionConfig) > 0 {
        vp.versionConfig = versionConfig
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "版本配置不合法", nil))
    }
}

func (vp *versionPut) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(vp.versionConfig) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "版本配置不能为空", nil))
    }

    xmlData := mxj.Map(vp.versionConfig)
    xmlStr, _ := xmlData.Xml("VersioningConfiguration")
    reqBody := project.DataPrefixXml + string(xmlStr)
    vp.ReqUrl = "http://" + vp.ReqHeader["Host"] + vp.ReqUri + "?versioning"
    client, req := vp.GetRequest()
    req.SetBody([]byte(reqBody))
    return client, req
}

func NewVersionPut() *versionPut {
    vp := &versionPut{qcloud.NewCos(), make(map[string]interface{})}
    vp.ReqMethod = fasthttp.MethodPut
    vp.SetParamData("versioning", "")
    return vp
}
