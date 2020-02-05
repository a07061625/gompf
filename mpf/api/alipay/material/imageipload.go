/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 17:49
 */
package material

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 上传门店照片
type imageUpload struct {
    alipay.BaseAliPay
    imageType    string // 图片格式
    imageName    string // 图片名称
    imageContent string // 图片内容
}

func (iu *imageUpload) SetImageType(imageType string) {
    if (len(imageType) > 0) && (len(imageType) <= 8) {
        iu.imageType = imageType
    } else {
        panic(mperr.NewAliPayMaterial(errorcode.AliPayMaterialParam, "图片格式不合法", nil))
    }
}

func (iu *imageUpload) SetImageName(imageName string) {
    if (len(imageName) > 0) && (len(imageName) <= 128) {
        iu.imageName = imageName
    } else {
        panic(mperr.NewAliPayMaterial(errorcode.AliPayMaterialParam, "图片名称不合法", nil))
    }
}

func (iu *imageUpload) SetImageContent(imageContent string) {
    if len(imageContent) > 0 {
        iu.imageContent = imageContent
    } else {
        panic(mperr.NewAliPayMaterial(errorcode.AliPayMaterialParam, "图片内容不合法", nil))
    }
}

func (iu *imageUpload) SetImagePid(imagePid string) {
    match, _ := regexp.MatchString(`^[0-9]{1,16}$`, imagePid)
    if match {
        iu.BizContent["image_pid"] = imagePid
    } else {
        panic(mperr.NewAliPayMaterial(errorcode.AliPayMaterialParam, "图片所属的partnerId不合法", nil))
    }
}

func (iu *imageUpload) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(iu.imageType) == 0 {
        panic(mperr.NewAliPayMaterial(errorcode.AliPayMaterialParam, "图片格式不能为空", nil))
    }
    if len(iu.imageName) == 0 {
        panic(mperr.NewAliPayMaterial(errorcode.AliPayMaterialParam, "图片名称不能为空", nil))
    }
    if len(iu.imageContent) == 0 {
        panic(mperr.NewAliPayMaterial(errorcode.AliPayMaterialParam, "图片内容不能为空", nil))
    }
    iu.BizContent["image_type"] = iu.imageType
    iu.BizContent["image_name"] = iu.imageName
    iu.BizContent["image_content"] = iu.imageContent

    return iu.GetRequest()
}

func NewImageUpload(appId string) *imageUpload {
    iu := &imageUpload{alipay.NewBase(appId), "", "", ""}
    iu.SetMethod("alipay.offline.material.image.upload")
    return iu
}
